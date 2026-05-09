package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	animLoading      = "loading"
	loadingFrameTime = 80 * time.Millisecond
)

var loadingFrames = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

// startLoadingCmd starts the spinner animation; returns nil when animations are disabled.
func (m *appModel) startLoadingCmd() tea.Cmd {
	if !m.animations {
		return nil
	}
	return m.anim.StartIndefinite(animLoading)
}

// renderLoading returns the "Searching connections" line with a spinner glyph.
func (m appModel) renderLoading() string {
	if !m.animations {
		return "\n  Searching connections..."
	}
	frame := loadingFrames[0]
	if elapsed, ok := m.anim.Elapsed(animLoading); ok {
		frame = loadingFrames[int(elapsed/loadingFrameTime)%len(loadingFrames)]
	}
	return "\n  Searching connections " + string(frame)
}
