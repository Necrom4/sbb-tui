package ui

import (
	"time"
	"unicode/utf8"

	"github.com/lucasb-eyer/go-colorful"
)

const (
	animTaglineBuild = "taglineBuild"

	taglineBuildDuration   = 350 * time.Millisecond
	taglineBuildFadeWindow = 0.20
)

// renderTaglineBuild types `text` left-to-right, with each rune
// fading in over a short window. Character 0 starts mid-fade so the
// reveal feels continuous with whatever animation preceded it.
func renderTaglineBuild(text string, base colorful.Color, progress float64) string {
	n := utf8.RuneCountInString(text)
	if n == 0 {
		return text
	}
	denom := float64(n)
	return renderFade(text, base, fadeOpts{
		progress: progress,
		window:   taglineBuildFadeWindow,
		shift:    -taglineBuildFadeWindow / 2,
		norm:     func(_, col int) float64 { return float64(col) / denom },
	})
}
