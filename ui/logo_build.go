package ui

import (
	"math"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

const (
	animLogoBuild = "logoBuild"

	logoBuildDuration   = 500 * time.Millisecond
	logoBuildFadeWindow = 0.35
)

// renderLogoBuild fades the logo in from the centre outward.
func renderLogoBuild(logo string, base colorful.Color, progress float64) string {
	lines, maxWidth := textBounds(logo)
	cx := float64(maxWidth-1) / 2
	cy := float64(len(lines)-1) / 2
	dMax := math.Max(cx, cy)
	if dMax <= 0 {
		dMax = 1
	}

	return renderFade(logo, base, fadeOpts{
		progress: progress,
		window:   logoBuildFadeWindow,
		norm: func(row, col int) float64 {
			dx := math.Abs(float64(col) - cx)
			dy := math.Abs(float64(row) - cy)
			return math.Min(math.Max(dx, dy)/dMax, 1)
		},
	})
}
