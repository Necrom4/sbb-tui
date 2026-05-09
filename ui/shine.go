package ui

import (
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	animLogoShine = "logoShine"
	animTextShine = "textShine"
)

const (
	shineDuration      = 800 * time.Millisecond
	shineRepeatGap     = 2 * time.Second
	shineDelta         = 0.30 // max lightness shift at the band's center (0..1)
	shineGradientSize  = 65
	logoShineBandWidth = 28.0
	textShineBandWidth = 8.0

	// shineLumaPivot picks shine polarity: bases brighter than this darken, the rest brighten.
	shineLumaPivot = 0.6
)

// Color and rendering primitives shared by every animation effect.

// relLuminance returns the WCAG relative luminance of c in [0, 1].
func relLuminance(c colorful.Color) float64 {
	return 0.2126*linearize(c.R) + 0.7152*linearize(c.G) + 0.0722*linearize(c.B)
}

func linearize(channel float64) float64 {
	if channel <= 0.03928 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

// shineTarget picks the neutral toward which a shine should blend
// `base`: black for already-bright bases, white for darker ones.
func shineTarget(base colorful.Color) colorful.Color {
	if relLuminance(base) >= shineLumaPivot {
		return colorBlack
	}
	return colorWhite
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

type paletteCell struct {
	prefix, suffix string
}

// palette holds pre-rendered ANSI prefix/suffix pairs for each step
// of a Lab-blended gradient between two colors.
type palette [shineGradientSize]paletteCell

// newPalette builds a gradient blending from `from` (factor 0) to `to` (factor 1) in Lab space.
func newPalette(from, to colorful.Color) palette {
	var p palette
	for i := range p {
		f := float64(i) / float64(shineGradientSize-1)
		c := from.BlendLab(to, f).Clamped()
		p[i] = makePaletteCell(lipgloss.Color(c.Hex()))
	}
	return p
}

func makePaletteCell(c lipgloss.Color) paletteCell {
	const sentinel = "\x00"
	rendered := lipgloss.NewStyle().Foreground(c).Render(sentinel)
	prefix, suffix, ok := strings.Cut(rendered, sentinel)
	if !ok {
		return paletteCell{}
	}
	return paletteCell{prefix: prefix, suffix: suffix}
}

// render returns the rune wrapped in the SGR pair for the given factor.
func (p palette) render(factor float64, r rune) string {
	idx := int(math.Round(clamp01(factor) * float64(shineGradientSize-1)))
	cell := p[idx]
	return cell.prefix + string(r) + cell.suffix
}

// textBounds returns the lines of `text` and the visible width of the widest one.
func textBounds(text string) (lines []string, maxWidth int) {
	lines = strings.Split(text, "\n")
	for _, ln := range lines {
		if w := lipgloss.Width(ln); w > maxWidth {
			maxWidth = w
		}
	}
	return lines, maxWidth
}

// renderGrid walks `text` line by line and writes each non-space rune
// styled by `pal` at the factor returned by `factor(row, col)`.
// Spaces are emitted unchanged; lines are right-padded to the widest
// line so lipgloss centring does not shift shorter rows.
func renderGrid(text string, pal palette, factor func(row, col int) float64) string {
	lines, maxWidth := textBounds(text)

	var b strings.Builder
	b.Grow(len(text) * 6)
	for row, line := range lines {
		col := 0
		for _, r := range line {
			if r == ' ' {
				b.WriteRune(r)
				col++
				continue
			}
			b.WriteString(pal.render(factor(row, col), r))
			col++
		}
		if pad := maxWidth - col; pad > 0 {
			b.WriteString(strings.Repeat(" ", pad))
		}
		if row < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// fadeOpts configures a windowed-fade render: each glyph reveals over
// a small `window` of progress, with its position computed by `norm`
// (0 = first to appear, 1 = last). `shift` slides every cell's window
// by a constant amount.
type fadeOpts struct {
	progress float64
	window   float64
	shift    float64
	norm     func(row, col int) float64
}

// renderFade paints `text` with each non-space rune fading in over
// its own window. When the terminal background is known the gradient
// runs from that background colour to `base`, producing a true
// invisible-to-base fade. Otherwise cells stay as spaces until their
// window opens and only then start fading from a dark anchor, so we
// never flash a stale colour against a mismatched background.
func (m appModel) renderFade(text string, base colorful.Color, opts fadeOpts) string {
	if m.styles.backgroundKnown {
		pal := newPalette(m.styles.background, base)
		return renderGrid(text, pal, func(row, col int) float64 {
			return windowedFade(opts.progress, opts.norm(row, col), opts.window, opts.shift)
		})
	}
	return renderFadeFallback(text, base, opts)
}

// renderFadeFallback runs the fade when the terminal background is
// not known: cells are spaces before their window opens, then fade
// from a dark anchor to `base` over the window.
func renderFadeFallback(text string, base colorful.Color, opts fadeOpts) string {
	pal := newPalette(colorBlack, base)
	lines, maxWidth := textBounds(text)

	var b strings.Builder
	b.Grow(len(text) * 6)
	for row, line := range lines {
		col := 0
		for _, r := range line {
			if r == ' ' {
				b.WriteRune(r)
				col++
				continue
			}
			f := windowedFade(opts.progress, opts.norm(row, col), opts.window, opts.shift)
			if f <= 0 {
				b.WriteByte(' ')
			} else {
				b.WriteString(pal.render(f, r))
			}
			col++
		}
		if pad := maxWidth - col; pad > 0 {
			b.WriteString(strings.Repeat(" ", pad))
		}
		if row < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// windowedFade returns a smoothstepped factor in [0, 1] for a glyph
// whose reveal window is centred around `norm * (1 - window) + shift`.
func windowedFade(progress, norm, window, shift float64) float64 {
	start := norm*(1-window) + shift
	end := start + window
	if progress <= start {
		return 0
	}
	if progress >= end {
		return 1
	}
	t := (progress - start) / window
	return t * t * (3 - 2*t)
}

// Shine animation orchestration.

type shineRestartMsg struct{}

func shineRestartCmd() tea.Cmd {
	return tea.Tick(shineRepeatGap, func(time.Time) tea.Msg {
		return shineRestartMsg{}
	})
}

// startShineCycle restarts both shine animations together, keeping their progress in lockstep.
func (m *appModel) startShineCycle() tea.Cmd {
	return tea.Batch(
		m.anim.Start(animLogoShine, shineDuration),
		m.anim.Start(animTextShine, shineDuration),
	)
}

// onAnimationsFinished walks the start-screen animation chain in
// response to the named animations completing on this tick.
func (m *appModel) onAnimationsFinished(finished []string) []tea.Cmd {
	var cmds []tea.Cmd
	var shineRestarted bool
	for _, name := range finished {
		switch name {
		case animLogoBuild:
			cmds = append(cmds, m.anim.Start(animTaglineBuild, taglineBuildDuration))
		case animTaglineBuild:
			cmds = append(cmds, m.startShineCycle())
		case animLogoShine, animTextShine:
			if !shineRestarted {
				cmds = append(cmds, shineRestartCmd())
				shineRestarted = true
			}
		}
	}
	return cmds
}

// shineDirection is the axis along which the band sweeps.
type shineDirection int

const (
	shineDiagonal shineDirection = iota
	shineHorizontal
)

type shineOpts struct {
	base      colorful.Color
	progress  float64
	bandWidth float64
	direction shineDirection
}

// applyShine paints each non-space rune of `s` with a color whose
// brightness depends on its position relative to a moving band.
func applyShine(s string, opts shineOpts) string {
	lines, maxWidth := textBounds(s)

	dMax := float64(maxWidth - 1)
	if opts.direction == shineDiagonal {
		dMax = float64(maxWidth + len(lines) - 2)
	}
	dCenter := -opts.bandWidth + opts.progress*(dMax+2*opts.bandWidth)

	return renderGrid(s, newShinePalette(opts.base), func(row, col int) float64 {
		d := float64(col)
		if opts.direction == shineDiagonal {
			d = float64(col + row)
		}
		return shineFactor(d-dCenter, opts.bandWidth)
	})
}

// shineFactor maps signed distance from the band center to a band intensity in [0, 1].
func shineFactor(signedDist, bandWidth float64) float64 {
	dist := math.Abs(signedDist)
	if dist >= bandWidth {
		return 0
	}
	return math.Cos(dist / bandWidth * math.Pi / 2)
}

// newShinePalette builds the gradient used by the shine effect: from
// the base color toward the auto-picked neutral, scaled so the peak
// only travels `shineDelta` of the way.
func newShinePalette(base colorful.Color) palette {
	peak := base.BlendLab(shineTarget(base), shineDelta)
	return newPalette(base, peak)
}

func (m appModel) renderLogo(logo string) string {
	if !m.animations {
		return m.styles.logo.Render(logo)
	}
	if progress, active := m.anim.Progress(animLogoBuild); active {
		return m.renderLogoBuild(logo, m.styles.logoBase, progress)
	}
	progress, active := m.anim.Progress(animLogoShine)
	if !active {
		return m.styles.logo.Render(logo)
	}
	return applyShine(logo, shineOpts{
		base:      m.styles.logoBase,
		progress:  progress,
		bandWidth: logoShineBandWidth,
		direction: shineDiagonal,
	})
}

func (m appModel) renderStartTagline(text string) string {
	if !m.animations {
		return m.styles.textMuted.Render(text)
	}
	// Hide the tagline until its build animation has been registered,
	// so the logo build plays alone first.
	if !m.anim.Registered(animTaglineBuild) {
		return strings.Repeat(" ", lipgloss.Width(text))
	}
	if progress, active := m.anim.Progress(animTaglineBuild); active {
		return m.renderTaglineBuild(text, m.styles.textMutedBase, progress)
	}
	progress, active := m.anim.Progress(animTextShine)
	if !active {
		return m.styles.textMuted.Render(text)
	}
	return applyShine(text, shineOpts{
		base:      m.styles.textMutedBase,
		progress:  progress,
		bandWidth: textShineBandWidth,
		direction: shineHorizontal,
	})
}
