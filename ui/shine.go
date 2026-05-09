package ui

import (
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	animLogoShine = "logoShine"
	animTextShine = "textShine"
)

const (
	shineDuration      = 800 * time.Millisecond
	shineRepeatGap     = 2 * time.Second
	shineDelta         = 0.30 // max lightness shift at the band's center (0..1)
	logoShineBandWidth = 28.0
	textShineBandWidth = 8.0

	// shineLumaPivot picks shine polarity: bases brighter than this darken, the rest brighten.
	shineLumaPivot = 0.6
)

type shineRestartMsg struct{}

func shineRestartCmd() tea.Cmd {
	return tea.Tick(shineRepeatGap, func(time.Time) tea.Msg {
		return shineRestartMsg{}
	})
}

// startShineCycle restarts both shine animations together, keeping their progress in lockstep.
func (m *appModel) startShineCycle() tea.Cmd {
	return tea.Batch(
		m.anim.Start(animLogoShine, shineDuration),
		m.anim.Start(animTextShine, shineDuration),
	)
}

// onAnimationsFinished walks the start-screen animation chain in
// response to the named animations completing on this tick.
func (m *appModel) onAnimationsFinished(finished []string) []tea.Cmd {
	var cmds []tea.Cmd
	var shineRestarted bool
	for _, name := range finished {
		switch name {
		case animLogoBuild:
			cmds = append(cmds, m.anim.Start(animTaglineBuild, taglineBuildDuration))
		case animTaglineBuild:
			cmds = append(cmds, m.startShineCycle())
		case animLogoShine, animTextShine:
			if !shineRestarted {
				cmds = append(cmds, shineRestartCmd())
				shineRestarted = true
			}
		}
	}
	return cmds
}

// shineDirection is the axis along which the band sweeps.
type shineDirection int

const (
	shineDiagonal shineDirection = iota
	shineHorizontal
)

type shineOpts struct {
	base      colorful.Color
	progress  float64
	bandWidth float64
	direction shineDirection
}

// applyShine paints each non-space rune of `s` with a color whose
// brightness depends on its position relative to a moving band.
func applyShine(s string, opts shineOpts) string {
	lines, maxWidth := textBounds(s)

	dMax := float64(maxWidth - 1)
	if opts.direction == shineDiagonal {
		dMax = float64(maxWidth + len(lines) - 2)
	}
	dCenter := -opts.bandWidth + opts.progress*(dMax+2*opts.bandWidth)

	return renderGrid(s, newShinePalette(opts.base), func(row, col int) float64 {
		d := float64(col)
		if opts.direction == shineDiagonal {
			d = float64(col + row)
		}
		return shineFactor(d-dCenter, opts.bandWidth)
	})
}

// shineFactor maps signed distance from the band center to a band intensity in [0, 1].
func shineFactor(signedDist, bandWidth float64) float64 {
	dist := math.Abs(signedDist)
	if dist >= bandWidth {
		return 0
	}
	return math.Cos(dist / bandWidth * math.Pi / 2)
}

// newShinePalette builds the gradient used by the shine effect: from
// the base color toward the auto-picked neutral, scaled so the peak
// only travels `shineDelta` of the way.
func newShinePalette(base colorful.Color) palette {
	peak := base.BlendLab(shineTarget(base), shineDelta)
	return newPalette(base, peak)
}

func (m appModel) renderLogo(logo string) string {
	if !m.animations {
		return m.styles.logo.Render(logo)
	}
	if progress, active := m.anim.Progress(animLogoBuild); active {
		return renderLogoBuild(logo, m.styles.logoBase, progress)
	}
	progress, active := m.anim.Progress(animLogoShine)
	if !active {
		return m.styles.logo.Render(logo)
	}
	return applyShine(logo, shineOpts{
		base:      m.styles.logoBase,
		progress:  progress,
		bandWidth: logoShineBandWidth,
		direction: shineDiagonal,
	})
}

func (m appModel) renderStartTagline(text string) string {
	if !m.animations {
		return m.styles.textMuted.Render(text)
	}
	// Hide the tagline until its build animation has been registered,
	// so the logo build plays alone first.
	if !m.anim.Registered(animTaglineBuild) {
		return strings.Repeat(" ", lipgloss.Width(text))
	}
	if progress, active := m.anim.Progress(animTaglineBuild); active {
		return renderTaglineBuild(text, m.styles.textMutedBase, progress)
	}
	progress, active := m.anim.Progress(animTextShine)
	if !active {
		return m.styles.textMuted.Render(text)
	}
	return applyShine(text, shineOpts{
		base:      m.styles.textMutedBase,
		progress:  progress,
		bandWidth: textShineBandWidth,
		direction: shineHorizontal,
	})
}
