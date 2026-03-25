// Package views
package views

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/necrom4/sbb-tui/api"
	"github.com/necrom4/sbb-tui/models"
	"github.com/necrom4/sbb-tui/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	//go:embed sbb-logo.txt
	sbbLogo string

	//go:embed sbb-logo-nerdfont.txt
	sbbLogoNerdFont string
)

const (
	// Focusable item kinds
	KindInput int = iota
	KindButton
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

var (
	// Colors
	sbbWhite      = lipgloss.Color("#FFFFFF")
	sbbMidWhite   = lipgloss.Color("#F6F6F6")
	sbbDarkWhite  = lipgloss.Color("#DDDDDD")
	sbbGray       = lipgloss.Color("#888888")
	sbbMidGray    = lipgloss.Color("#484848")
	sbbDarkGray   = lipgloss.Color("#333333")
	sbbLightBlack = lipgloss.Color("#212121")
	sbbBlack      = lipgloss.Color("#141414")
	sbbRed        = lipgloss.Color("#D82E20")
	sbbMidRed     = lipgloss.Color("#B52C24")
	sbbDarkRed    = lipgloss.Color("#862010")
	sbbLightBlue  = lipgloss.Color("#315086")
	sbbBlue       = lipgloss.Color("#2E3279")
	sbbGreen      = lipgloss.Color("#3A7446")
)

var (
	// Styles
	noStyle = lipgloss.NewStyle()

	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(sbbRed).
			Padding(0, 1)

	blurredStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(sbbMidGray).
			Padding(0, 1)

	detailedResultStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(sbbRed).
				Padding(fullConnPaddV, fullConnPaddH)

	helpKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(sbbWhite).
			Background(sbbMidGray).
			Padding(0, 1)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(sbbGray)
)

type focusable struct {
	kind  int
	id    string
	index int
}

// Config holds CLI flag values to pre-fill the TUI form.
type Config struct {
	From          string
	To            string
	Date          string
	Time          string
	IsArrivalTime bool
	NoNerdFont    bool
}

type iconSet struct {
	// Mode-dependent (Nerd Font vs Unicode fallback)
	arr    string
	dpt    string
	plt    string
	srch   string
	swp    string
	vhc    string
	wlk    string
	prompt string

	// Mode-invariant
	twrds     string
	filledDot string
	hollowDot string
	horzLine  string
	vertLine  string
	keyTab    string
	keyEnter  string
	keySpace  string
	keyUpDw   string
	keyRight  string
	keyEsc    string
}

func newIconSet(noNerdFont bool) iconSet {
	icons := iconSet{
		// Shared symbols
		twrds:     "→",
		filledDot: "●",
		hollowDot: "○",
		horzLine:  "─",
		vertLine:  "│",
		keyTab:    "⇥",
		keyEnter:  "↵",
		keySpace:  "␣",
		keyUpDw:   "↕",
		keyRight:  "→",
		keyEsc:    "⎋",
	}

	if noNerdFont {
		icons.arr = "↘"
		icons.dpt = "↗"
		icons.plt = "Pl."
		icons.srch = "⌕"
		icons.swp = "⇋"
		icons.vhc = "×"
		icons.wlk = "Walk:"
		icons.prompt = "> "
	} else {
		icons.arr = "󰗔"
		icons.dpt = ""
		icons.plt = "󱀓"
		icons.srch = ""
		icons.swp = ""
		icons.vhc = ""
		icons.wlk = ""
		icons.prompt = " "
	}

	return icons
}

type DataMsg struct {
	connections []models.Connection
	err         error
}

type SuggestionsMsg struct {
	inputIndex int
	names      []string
	err        error
}

type model struct {
	width, height int
	tabIndex      int
	resultIndex   int
	detailScrollY int
	headerOrder   []focusable
	inputs        []textinput.Model
	icons         iconSet
	noNerdFont    bool
	isArrivalTime bool
	connections   []models.Connection
	loading       bool
	errorMsg      string
	searched      bool
	lastFromQuery string
	lastToQuery   string
}

func InitialModel(cfg Config) model {
	// Define input prompts
	m := model{
		headerOrder: []focusable{
			{KindInput, "from", 0},
			{KindInput, "to", 1},
			{KindButton, "swap", -1},
			{KindButton, "isArrivalTime", -1},
			{KindInput, "date", 2},
			{KindInput, "time", 3},
			{KindButton, "search", -1},
		},
		inputs:        make([]textinput.Model, 4),
		icons:         newIconSet(cfg.NoNerdFont),
		noNerdFont:    cfg.NoNerdFont,
		isArrivalTime: cfg.IsArrivalTime,
	}

	now := time.Now()

	for i := range m.inputs {
		t := textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "From"
			t.Prompt = m.icons.prompt
			t.ShowSuggestions = true
			t.KeyMap.AcceptSuggestion = key.NewBinding(key.WithKeys("right"))
			if cfg.From != "" {
				t.SetValue(cfg.From)
			}
			t.Focus()
		case 1:
			t.Placeholder = "To"
			t.Prompt = m.icons.prompt
			t.ShowSuggestions = true
			t.KeyMap.AcceptSuggestion = key.NewBinding(key.WithKeys("right"))
			if cfg.To != "" {
				t.SetValue(cfg.To)
			}
		case 2:
			t.Placeholder = now.Format("2006-01-02")
			t.Prompt = m.icons.prompt
			t.CharLimit = 10
			t.Width = t.CharLimit - 1
			if cfg.Date != "" {
				t.SetValue(cfg.Date)
			}
		case 3:
			t.Placeholder = now.Format("15:04")
			t.Prompt = m.icons.prompt
			t.CharLimit = 5
			t.Width = t.CharLimit - 1
			if cfg.Time != "" {
				t.SetValue(cfg.Time)
			}
		}
		m.inputs[i] = t
	}
	return m
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Define keymaps
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// Set initial input widths based on fixed header elements, then
		// measure the actual rendered header and correct for any difference.
		remaining := m.width - m.headerFixedWidth()
		inputWidth := max(remaining/2, 1)
		m.inputs[0].Width = inputWidth
		m.inputs[1].Width = inputWidth
		if gap := m.width - lipgloss.Width(m.renderHeader()); gap > 0 {
			m.inputs[0].Width += (gap + 1) / 2
			m.inputs[1].Width += gap / 2
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "q":
			active := m.headerOrder[m.tabIndex]
			if active.kind == KindButton {
				return m, tea.Quit
			}

		case "enter":
			if err := m.validateInputs(); err != "" {
				m.errorMsg = err
				m.connections = nil
				m.searched = false
				m.resultIndex = 0
				return m, nil
			}
			m.loading = true
			m.connections = nil
			m.errorMsg = ""
			m.searched = true
			return m, m.searchCmd()

		case " ":
			active := m.headerOrder[m.tabIndex]
			switch active.id {
			case "swap":
				tmp := m.inputs[0].Value()
				m.inputs[0].SetValue(m.inputs[1].Value())
				m.inputs[1].SetValue(tmp)
			case "isArrivalTime":
				m.isArrivalTime = !m.isArrivalTime
			case "search":
				if err := m.validateInputs(); err != "" {
					m.errorMsg = err
					m.connections = nil
					m.searched = false
					m.resultIndex = 0
					return m, nil
				}
				m.loading = true
				m.connections = nil
				m.errorMsg = ""
				m.searched = true
				return m, m.searchCmd()
			}

		case "tab", "shift+tab":
			if msg.String() == "shift+tab" {
				m.tabIndex--
			} else {
				m.tabIndex++
			}

			if m.tabIndex >= len(m.headerOrder) {
				m.tabIndex = 0
			}
			if m.tabIndex < 0 {
				m.tabIndex = len(m.headerOrder) - 1
			}

			var cmds []tea.Cmd
			for _, item := range m.headerOrder {
				if item.kind == KindInput {
					if item.index == m.headerOrder[m.tabIndex].index {
						cmds = append(cmds, m.inputs[item.index].Focus())
					} else {
						m.inputs[item.index].Blur()
					}
				}
			}
			return m, tea.Batch(cmds...)

		case "up":
			if len(m.connections) > 0 && m.resultIndex > 0 {
				m.resultIndex--
				m.detailScrollY = 0
			}
		case "down":
			if len(m.connections) > 0 && m.resultIndex < len(m.connections)-1 {
				m.resultIndex++
				m.detailScrollY = 0
			}
		case "shift+up":
			if m.detailScrollY > 0 {
				m.detailScrollY--
			}
		case "shift+down":
			m.detailScrollY++
		}

	case SuggestionsMsg:
		if msg.err == nil {
			m.inputs[msg.inputIndex].SetSuggestions(msg.names)
		}
		return m, nil

	case DataMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = "Failed to fetch connections. Check your internet connection."
			return m, nil
		}
		m.connections = msg.connections
		m.resultIndex = 0
		m.detailScrollY = 0
		if len(m.connections) == 0 {
			m.errorMsg = "No connections found for the specified route."
		}
		return m, nil
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width < minTermWidth || m.height < minTermHeight {
		minW := minTermWidth
		msg := fmt.Sprintf("Terminal too small (%dx%d)\nMinimum size: %dx%d", m.width, m.height, minW, minTermHeight)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			noStyle.Foreground(sbbRed).Bold(true).Render(msg))
	}

	header := m.renderHeader()
	results := lipgloss.JoinHorizontal(lipgloss.Top,
		noStyle.
			Height(m.resultsHeight()).
			Render(m.renderResults()),
		noStyle.
			Height(m.resultsHeight()).
			Render(m.renderDetailedResult()),
	)

	helpBar := m.renderHelpBar()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		noStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(sbbDarkRed).
			Width(m.contentWidth()).
			Height(m.resultsHeight()).
			Padding(0, rsltMrgn).
			Render(results),
		helpBar,
	)
}

func (m model) contentWidth() int {
	return max(m.width-borderSize, 0)
}

func (m model) resultsHeight() int {
	return max(m.height-hdrHeight-borderSize-helpBarHeight, 0)
}

func (m model) maxVisibleConnections() int {
	return max(m.resultsHeight()/smplConnHeight, 1)
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check key input in input fields
		switch m.headerOrder[m.tabIndex].id {
		case "date":
			t := &m.inputs[2]
			s := msg.String()
			val := t.Value()

			if msg.Type == tea.KeyBackspace {
				if len(val) == 6 || len(val) == 9 {
					t.SetValue(val[:len(val)-2])
					return nil
				}
			}

			if len(s) == 1 && s >= "0" && s <= "9" {
				switch len(val) {
				// Year: YYYY (positions 0-3)
				case 0:
					if s > "2" {
						return nil
					}
				case 1, 2, 3:
				// 1st month digit: auto-insert "-", validate (0 or 1)
				case 4:
					if s > "1" {
						return nil
					}
					t.SetValue(val + "-" + s)
					t.SetCursor(len(val) + 2)
					return nil
				case 5:
				// 2nd month digit: val = "YYYY-M#" where val[5] is 1st month digit
				case 6:
					if val[5] == '0' && s == "0" {
						return nil // block month 00
					}
					if val[5] == '1' && s > "2" {
						return nil // block months 13-19
					}
				// 1st day digit: auto-insert "-", validate (0-3)
				case 7:
					if s > "3" {
						return nil
					}
					t.SetValue(val + "-" + s)
					t.SetCursor(len(val) + 2)
					return nil
				case 8:
				// 2nd day digit: val = "YYYY-MM-D#" where val[8] is 1st day digit
				case 9:
					if val[8] == '0' && s == "0" {
						return nil // block day 00
					}
					if val[8] == '3' && s > "1" {
						return nil // block days 32-39
					}
				default:
					return nil
				}
			} else if msg.Type == tea.KeyRunes {
				return nil
			}

		case "time":
			t := &m.inputs[3]
			s := msg.String()
			val := t.Value()

			if msg.Type == tea.KeyBackspace && len(val) == 4 {
				t.SetValue(val[:2])
				return nil
			}

			// Only process numeric runes for the following logic
			if len(s) == 1 && s >= "0" && s <= "9" {
				switch len(val) {
				// Hours: HH
				case 0:
					if s > "2" {
						return nil
					}
				case 1:
					if val == "2" && s > "3" {
						return nil
					}
				// Minutes: auto-insert ":", validate 1st minute digit (0-5)
				case 2:
					if s > "5" {
						return nil
					}
					t.SetValue(val + ":" + s)
					t.SetCursor(5)
					return nil
				// 2nd minute digit (0-9, no restriction)
				case 3, 4:
				default:
					return nil
				}
			} else if msg.Type == tea.KeyRunes {
				return nil
			}
		}
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// Trigger suggestion fetches for from/to inputs when value changes
	if fromVal := m.inputs[0].Value(); fromVal != m.lastFromQuery {
		m.lastFromQuery = fromVal
		if len(fromVal) >= 2 {
			cmds = append(cmds, fetchSuggestionsCmd(0, fromVal))
		} else {
			m.inputs[0].SetSuggestions(nil)
		}
	}
	if toVal := m.inputs[1].Value(); toVal != m.lastToQuery {
		m.lastToQuery = toVal
		if len(toVal) >= 2 {
			cmds = append(cmds, fetchSuggestionsCmd(1, toVal))
		} else {
			m.inputs[1].SetSuggestions(nil)
		}
	}

	return tea.Batch(cmds...)
}

func (m model) validateInputs() string {
	if m.inputs[0].Value() == "" {
		return "Please enter a departure station."
	}
	if m.inputs[1].Value() == "" {
		return "Please enter an arrival station."
	}
	return ""
}

func fetchSuggestionsCmd(inputIndex int, query string) tea.Cmd {
	return func() tea.Msg {
		names, err := api.FetchLocations(query)
		return SuggestionsMsg{inputIndex: inputIndex, names: names, err: err}
	}
}

func (m model) searchCmd() tea.Cmd {
	maxConnections := m.maxVisibleConnections()
	return func() tea.Msg {
		res, err := api.FetchConnections(
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			m.inputs[2].Value(),
			m.inputs[3].Value(),
			m.isArrivalTime,
			maxConnections,
		)
		return DataMsg{connections: res, err: err}
	}
}

func (m model) headerFixedWidth() int {
	width := 0
	for i, item := range m.headerOrder {
		if item.kind == KindInput && m.inputs[item.index].ShowSuggestions {
			// From/To: only count the per-item overhead (border + padding + prompt).
			width += borderSize + 2 + lipgloss.Width(m.inputs[item.index].Prompt)
			continue
		}
		width += lipgloss.Width(m.renderHeaderItem(i))
	}
	return width
}

func (m model) renderHeader() string {
	var headerItems []string
	for i := range m.headerOrder {
		headerItems = append(headerItems, m.renderHeaderItem(i))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, headerItems...)
}

func (m model) renderHelpBar() string {
	keys := []struct{ key, desc string }{
		{m.icons.keyTab, "navigate"},
		{m.icons.keyEnter, "search"},
		{m.icons.keySpace, "toggle"},
		{m.icons.keyUpDw, "results"},
		{m.icons.keyRight, "complete"},
		{m.icons.keyEsc, "quit"},
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, helpKeyStyle.Render(k.key)+" "+helpDescStyle.Render(k.desc))
	}

	return " " + strings.Join(parts, "   ")
}

func (m model) renderHeaderItem(idx int) string {
	item := m.headerOrder[idx]
	style := blurredStyle
	if m.tabIndex == idx {
		style = focusedStyle
	}

	if item.kind == KindInput {
		input := m.inputs[item.index]
		view := input.View()
		if input.ShowSuggestions {
			// Clip text to prevent line wrapping.
			maxView := lipgloss.Width(input.Prompt) + input.Width
			view = ansi.Truncate(view, maxView, "")
		}
		return style.Render(view)
	}

	icon := " "
	switch item.id {
	case "swap":
		icon = m.icons.swp
	case "isArrivalTime":
		if m.isArrivalTime {
			icon = m.icons.arr
		} else {
			icon = m.icons.dpt
		}
	case "search":
		icon = m.icons.srch
	}
	return style.Render(icon)
}

func (m model) resultBoxWidth() int {
	return max((m.width-smplConnMrgn)/2, rsltMrgn+stopsLineMinWidth+stopsLineFixedWidth)
}

func (m model) renderResults() string {
	if m.loading {
		return "\n  Searching connections..."
	}

	if m.errorMsg != "" {
		return "\n  " + noStyle.Foreground(sbbRed).Render(m.errorMsg)
	}

	if len(m.connections) == 0 {
		if m.searched {
			return "\n  No connections found."
		}
		return m.renderStartScreen()
	}

	var boxes []string
	boxWidth := m.resultBoxWidth()

	for i, c := range m.connections {
		boxes = append(boxes, m.renderSimpleConnection(c, i, boxWidth))
	}

	return lipgloss.JoinVertical(lipgloss.Left, boxes...)
}

func (m model) renderStartScreen() string {
	logo := sbbLogoNerdFont
	if m.noNerdFont {
		logo = sbbLogo
	}
	logo = strings.TrimRight(logo, "\n")

	coloredLogo := noStyle.Foreground(sbbWhite).Render(logo)

	text := noStyle.Foreground(sbbGray).Render("Enter stations above to see timetables")

	block := lipgloss.JoinVertical(lipgloss.Center, text, "", coloredLogo)

	width := max(m.contentWidth()-borderSize-(rsltMrgn*2), 0)
	height := m.resultsHeight()

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, block)
}

func (m model) renderDetailedResult() string {
	if len(m.connections) == 0 {
		return ""
	}

	boxWidth := max(m.width-borderSize*4-m.resultBoxWidth(), 0)
	return m.renderFullConnection(m.connections[m.resultIndex], boxWidth)
}

func (m model) renderFullConnection(c models.Connection, width int) string {
	var lines []string
	innerWidth := max(width-borderSize-(fullConnPaddH*2), 0)

	for i, section := range c.Sections {
		isFirst := i == 0
		isLast := i == len(c.Sections)-1

		if section.Walk != nil {
			lines = append(lines, m.renderWalkSection(section)...)
		} else if section.Journey != nil {
			lines = append(lines, m.renderJourneySection(section, innerWidth, isFirst, isLast)...)
		}

		if !isLast {
			lines = append(lines, "", "")
		}
	}

	boxHeight := max(m.resultsHeight()-borderSize-(fullConnPaddV*2), 0)

	// Wrap and split into visual lines for scrolling.
	content := strings.Join(lines, "\n")
	wrapped := noStyle.Width(innerWidth).Render(content)
	visLines := strings.Split(wrapped, "\n")

	// Scroll and clamp to the visible area.
	if len(visLines) > boxHeight {
		scrollY := min(m.detailScrollY, len(visLines)-boxHeight)
		visLines = visLines[scrollY : scrollY+boxHeight]
	}

	return detailedResultStyle.Width(width).Height(boxHeight).Render(strings.Join(visLines, "\n"))
}

func (m model) renderJourneySection(section models.Section, width int, isFirst, isLast bool) []string {
	var lines []string

	const timeCol = 5
	const delayCol = 4
	const symbolCol = 5

	depTime := section.Departure.Departure.Local().Format("15:04")
	depDelay := section.Departure.Delay
	depStation := section.Departure.Station.Name
	depPlatform := section.Departure.Platform

	depDot := m.icons.hollowDot
	if isFirst {
		depDot = m.icons.filledDot
	}

	depLine := m.formatStationLine(depTime, depDelay, depDot, depStation, depPlatform, width, timeCol, delayCol, symbolCol, true)
	lines = append(lines, depLine)

	indent := strings.Repeat(" ", timeCol+delayCol)
	spacingLine := fmt.Sprintf("%s  %s", indent, m.icons.vertLine)
	lines = append(lines, spacingLine)

	vehicleIcon := noStyle.Background(sbbBlue).Foreground(sbbWhite).Render(" " + m.icons.vhc + " ")
	vehicleCategory := noStyle.Background(sbbRed).Foreground(sbbWhite).Bold(true).
		Render(section.Journey.Category + " " + section.Journey.Number)
	company := noStyle.Background(sbbWhite).Foreground(sbbBlack).
		Render(section.Journey.Operator)
	vehicleLine := fmt.Sprintf("%s  %s  %s %s %s", indent, m.icons.vertLine, vehicleIcon, vehicleCategory, company)
	lines = append(lines, vehicleLine)

	destLine := fmt.Sprintf("%s  %s   %s %s", indent, m.icons.vertLine, m.icons.twrds, section.Journey.To)
	lines = append(lines, destLine)

	lines = append(lines, spacingLine)

	arrTime := section.Arrival.Arrival.Local().Format("15:04")
	arrDelay := section.Arrival.Delay
	arrStation := section.Arrival.Station.Name
	arrPlatform := section.Arrival.Platform

	arrSymbol := m.icons.vertLine
	if isLast {
		arrSymbol = m.icons.filledDot
	}

	arrLine := m.formatStationLine(arrTime, arrDelay, arrSymbol, arrStation, arrPlatform, width, timeCol, delayCol, symbolCol, false)
	lines = append(lines, arrLine)

	return lines
}

func getGoogleMapsURL(s models.Section) string {
	dep := s.Departure.Station.Coordinate
	arr := s.Arrival.Station.Coordinate
	return fmt.Sprintf("https://www.google.com/maps/dir/?api=1&origin=%f,%f&destination=%f,%f&travelmode=walking",
		dep.X, dep.Y, arr.X, arr.Y)
}

func (m model) renderWalkSection(section models.Section) []string {
	var lines []string

	walkDuration := ""
	if section.Walk != nil {
		dur := section.Walk.Duration
		if dur > 0 {
			walkDuration = fmt.Sprintf("%d min", dur/60)
		} else {
			depTime := section.Departure.Departure.Time
			arrTime := section.Arrival.Arrival.Time
			if !depTime.IsZero() && !arrTime.IsZero() {
				walkDuration = fmt.Sprintf("%d min", int(arrTime.Sub(depTime).Minutes()))
			}
		}
		url := getGoogleMapsURL(section)

		walkDuration = utils.RenderLink(walkDuration, url)
	}

	walkLine := fmt.Sprintf("           %s %s", m.icons.wlk, walkDuration)
	lines = append(lines, walkLine)

	return lines
}

func (m model) formatStationLine(timeStr string, delay int, symbol, station, platform string, width, timeCol, delayCol, symbolCol int, bold bool) string {
	textStyle := noStyle
	if bold {
		textStyle = noStyle.Bold(true)
	}

	timePart := textStyle.Render(timeStr)

	delayPart := ""
	if delay > 0 {
		delayStr := fmt.Sprintf("+%d", delay)
		delayPart = noStyle.Foreground(sbbRed).Bold(true).Render(fmt.Sprintf("%*s", delayCol, delayStr))
	} else {
		delayPart = strings.Repeat(" ", delayCol)
	}

	symbolPart := fmt.Sprintf("  %s  ", symbol)

	platformPart := ""
	platformVisibleLen := 0
	if platform != "" {
		platformPart = textStyle.Render(fmt.Sprintf("%s %s", m.icons.plt, platform))
		platformVisibleLen = len(platform) + len(m.icons.plt) + 1
	}

	fixedWidth := timeCol + delayCol + symbolCol + platformVisibleLen
	availableForStation := max(width-fixedWidth-1, 5)

	truncatedStation := truncateString(station, availableForStation)
	stationPart := textStyle.Render(truncatedStation)

	stationLen := len(truncatedStation)
	padding := max(availableForStation-stationLen, 1)

	if platformPart != "" {
		return fmt.Sprintf("%s%s%s%s%s%s",
			timePart, delayPart, symbolPart, stationPart, strings.Repeat(" ", padding), platformPart)
	}
	return fmt.Sprintf("%s%s%s%s", timePart, delayPart, symbolPart, stationPart)
}

func truncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if maxLen <= 3 {
		return s[:min(len(s), maxLen)]
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func (m model) renderSimpleConnection(c models.Connection, index int, width int) string {
	firstVehicle := 0
	for x := range c.Sections {
		if c.Sections[x].Journey != nil {
			firstVehicle = x
			break
		}
	}

	vehicleIcon := noStyle.Background(sbbBlue).Foreground(sbbWhite).Render(" " + m.icons.vhc + " ")
	vehicleCategory := noStyle.Background(sbbRed).Foreground(sbbWhite).Bold(true).
		Render(c.Sections[firstVehicle].Journey.Category + " " + c.Sections[firstVehicle].Journey.Number)
	company := noStyle.Background(sbbWhite).Foreground(sbbBlack).
		Render(c.Sections[firstVehicle].Journey.Operator)
	endStop := noStyle.Render(c.Sections[firstVehicle].Journey.To)

	dep := c.FromData.Departure.Local().Format("15:04")
	arr := c.ToData.Arrival.Local().Format("15:04")
	departure := noStyle.Bold(true).Render(dep)
	arrival := noStyle.Bold(true).Render(arr)

	departureDelay := formatDelay(c.Sections[firstVehicle].Departure.Delay)
	arrivalDelay := formatDelay(c.Sections[firstVehicle].Arrival.Delay)

	stopsLineWidth := max(width-stopsLineFixedWidth, stopsLineMinWidth)
	stopsLine := noStyle.Bold(true).Render(m.renderStopsLine(c, stopsLineWidth))

	platformOrWalk := ""
	if len(c.FromData.Platform) > 0 {
		platformOrWalk = m.icons.plt + " " + noStyle.Render(c.FromData.Platform)
	} else if c.Sections[0].Walk != nil {
		platformOrWalk = m.icons.wlk + " " + noStyle.Render(
			fmt.Sprintf("%vm", c.Sections[0].Arrival.Arrival.Sub(c.Sections[0].Departure.Departure).Minutes()),
		)
	}

	duration := noStyle.Render(formatDuration(c.Duration))

	bottomLinePadding := max(width-(borderSize*2+smplConnMrgn*2+smplConnMrgn*2+3+5), 1)

	content := fmt.Sprintf("\n  %s %s %s  %s\n\n  %s%s  %s  %s%s\n\n  %s%s%v\n",
		vehicleIcon,
		vehicleCategory,
		company,
		endStop,
		departure,
		departureDelay,
		stopsLine,
		arrival,
		arrivalDelay,
		platformOrWalk,
		strings.Repeat(" ", bottomLinePadding),
		duration,
	)

	style := blurredStyle.Width(width)
	if index == m.resultIndex {
		style = focusedStyle.Width(width)
	}

	return style.Render(content)
}

// 00d01:15:00" -> "1h 15m" or "15 min".
func formatDuration(duration string) string {
	parts := strings.Split(duration, ":")
	if len(parts) < 2 {
		return duration
	}

	minutes := parts[1]
	if len(parts[0]) > 3 && parts[0][3:] != "00" {
		hours := parts[0][3:]
		return hours + "h " + minutes + "m"
	}
	return minutes + "min"
}

func formatDelay(delay int) string {
	if delay > 0 {
		return noStyle.Foreground(sbbRed).Bold(true).Render(fmt.Sprintf(" +%d", delay))
	}
	return ""
}

func (m model) renderStopsLine(c models.Connection, totalWidth int) string {
	if len(c.Sections) == 0 {
		return m.icons.filledDot + m.icons.horzLine + m.icons.horzLine + m.icons.filledDot
	}

	var sectionDurations []time.Duration
	var totalSectionDuration time.Duration
	for _, s := range c.Sections {
		// Skip walking sections
		if s.Journey == nil {
			continue
		}
		dep := s.Departure.Departure.Time
		arr := s.Arrival.Arrival.Time
		if !dep.IsZero() && !arr.IsZero() {
			dur := arr.Sub(dep)
			sectionDurations = append(sectionDurations, dur)
			totalSectionDuration += dur
		}
	}

	if totalSectionDuration == 0 || len(sectionDurations) == 0 {
		// Fallback to equal distribution
		return m.icons.filledDot + strings.Repeat(m.icons.horzLine+m.icons.horzLine+m.icons.hollowDot, c.Transfers) + m.icons.horzLine + m.icons.horzLine + m.icons.filledDot
	}

	var sb strings.Builder
	sb.WriteString(m.icons.filledDot)

	usedChars := 0
	for i, secDur := range sectionDurations {
		var lineChars int
		if i == len(sectionDurations)-1 {
			// Last section gets remaining chars to avoid rounding errors
			lineChars = totalWidth - usedChars
		} else {
			proportion := float64(secDur) / float64(totalSectionDuration)
			lineChars = int(proportion*float64(totalWidth) + 0.5)
		}
		lineChars = max(lineChars, 1)
		usedChars += lineChars

		sb.WriteString(strings.Repeat(m.icons.horzLine, lineChars))
		if i < len(sectionDurations)-1 {
			sb.WriteString(m.icons.hollowDot)
		} else {
			sb.WriteString(m.icons.filledDot)
		}
	}

	return sb.String()
}
