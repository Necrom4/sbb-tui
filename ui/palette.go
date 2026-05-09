package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// shineGradientSize is the number of pre-built color steps animations
// pick from. Higher values produce smoother gradients at the cost of
// more SGR escape sequences cached per palette.
const shineGradientSize = 65

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
	idx := strings.Index(rendered, sentinel)
	if idx < 0 {
		return paletteCell{}
	}
	return paletteCell{prefix: rendered[:idx], suffix: rendered[idx+len(sentinel):]}
}

// render returns the rune wrapped in the SGR pair for the given factor.
func (p palette) render(factor float64, r rune) string {
	idx := int(math.Round(clamp01(factor) * float64(shineGradientSize-1)))
	cell := p[idx]
	return cell.prefix + string(r) + cell.suffix
}
