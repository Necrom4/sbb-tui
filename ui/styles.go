package ui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/necrom4/sbb-tui/config"
)

const (
	// Layout dimensions
	borderSize          = 2
	headerHeight        = 3
	resultMargin        = 1
	simpleConnHeight    = 9
	simpleConnMargin    = 3
	helpBarHeight       = 1
	stopsLineFixedWidth = (borderSize * 2) + (simpleConnMargin * 2) + (2+5)*2 + 6
	stopsLineMinWidth   = 10
	detailPaddingH      = 3
	detailPaddingV      = 1
	minTermWidth        = 80
	minTermHeight       = 24
)

type styles struct {
	text           lipgloss.Style
	error          lipgloss.Style
	textMuted      lipgloss.Style
	active         lipgloss.Style
	inactive       lipgloss.Style
	detailedResult lipgloss.Style
	helpKey        lipgloss.Style
	helpDesc       lipgloss.Style
	warning        lipgloss.Style
	warningBold    lipgloss.Style
	vehicleIcon    lipgloss.Style
	vehicleModel   lipgloss.Style
	company        lipgloss.Style
	logo           lipgloss.Style
	bold           lipgloss.Style
}

func newStyles(theme config.Theme) styles {
	return styles{
		text: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Text)),
		error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Error)),
		textMuted: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.TextMuted)),
		active: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.BorderFocused)).
			Foreground(lipgloss.Color(theme.Text)).
			Padding(0, 1),
		inactive: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.BorderUnfocused)).
			Foreground(lipgloss.Color(theme.Text)).
			Padding(0, 1),
		detailedResult: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.BorderFocused)).
			Padding(detailPaddingV, detailPaddingH),
		helpKey: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(theme.BadgeKeyFg)).
			Background(lipgloss.Color(theme.BadgeKeyBg)).
			Padding(0, 1),
		helpDesc: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.TextMuted)),
		warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Warning)),
		warningBold: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Warning)).
			Bold(true),
		vehicleIcon: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.BadgeVehicleBg)).
			Foreground(lipgloss.Color(theme.BadgeVehicleFg)),
		vehicleModel: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.BadgeModelBg)).
			Foreground(lipgloss.Color(theme.BadgeBadgeModelFg)).
			Bold(true),
		company: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.BadgeCompanyBg)).
			Foreground(lipgloss.Color(theme.BadgeCompanyFg)),
		logo: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Logo)),
		bold: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Text)).
			Bold(true),
	}
}
