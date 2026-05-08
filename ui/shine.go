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
	shinePaletteLen    = 65   // pre-built color steps; higher = smoother
	logoShineBandWidth = 28.0
	textShineBandWidth = 8.0

	// shineLumaPivot is the perceptual-luminance threshold above which
	// the shine band darkens the base color, and below which it
	// brightens it. Sits a bit above 0.5 because saturated colors
	// already feel vivid and reading "shine" as brighter on near-white
	// looks washed out.
	shineLumaPivot = 0.6
)

type shineRestartMsg struct{}

func shineRestartCmd() tea.Cmd {
	return tea.Tick(shineRepeatGap, func(time.Time) tea.Msg {
		return shineRestartMsg{}
	})
}

// shineCycleFinished reports whether the finished slice contains a
// shine-cycle animation. Both shines share the same duration so they
// complete on the same tick; firing the restart on either one is
// enough.
func shineCycleFinished(finished []string) bool {
	for _, name := range finished {
		if name == animLogoShine || name == animTextShine {
			return true
		}
	}
	return false
}

// shineDirection is the axis along which the band sweeps.
type shineDirection int

const (
	shineDiagonal shineDirection = iota
	shineHorizontal
)

// shineOpts bundles all per-effect knobs.
type shineOpts struct {
	base      colorful.Color
	progress  float64
	bandWidth float64
	direction shineDirection
}

// applyShine paints each non-space rune of `s` with a color whose
// brightness depends on its position relative to a moving band. The
// band's polarity (lighter or darker than the base) is decided
// automatically from the base color's perceptual luminance.
func applyShine(s string, opts shineOpts) string {
	lines := strings.Split(s, "\n")
	maxWidth := 0
	for _, ln := range lines {
		if w := lipgloss.Width(ln); w > maxWidth {
			maxWidth = w
		}
	}

	var dMin, dMax float64
	switch opts.direction {
	case shineHorizontal:
		dMin = 0
		dMax = float64(maxWidth - 1)
	default: // shineDiagonal
		dMin = 0
		dMax = float64(maxWidth + len(lines) - 2)
	}
	span := (dMax - dMin) + 2*opts.bandWidth
	dCenter := dMin - opts.bandWidth + opts.progress*span

	palette := buildShinePalette(opts.base)

	var b strings.Builder
	b.Grow(len(s) * 6)
	for row, line := range lines {
		col := 0
		for _, r := range line {
			if r == ' ' {
				b.WriteRune(r)
				col++
				continue
			}
			var d float64
			switch opts.direction {
			case shineHorizontal:
				d = float64(col)
			default:
				d = float64(col + row)
			}
			factor := shineFactor(d-dCenter, opts.bandWidth)
			b.WriteString(palette.render(factor, r))
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

// shineFactor maps signed distance from the band center to a band intensity in [0, 1].
func shineFactor(signedDist, bandWidth float64) float64 {
	dist := math.Abs(signedDist)
	if dist >= bandWidth {
		return 0
	}
	return math.Cos(dist / bandWidth * math.Pi / 2)
}

type shinePalette struct {
	cells [shinePaletteLen]paletteCell
}

type paletteCell struct {
	prefix string
	suffix string
}

// buildShinePalette generates color steps from base (factor=0) to the
// shine peak (factor=1). The peak is produced by blending the base
// toward black or white in Lab space, which naturally desaturates as
// it darkens/brightens. This avoids HSL's tendency to reveal a
// saturated hue when shifting the lightness of a near-neutral color
// (e.g. shining #E9F7EF would otherwise pull a vivid green band).
func buildShinePalette(base colorful.Color) shinePalette {
	var p shinePalette
	target := colorful.Color{R: 1, G: 1, B: 1}
	if relLuminance(base) >= shineLumaPivot {
		target = colorful.Color{R: 0, G: 0, B: 0}
	}
	for i := 0; i < shinePaletteLen; i++ {
		f := float64(i) / float64(shinePaletteLen-1)
		c := base.BlendLab(target, f*shineDelta)
		hex := c.Clamped().Hex()
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
		const sentinel = "\x00"
		rendered := style.Render(sentinel)
		idx := strings.Index(rendered, sentinel)
		if idx < 0 {
			p.cells[i] = paletteCell{}
			continue
		}
		p.cells[i] = paletteCell{
			prefix: rendered[:idx],
			suffix: rendered[idx+len(sentinel):],
		}
	}
	return p
}

func (p shinePalette) render(factor float64, r rune) string {
	if factor < 0 {
		factor = 0
	} else if factor > 1 {
		factor = 1
	}
	idx := int(math.Round(factor * float64(shinePaletteLen-1)))
	cell := p.cells[idx]
	return cell.prefix + string(r) + cell.suffix
}

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

func (m appModel) renderLogo(logo string) string {
	if !m.animations {
		return m.styles.logo.Render(logo)
	}
	if progress, active := m.anim.Progress(animLogoBuild); active {
		return renderLogoBuild(logo, m.styles.logoBase, progress)
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
	// Hide the tagline until its build animation has started, so the
	// logo build animation plays alone first.
	if !m.anim.Registered(animTaglineBuild) {
		return strings.Repeat(" ", lipgloss.Width(text))
	}
	if progress, active := m.anim.Progress(animTaglineBuild); active {
		return renderTaglineBuild(text, m.styles.textMutedBase, progress)
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
