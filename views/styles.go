package views

import (
	"github.com/necrom4/sbb-tui/config"

	"github.com/charmbracelet/lipgloss"
)

const (
	// Layout dimensions
	borderSize     = 2
	hdrHeight      = 3
	rsltMrgn       = 1
	smplConnHeight = 9
	smplConnMrgn   = 3
	helpBarHeight  = 1

	stopsLineFixedWidth = (borderSize * 2) + (smplConnMrgn * 2) + (2+5)*2 + 6
	stopsLineMinWidth   = 10

	fullConnPaddH = 3
	fullConnPaddV = 1

	minTermWidth  = 80
	minTermHeight = 24
)

type styles struct {
	text           lipgloss.Style
	active         lipgloss.Style
	inactive       lipgloss.Style
	detailedResult lipgloss.Style
	dimmedBorder   lipgloss.Style
	helpKey        lipgloss.Style
	helpDesc       lipgloss.Style
	ghostText      lipgloss.Style
	warning        lipgloss.Style
	warningBold    lipgloss.Style
	vehicleIcon    lipgloss.Style
	vehicleModel   lipgloss.Style
	company        lipgloss.Style
	logo           lipgloss.Style
	bold           lipgloss.Style
}

func newStyles(theme config.Theme) styles {
	textColor := lipgloss.Color(theme.Text)
	return styles{
		text: lipgloss.NewStyle().
			Foreground(textColor),
		active: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.ActiveBorder)).
			Padding(0, 1),
		inactive: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.InactiveBorder)).
			Padding(0, 1),
		detailedResult: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.ActiveBorder)).
			Padding(fullConnPaddV, fullConnPaddH),
		dimmedBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.DimmedBorder)).
			Padding(0, rsltMrgn),
		helpKey: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(theme.KeysFg)).
			Background(lipgloss.Color(theme.KeysBg)).
			Padding(0, 1),
		helpDesc: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.GhostText)),
		ghostText: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.GhostText)),
		warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.WarningFlag)),
		warningBold: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.WarningFlag)).
			Bold(true),
		vehicleIcon: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.VehicleBg)).
			Foreground(lipgloss.Color(theme.VehicleFg)),
		vehicleModel: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.ModelBg)).
			Foreground(lipgloss.Color(theme.ModelFg)).
			Bold(true),
		company: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.CompanyBg)).
			Foreground(lipgloss.Color(theme.CompanyFg)),
		logo: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Logo)),
		bold: lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true),
	}
}
