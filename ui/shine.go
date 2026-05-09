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

// ----- Color and rendering primitives shared by every animation effect -----

// relLuminance returns the WCAG relative luminance of c in [0, 1].
func relLuminance(c colorful.Color) float64 {
	return 0.2126*linearize(c.R) + 0.7152*linearize(c.G) + 0.0722*linearize(c.B)
}

// linearize converts an sRGB channel to linear-light per the WCAG formula.
func linearize(channel float64) float64 {
	if channel <= 0.03928 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

// shineTarget returns the neutral the shine should blend toward: black
// for already-bright bases, white for darker ones.
func shineTarget(base colorful.Color) colorful.Color {
	if relLuminance(base) >= shineLumaPivot {
		return colorBlack
	}
	return colorWhite
}

// clamp01 clamps x into the [0, 1] range.
func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// paletteCell holds the SGR prefix and reset suffix for a single gradient step.
type paletteCell struct {
	prefix, suffix string
}

// palette is a pre-rendered Lab gradient between two colors,
// indexed by a [0, 1] factor (mapped onto shineGradientSize steps).
type palette [shineGradientSize]paletteCell

// newPalette builds a Lab-space gradient from `from` (factor 0) to `to` (factor 1).
func newPalette(from, to colorful.Color) palette {
	var p palette
	for i := range p {
		f := float64(i) / float64(shineGradientSize-1)
		c := from.BlendLab(to, f).Clamped()
		p[i] = makePaletteCell(lipgloss.Color(c.Hex()))
	}
	return p
}

// makePaletteCell extracts the SGR prefix/suffix lipgloss emits for one color.
func makePaletteCell(c lipgloss.Color) paletteCell {
	const sentinel = "\x00"
	rendered := lipgloss.NewStyle().Foreground(c).Render(sentinel)
	prefix, suffix, ok := strings.Cut(rendered, sentinel)
	if !ok {
		return paletteCell{}
	}
	return paletteCell{prefix: prefix, suffix: suffix}
}

// render returns r wrapped in the SGR pair for the gradient step at factor.
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

// renderGrid styles each non-space rune of text via pal at factor(row, col).
// Lines are right-padded to the widest line so JoinVertical's centring
// does not shift shorter rows.
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
// `window` of progress, with its position computed by `norm` (0 = first
// to appear, 1 = last). `shift` slides every cell's window by a constant.
type fadeOpts struct {
	progress float64
	window   float64
	shift    float64
	norm     func(row, col int) float64
}

// renderFade reveals each non-space rune over its own windowed fade. Cells
// stay as literal spaces until the window opens, then fade from the detected
// terminal background (or black when unknown) to base. The space gate keeps
// transparent terminals fully pass-through during the build.
func (m appModel) renderFade(text string, base colorful.Color, opts fadeOpts) string {
	from := colorBlack
	if m.styles.backgroundKnown {
		from = m.styles.background
	}
	pal := newPalette(from, base)
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
// whose reveal window is centred around norm * (1 - window) + shift.
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

// ----- Shine animation orchestration -----

// shineRestartMsg fires after the post-cycle gap to restart the shine cycle.
type shineRestartMsg struct{}

// shineRestartCmd schedules a shineRestartMsg shineRepeatGap from now.
func shineRestartCmd() tea.Cmd {
	return tea.Tick(shineRepeatGap, func(time.Time) tea.Msg {
		return shineRestartMsg{}
	})
}

// startShineCycle restarts both shine animations together so they stay in lockstep.
func (m *appModel) startShineCycle() tea.Cmd {
	return tea.Batch(
		m.anim.Start(animLogoShine, shineDuration),
		m.anim.Start(animTextShine, shineDuration),
	)
}

// onAnimationsFinished advances the start-screen animation chain when the
// named animations complete on the current tick.
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
			// Both shines finish on the same tick; only schedule one restart.
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

// shineOpts configures one shine pass.
type shineOpts struct {
	base      colorful.Color
	progress  float64
	bandWidth float64
	direction shineDirection
}

// applyShine paints each non-space rune of s with a color whose
// brightness depends on its position relative to a moving band.
func applyShine(s string, opts shineOpts) string {
	lines, maxWidth := textBounds(s)

	dMax := float64(maxWidth - 1)
	if opts.direction == shineDiagonal {
		dMax = float64(maxWidth + len(lines) - 2)
	}
	// Start the band fully off-screen at progress 0 and end it fully off-screen at 1.
	dCenter := -opts.bandWidth + opts.progress*(dMax+2*opts.bandWidth)

	return renderGrid(s, newShinePalette(opts.base), func(row, col int) float64 {
		d := float64(col)
		if opts.direction == shineDiagonal {
			d = float64(col + row)
		}
		return shineFactor(d-dCenter, opts.bandWidth)
	})
}

// shineFactor maps signed distance from the band center to a [0, 1] band intensity.
func shineFactor(signedDist, bandWidth float64) float64 {
	dist := math.Abs(signedDist)
	if dist >= bandWidth {
		return 0
	}
	return math.Cos(dist / bandWidth * math.Pi / 2)
}

// newShinePalette returns the gradient used by the shine effect: base
// to base-blended-toward-neutral by shineDelta.
func newShinePalette(base colorful.Color) palette {
	peak := base.BlendLab(shineTarget(base), shineDelta)
	return newPalette(base, peak)
}

// renderLogo returns the logo styled, building or shining as the model state dictates.
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

// renderStartTagline returns the start-screen tagline styled, building or
// shining alongside the logo. While the logo is still building, the
// tagline is intentionally rendered as blank spaces.
func (m appModel) renderStartTagline(text string) string {
	if !m.animations {
		return m.styles.textMuted.Render(text)
	}
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
