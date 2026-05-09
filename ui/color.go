package ui

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

var (
	colorBlack = colorful.Color{R: 0, G: 0, B: 0}
	colorWhite = colorful.Color{R: 1, G: 1, B: 1}
)

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

// shineTarget returns the neutral toward which a shine effect should
// blend `base`: black for already-bright bases, white for darker ones.
func shineTarget(base colorful.Color) colorful.Color {
	if relLuminance(base) >= shineLumaPivot {
		return colorBlack
	}
	return colorWhite
}

// parseColor converts a theme color string (hex or named) to a colorful.Color.
func parseColor(s string) colorful.Color {
	if c, err := colorful.Hex(s); err == nil {
		return c
	}
	switch s {
	case "white":
		return colorWhite
	case "black":
		return colorBlack
	case "red":
		return colorful.Color{R: 1, G: 0, B: 0}
	case "green":
		return colorful.Color{R: 0, G: 1, B: 0}
	case "blue":
		return colorful.Color{R: 0, G: 0, B: 1}
	case "yellow":
		return colorful.Color{R: 1, G: 1, B: 0}
	case "cyan":
		return colorful.Color{R: 0, G: 1, B: 1}
	case "magenta":
		return colorful.Color{R: 1, G: 0, B: 1}
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
