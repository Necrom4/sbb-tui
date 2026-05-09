package ui

import (
	"math"
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

// animator runs a set of named time-bounded animations off a single shared ticker.
type animator struct {
	anims map[string]*animation
}

func newAnimator() animator {
	return animator{anims: map[string]*animation{}}
}

// Start (re)starts a named animation and returns the Cmd that schedules the next frame.
func (a *animator) Start(name string, duration time.Duration) tea.Cmd {
	if a.anims == nil {
		a.anims = map[string]*animation{}
	}
	a.anims[name] = &animation{start: time.Now(), duration: duration}
	return animationTickCmd()
}

// StartIndefinite starts an animation that never auto-finishes; callers must Stop it.
func (a *animator) StartIndefinite(name string) tea.Cmd {
	return a.Start(name, math.MaxInt64)
}

// Stop marks the named animation as finished without firing a completion event.
func (a *animator) Stop(name string) {
	if an, ok := a.anims[name]; ok {
		an.done = true
	}
}

// Registered reports whether the animation has ever been started.
func (a animator) Registered(name string) bool {
	_, ok := a.anims[name]
	return ok
}

// Elapsed returns how long the named animation has been running, and whether it still is.
func (a animator) Elapsed(name string) (time.Duration, bool) {
	an, ok := a.anims[name]
	if !ok || an.done {
		return 0, false
	}
	return time.Since(an.start), true
}

// Tick advances every animation and returns the names that just finished
// plus the Cmd for the next frame (nil once nothing is running).
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

// Progress returns the [0,1] progress of the named animation and whether it is still active.
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
