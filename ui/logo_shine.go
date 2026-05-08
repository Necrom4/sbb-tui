package ui

import (
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const animLogoShine = "logoShine"

const (
	logoShineDuration   = 800 * time.Millisecond
	logoShineRepeatGap  = 2 * time.Second // pause between successive sweeps on the start screen
	logoShineBandWidth  = 28.0            // diagonal half-width of the dark band
	logoShineDarkDelta  = 0.30            // max lightness reduction at the band's center (0..1)
	logoShinePaletteLen = 65              // pre-built color steps; higher = smoother
)

type logoShineRestartMsg struct{}

func logoShineRestartCmd() tea.Cmd {
	return tea.Tick(logoShineRepeatGap, func(time.Time) tea.Msg {
		return logoShineRestartMsg{}
	})
}

func (m appModel) renderLogo(logo string) string {
	if !m.animations {
		return m.styles.logo.Render(logo)
	}
	progress, active := m.anim.Progress(animLogoShine)
	if !active {
		return m.styles.logo.Render(logo)
	}
	return shineLogo(logo, m.styles.logoBase, progress)
}

// shineLogo paints each non-space rune with a color that depends on
// its diagonal distance from a band sweeping from upper-left to
// lower-right as `progress` goes from 0 to 1.
func shineLogo(logo string, base colorful.Color, progress float64) string {
	lines := strings.Split(logo, "\n")
	maxWidth := 0
	for _, ln := range lines {
		if w := lipgloss.Width(ln); w > maxWidth {
			maxWidth = w
		}
	}

	dMin := 0.0
	dMax := float64(maxWidth + len(lines) - 2)
	span := (dMax - dMin) + 2*logoShineBandWidth
	dCenter := dMin - logoShineBandWidth + progress*span

	palette := buildShinePalette(base)

	var b strings.Builder
	b.Grow(len(logo) * 6)
	for row, line := range lines {
		col := 0
		for _, r := range line {
			if r == ' ' {
				b.WriteRune(r)
				col++
				continue
			}
			d := float64(col + row)
			factor := shineFactor(d - dCenter)
			b.WriteString(palette.render(factor, r))
			col++
		}
		// Right-pad to the widest line so JoinVertical's center alignment
		// doesn't shift shorter rows and distort the logo.
		if pad := maxWidth - col; pad > 0 {
			b.WriteString(strings.Repeat(" ", pad))
		}
		if row < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// shineFactor maps signed diagonal distance from the band center to a darkness factor in [0, 1].
func shineFactor(signedDist float64) float64 {
	dist := math.Abs(signedDist)
	if dist >= logoShineBandWidth {
		return 0
	}
	return math.Cos(dist / logoShineBandWidth * math.Pi / 2)
}

type shinePalette struct {
	cells [logoShinePaletteLen]paletteCell
}

type paletteCell struct {
	prefix string
	suffix string
}

func buildShinePalette(base colorful.Color) shinePalette {
	var p shinePalette
	bh, bs, bl := base.Hsl()
	for i := 0; i < logoShinePaletteLen; i++ {
		f := float64(i) / float64(logoShinePaletteLen-1)
		l := bl - f*logoShineDarkDelta
		if l < 0 {
			l = 0
		}
		if l > 1 {
			l = 1
		}
		c := colorful.Hsl(bh, bs, l)
		hex := c.Clamped().Hex()
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
		const sentinel = "\x00"
		rendered := style.Render(sentinel)
		idx := strings.Index(rendered, sentinel)
		if idx < 0 {
			p.cells[i] = paletteCell{}
			continue
		}
		p.cells[i] = paletteCell{
			prefix: rendered[:idx],
			suffix: rendered[idx+len(sentinel):],
		}
	}
	return p
}

func (p shinePalette) render(factor float64, r rune) string {
	if factor < 0 {
		factor = 0
	} else if factor > 1 {
		factor = 1
	}
	idx := int(math.Round(factor * float64(logoShinePaletteLen-1)))
	cell := p.cells[idx]
	return cell.prefix + string(r) + cell.suffix
}
