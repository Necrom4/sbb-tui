package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const animationFPS = 30

var animationTick = time.Second / animationFPS

type animationTickMsg struct{}

type animation struct {
	start    time.Time
	duration time.Duration
	done     bool
}

// animator manages a set of named time-bounded animations driven by a single shared ticker.
type animator struct {
	anims map[string]*animation
}

func newAnimator() animator {
	return animator{anims: map[string]*animation{}}
}

// Start (re)starts a named animation and returns a Cmd to schedule the next frame.
func (a *animator) Start(name string, duration time.Duration) tea.Cmd {
	if a.anims == nil {
		a.anims = map[string]*animation{}
	}
	a.anims[name] = &animation{
		start:    time.Now(),
		duration: duration,
	}
	return animationTickCmd()
}

// Tick advances all animations. It returns the names of animations
// that just finished on this tick and a Cmd to schedule the next
// frame (nil when no animation is still running).
func (a *animator) Tick() (finished []string, next tea.Cmd) {
	anyRunning := false
	now := time.Now()
	for name, an := range a.anims {
		if an.done {
			continue
		}
		if now.Sub(an.start) >= an.duration {
			an.done = true
			finished = append(finished, name)
			continue
		}
		anyRunning = true
	}
	if anyRunning {
		next = animationTickCmd()
	}
	return finished, next
}

// Progress returns the [0,1] progress of a named animation and whether it is currently active.
func (a animator) Progress(name string) (float64, bool) {
	an, ok := a.anims[name]
	if !ok || an.done {
		return 1, false
	}
	elapsed := time.Since(an.start)
	if elapsed >= an.duration {
		return 1, false
	}
	if elapsed < 0 {
		return 0, true
	}
	return float64(elapsed) / float64(an.duration), true
}

func animationTickCmd() tea.Cmd {
	return tea.Tick(animationTick, func(time.Time) tea.Msg {
		return animationTickMsg{}
	})
}
