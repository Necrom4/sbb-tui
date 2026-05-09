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

// renderTaglineBuild types text left-to-right, fading each rune in.
// The reveal window is shifted so character 0 starts mid-fade and the
// animation feels continuous with the logo build that precedes it.
func (m appModel) renderTaglineBuild(text string, base colorful.Color, progress float64) string {
	n := utf8.RuneCountInString(text)
	if n == 0 {
		return text
	}
	denom := float64(n)
	return m.renderFade(text, base, fadeOpts{
		progress: progress,
		window:   taglineBuildFadeWindow,
		shift:    -taglineBuildFadeWindow / 2,
		norm:     func(_, col int) float64 { return float64(col) / denom },
	})
}
