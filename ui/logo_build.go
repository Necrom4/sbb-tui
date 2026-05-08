package ui

import (
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	animLogoBuild = "logoBuild"

	logoBuildDuration = 500 * time.Millisecond

	// logoBuildFadeWindow is the fraction of the total animation each
	// character spends fading in. 0.35 means a glyph takes ~175ms to
	// go from invisible to fully colored.
	logoBuildFadeWindow = 0.35

	// logoBuildPaletteLen controls fade smoothness; the palette goes
	// from the terminal background (black) to the base logo color.
	// Reuses shinePaletteLen so the existing shinePalette type can
	// host the entries without bounds gymnastics.
	logoBuildPaletteLen = shinePaletteLen
)

// logoBuildFinished reports whether the build animation just finished
// on this tick. Used to trigger the follow-up shines.
func logoBuildFinished(finished []string) bool {
	for _, name := range finished {
		if name == animLogoBuild {
			return true
		}
	}
	return false
}

// renderLogoBuild returns the logo rendered at the given build
// progress in [0, 1]. Characters fade in (Lab-blended from black to
// the base color) starting from the logo's center and spreading
// outward. Whitespace cells remain spaces throughout.
func renderLogoBuild(logo string, base colorful.Color, progress float64) string {
	lines := strings.Split(logo, "\n")
	maxWidth := 0
	for _, ln := range lines {
		if w := lipgloss.Width(ln); w > maxWidth {
			maxWidth = w
		}
	}
	cx := float64(maxWidth-1) / 2
	cy := float64(len(lines)-1) / 2
	// Maximum Chebyshev distance from center to any corner.
	dMax := math.Max(cx, cy)
	if dMax <= 0 {
		dMax = 1
	}

	palette := buildLogoBuildPalette(base)

	var b strings.Builder
	b.Grow(len(logo) * 6)
	for row, line := range lines {
		col := 0
		for _, r := range line {
			if r == ' ' {
				b.WriteRune(r)
				col++
				continue
			}
			// Chebyshev (max-norm) distance from logo center,
			// normalized to [0, 1]. Center → 0, corners → 1.
			dx := math.Abs(float64(col) - cx)
			dy := math.Abs(float64(row) - cy)
			norm := math.Max(dx, dy) / dMax
			if norm > 1 {
				norm = 1
			}

			factor := logoBuildFactor(progress, norm)
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

// logoBuildFactor returns a per-cell brightness factor in [0, 1]
// where 0 = fully transparent (background color) and 1 = fully
// rendered at the base color. The cell starts fading in when
// `progress` enters its reveal window.
func logoBuildFactor(progress, norm float64) float64 {
	// Each cell's window starts at `start` and ends `fadeWindow`
	// later, with start placed so the corner cells finish exactly
	// at progress=1.
	start := norm * (1 - logoBuildFadeWindow)
	end := start + logoBuildFadeWindow
	if progress <= start {
		return 0
	}
	if progress >= end {
		return 1
	}
	t := (progress - start) / logoBuildFadeWindow
	// Smoothstep for a gentler fade.
	return t * t * (3 - 2*t)
}

// buildLogoBuildPalette generates color steps blending from black
// (factor=0, "transparent") to the base color (factor=1).
func buildLogoBuildPalette(base colorful.Color) shinePalette {
	var p shinePalette
	bg := colorful.Color{R: 0, G: 0, B: 0}
	for i := 0; i < logoBuildPaletteLen; i++ {
		f := float64(i) / float64(logoBuildPaletteLen-1)
		c := bg.BlendLab(base, f).Clamped()
		hex := c.Hex()
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
