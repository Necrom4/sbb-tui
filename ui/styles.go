package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"

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

// Foreground hex used by the adaptive themeColor helper. Picked to be
// readable on either background extreme.
const (
	adaptiveLight = "#1A1A1A"
	adaptiveDark  = "#FFFFFF"
)

type styles struct {
	text            lipgloss.Style
	error           lipgloss.Style
	textMuted       lipgloss.Style
	active          lipgloss.Style
	inactive        lipgloss.Style
	detailedResult  lipgloss.Style
	helpKey         lipgloss.Style
	helpDesc        lipgloss.Style
	warning         lipgloss.Style
	warningBold     lipgloss.Style
	vehicleIcon     lipgloss.Style
	vehicleModel    lipgloss.Style
	company         lipgloss.Style
	logo            lipgloss.Style
	bold            lipgloss.Style
	logoBase        colorful.Color
	textMutedBase   colorful.Color
	background      colorful.Color
	backgroundKnown bool
}

func newStyles(theme config.Theme) styles {
	bg, bgKnown := detectBackground()

	return styles{
		text: lipgloss.NewStyle().
			Foreground(themeColor(theme.Text)),
		error: lipgloss.NewStyle().
			Foreground(themeColor(theme.Error)),
		textMuted: lipgloss.NewStyle().
			Foreground(themeColor(theme.TextMuted)),
		active: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(themeColor(theme.BorderFocused)).
			Foreground(themeColor(theme.Text)).
			Padding(0, 1),
		inactive: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(themeColor(theme.BorderUnfocused)).
			Foreground(themeColor(theme.Text)).
			Padding(0, 1),
		detailedResult: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(themeColor(theme.BorderFocused)).
			Padding(detailPaddingV, detailPaddingH),
		helpKey: lipgloss.NewStyle().
			Bold(true).
			Foreground(themeColor(theme.BadgeKeyFg)).
			Background(themeColor(theme.BadgeKeyBg)).
			Padding(0, 1),
		helpDesc: lipgloss.NewStyle().
			Foreground(themeColor(theme.TextMuted)),
		warning: lipgloss.NewStyle().
			Foreground(themeColor(theme.Warning)),
		warningBold: lipgloss.NewStyle().
			Foreground(themeColor(theme.Warning)).
			Bold(true),
		vehicleIcon: lipgloss.NewStyle().
			Background(themeColor(theme.BadgeVehicleBg)).
			Foreground(themeColor(theme.BadgeVehicleFg)),
		vehicleModel: lipgloss.NewStyle().
			Background(themeColor(theme.BadgeModelBg)).
			Foreground(themeColor(theme.BadgeBadgeModelFg)).
			Bold(true),
		company: lipgloss.NewStyle().
			Background(themeColor(theme.BadgeCompanyBg)).
			Foreground(themeColor(theme.BadgeCompanyFg)),
		logo: lipgloss.NewStyle().
			Foreground(themeColor(theme.Logo)),
		bold: lipgloss.NewStyle().
			Foreground(themeColor(theme.Text)).
			Bold(true),
		logoBase:        parseColor(theme.Logo),
		textMutedBase:   parseColor(theme.TextMuted),
		background:      bg,
		backgroundKnown: bgKnown,
	}
}

// detectBackground returns the terminal's background color and
// whether the OSC 11 query produced a usable value.
func detectBackground() (colorful.Color, bool) {
	raw := termenv.BackgroundColor()
	if _, isNo := raw.(termenv.NoColor); isNo {
		return colorful.Color{}, false
	}
	return termenv.ConvertToRGB(raw), true
}

// themeColor accepts a theme value (hex like "#RRGGBB" or the
// adaptive sentinels "white"/"black") and returns a lipgloss color
// that resolves correctly on both light and dark terminals.
func themeColor(s string) lipgloss.TerminalColor {
	switch s {
	case "white":
		return lipgloss.AdaptiveColor{Light: adaptiveLight, Dark: adaptiveDark}
	case "black":
		return lipgloss.AdaptiveColor{Light: adaptiveDark, Dark: adaptiveLight}
	}
	return lipgloss.Color(s)
}

var (
	colorBlack = colorful.Color{R: 0, G: 0, B: 0}
	colorWhite = colorful.Color{R: 1, G: 1, B: 1}
)

// parseColor converts a theme color string into the RGB color the
// terminal will actually display, expanding the "white"/"black"
// adaptive sentinels via the detected background.
func parseColor(s string) colorful.Color {
	if c, err := colorful.Hex(s); err == nil {
		return c
	}
	switch s {
	case "white":
		if termenv.HasDarkBackground() {
			return mustHex(adaptiveDark)
		}
		return mustHex(adaptiveLight)
	case "black":
		if termenv.HasDarkBackground() {
			return mustHex(adaptiveLight)
		}
		return mustHex(adaptiveDark)
	case "red":
		return colorful.Color{R: 1, G: 0, B: 0}
	case "green":
		return colorful.Color{R: 0, G: 1, B: 0}
	case "blue":
		return colorful.Color{R: 0, G: 0, B: 1}
	case "yellow":
		return colorful.Color{R: 1, G: 1, B: 0}
	case "cyan":
		return colorful.Color{R: 0, G: 1, B: 1}
	case "magenta":
		return colorful.Color{R: 1, G: 0, B: 1}
	}
	return colorWhite
}

func mustHex(s string) colorful.Color {
	c, _ := colorful.Hex(s)
	return c
}
