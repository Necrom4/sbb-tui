package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// textBounds returns the lines of `text` and the visible width of the
// widest one.
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

// renderFade paints `text` against a black→base gradient, with each
// non-space rune fading in over its own window.
func renderFade(text string, base colorful.Color, opts fadeOpts) string {
	pal := newPalette(colorBlack, base)
	return renderGrid(text, pal, func(row, col int) float64 {
		return windowedFade(opts.progress, opts.norm(row, col), opts.window, opts.shift)
	})
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
