// Package views
package views

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/necrom4/sbb-tui/api"
	"github.com/necrom4/sbb-tui/config"
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

type focusable struct {
	kind  int
	id    string
	index int
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
	keyUPDW   string
	keyRight  string
	keyEsc    string
}

func newIconSet(noNerdFont bool) iconSet {
	icons := iconSet{
		// Shared symbols
		plt:   "Pl.",
		twrds: "→",

		filledDot: "●",
		hollowDot: "○",
		horzLine:  "─",
		vertLine:  "│",

		keyTab:   "⇥",
		keyEnter: "↵",
		keySpace: "␣",
		keyUpDw:  "↕",
		keyUPDW:  "⇧↕",
		keyRight: "→",
		keyEsc:   "⎋",
	}

	if noNerdFont {
		icons.arr = "↘"
		icons.dpt = "↗"
		icons.srch = "⌕"
		icons.swp = "⇋"
		icons.vhc = "×"
		icons.wlk = "Walk:"
		icons.prompt = "> "
	} else {
		icons.arr = "󰗔"
		icons.dpt = ""
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

const suggestDebounce = 300 * time.Millisecond

type suggestTickMsg struct {
	inputIndex int
	seq        int
}

type model struct {
	width, height int
	tabIndex      int
	resultIndex   int
	detailScrollY int
	headerOrder   []focusable
	inputs        []textinput.Model
	icons         iconSet
	styles        styles
	noNerdFont    bool
	isArrivalTime bool
	connections   []models.Connection
	loading       bool
	errorMsg      string
	searched      bool
	lastFromQuery string
	lastToQuery   string
	suggestSeq    [2]int
}

func InitialModel(cfg config.Config) model {
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
		styles:        newStyles(cfg.Theme),
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
			t.Prompt = m.icons.prompt
			t.CharLimit = 10
			t.Width = t.CharLimit
			t.ShowSuggestions = true
			t.CompletionStyle = m.styles.ghostText
			t.KeyMap.AcceptSuggestion = key.NewBinding(key.WithKeys("right"))
			if cfg.Date != "" {
				t.SetValue(cfg.Date)
			} else {
				t.SetValue(now.Format("02.01.2006"))
			}
		case 3:
			t.Prompt = m.icons.prompt
			t.CharLimit = 5
			t.Width = t.CharLimit
			t.ShowSuggestions = true
			t.CompletionStyle = m.styles.ghostText
			t.KeyMap.AcceptSuggestion = key.NewBinding(key.WithKeys("right"))
			if cfg.Time != "" {
				t.SetValue(cfg.Time)
			} else {
				t.SetValue(now.Format("15:04"))
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
		remaining := m.width - m.headerFixedWidth()
		inputWidth := max(remaining/2, 1)
		m.inputs[0].Width = inputWidth
		m.inputs[1].Width = inputWidth

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

	case suggestTickMsg:
		// Fetch if no new keystroke has occurred since
		if msg.seq == m.suggestSeq[msg.inputIndex] {
			query := m.inputs[msg.inputIndex].Value()
			return m, fetchSuggestionsCmd(msg.inputIndex, query)
		}
		return m, nil

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
			m.styles.warningBold.Render(msg))
	}

	header := m.renderHeader()
	results := lipgloss.JoinHorizontal(lipgloss.Top,
		m.styles.text.
			Height(m.resultsHeight()).
			Render(m.renderResults()),
		m.styles.text.
			Height(m.resultsHeight()).
			Render(m.renderDetailedResult()),
	)

	helpBar := m.renderHelpBar()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.styles.dimmedBorder.
			Width(m.contentWidth()).
			Height(m.resultsHeight()).
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
				if len(val) == 4 || len(val) == 7 {
					t.SetValue(val[:len(val)-2])
					return nil
				}
			}

			// date/time input: auto-insert `.`/`:` and block non existent values
			if len(s) == 1 && s >= "0" && s <= "9" {
				switch len(val) {
				case 0:
					if s > "3" {
						return nil
					}
				case 1:
					if val[0] == '0' && s == "0" {
						return nil
					}
					if val[0] == '3' && s > "1" {
						return nil
					}
				case 2:
					if s > "1" {
						return nil
					}
					t.SetValue(val + "." + s)
					t.SetCursor(len(val) + 2)
					return nil
				case 3:
				case 4:
					if val[3] == '0' && s == "0" {
						return nil
					}
					if val[3] == '1' && s > "2" {
						return nil
					}
				case 5:
					if s > "2" {
						return nil
					}
					t.SetValue(val + "." + s)
					t.SetCursor(len(val) + 2)
					return nil
				case 6, 7, 8, 9:
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

			if len(s) == 1 && s >= "0" && s <= "9" {
				switch len(val) {
				case 0:
					if s > "2" {
						return nil
					}
				case 1:
					if val == "2" && s > "3" {
						return nil
					}
				case 2:
					if s > "5" {
						return nil
					}
					t.SetValue(val + ":" + s)
					t.SetCursor(5)
					return nil
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

	// Debounce suggestion fetches for from/to inputs when value changes
	if fromVal := m.inputs[0].Value(); fromVal != m.lastFromQuery {
		m.lastFromQuery = fromVal
		if len(fromVal) >= 2 {
			m.suggestSeq[0]++
			seq := m.suggestSeq[0]
			cmds = append(cmds, tea.Tick(suggestDebounce, func(time.Time) tea.Msg {
				return suggestTickMsg{inputIndex: 0, seq: seq}
			}))
		} else {
			m.inputs[0].SetSuggestions(nil)
		}
	}
	if toVal := m.inputs[1].Value(); toVal != m.lastToQuery {
		m.lastToQuery = toVal
		if len(toVal) >= 2 {
			m.suggestSeq[1]++
			seq := m.suggestSeq[1]
			cmds = append(cmds, tea.Tick(suggestDebounce, func(time.Time) tea.Msg {
				return suggestTickMsg{inputIndex: 1, seq: seq}
			}))
		} else {
			m.inputs[1].SetSuggestions(nil)
		}
	}

	// Update date/time inputs' ghost completion
	m.inputs[2].SetSuggestions([]string{completeDate(m.inputs[2].Value())})
	m.inputs[3].SetSuggestions([]string{completeTime(m.inputs[3].Value())})

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

func completeDate(partial string) string {
	now := time.Now()
	full := now.Format("02.01.2006")
	if len(partial) < len(full) {
		return partial + full[len(partial):]
	}
	return partial
}

func completeTime(partial string) string {
	if len(partial) < 5 {
		full := partial + "00:00"[len(partial):]
		return full
	}
	return partial
}

// Convert swiss date format to API complient
func toAPIDate(swiss string) string {
	parts := strings.Split(swiss, ".")
	if len(parts) != 3 {
		return swiss
	}
	return parts[2] + "-" + parts[1] + "-" + parts[0]
}

func (m model) searchCmd() tea.Cmd {
	return func() tea.Msg {
		res, err := api.FetchConnections(
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			toAPIDate(completeDate(m.inputs[2].Value())),
			completeTime(m.inputs[3].Value()),
			m.isArrivalTime,
			m.maxVisibleConnections(),
		)
		return DataMsg{connections: res, err: err}
	}
}

func (m model) headerFixedWidth() int {
	width := 0
	for i, item := range m.headerOrder {
		if item.id == "from" || item.id == "to" {
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
		{m.icons.keyUPDW, "scroll"},
		{m.icons.keyRight, "complete"},
		{m.icons.keyEsc, "quit"},
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, m.styles.helpKey.Render(k.key)+" "+m.styles.helpDesc.Render(k.desc))
	}

	return " " + strings.Join(parts, "   ")
}

func (m model) renderHeaderItem(idx int) string {
	item := m.headerOrder[idx]
	style := m.styles.inactive
	if m.tabIndex == idx {
		style = m.styles.active
	}

	if item.kind == KindInput {
		input := m.inputs[item.index]
		view := input.View()
		if input.ShowSuggestions {
			// Clip text to prevent suggestion overflow.
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
		return "\n  " + m.styles.warning.Render(m.errorMsg)
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

	coloredLogo := m.styles.logo.Render(logo)

	text := m.styles.ghostText.Render("Enter stations above to see timetables")

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
	wrapped := m.styles.text.Width(innerWidth).Render(content)
	visLines := strings.Split(wrapped, "\n")

	// Scroll and clamp to the visible area.
	if len(visLines) > boxHeight {
		scrollY := min(m.detailScrollY, len(visLines)-boxHeight)
		visLines = visLines[scrollY : scrollY+boxHeight]
	}

	return m.styles.detailedResult.Width(width).Height(boxHeight).Render(strings.Join(visLines, "\n"))
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

	vehicleIcon := m.styles.vehicleIcon.Render(" " + m.icons.vhc + " ")
	vehicleModel := m.styles.vehicleModel.Render(section.Journey.Category + " " + section.Journey.Number)
	company := m.styles.company.Render(section.Journey.Operator)
	vehicleLine := fmt.Sprintf("%s  %s  %s %s %s", indent, m.icons.vertLine, vehicleIcon, vehicleModel, company)
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
			walkDuration = fmt.Sprintf("%d", dur/60)
		} else {
			depTime := section.Departure.Departure.Time
			arrTime := section.Arrival.Arrival.Time
			if !depTime.IsZero() && !arrTime.IsZero() {
				walkDuration = fmt.Sprintf("%d", int(arrTime.Sub(depTime).Minutes()))
			}
		}
		url := getGoogleMapsURL(section)

		// TODO: add `` icon and set that as clickable url link instead of the time
		walkDuration = utils.RenderLink(walkDuration, url)
	}

	walkLine := fmt.Sprintf("           %s %s'", m.icons.wlk, walkDuration)
	lines = append(lines, walkLine)

	return lines
}

func (m model) formatStationLine(timeStr string, delay int, symbol, station, platform string, width, timeCol, delayCol, symbolCol int, bold bool) string {
	textStyle := m.styles.text
	if bold {
		textStyle = m.styles.bold
	}

	timePart := textStyle.Render(timeStr)

	delayPart := ""
	if delay > 0 {
		delayStr := fmt.Sprintf("+%d'", delay)
		delayPart = m.styles.warningBold.Render(fmt.Sprintf("%*s", delayCol, delayStr))
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
	firstVehicle := -1
	lastVehicle := -1
	for x := range c.Sections {
		if c.Sections[x].Journey != nil {
			if firstVehicle == -1 {
				firstVehicle = x
			}
			lastVehicle = x
		}
	}

	style := m.styles.inactive.Width(width)
	if index == m.resultIndex {
		style = m.styles.active.Width(width)
	}

	if firstVehicle == -1 {
		return style.Render("\n  Connection details unavailable\n")
	}

	lineContentWidth := max(width-style.GetHorizontalFrameSize()-2, 0)

	vehicleIcon := m.styles.vehicleIcon.Render(" " + m.icons.vhc + " ")
	vehicleModel := m.styles.vehicleModel.Render(c.Sections[firstVehicle].Journey.Category + " " + c.Sections[firstVehicle].Journey.Number)
	company := m.styles.company.Render(c.Sections[firstVehicle].Journey.Operator)
	endStop := m.styles.text.Render(c.Sections[firstVehicle].Journey.To)

	dep := c.Sections[firstVehicle].Departure.Departure.Local().Format("15:04")
	arr := c.ToData.Arrival.Local().Format("15:04")
	departure := m.styles.bold.Render(dep)
	arrival := m.styles.bold.Render(arr)

	departureDelay := m.formatDelay(c.Sections[firstVehicle].Departure.Delay)
	arrivalDelay := m.formatDelay(c.Sections[lastVehicle].Arrival.Delay)

	timelinePrefix := ""
	if c.Sections[0].Walk != nil {
		walkMinutes := int(c.Sections[0].Arrival.Arrival.Sub(c.Sections[0].Departure.Departure).Minutes())
		if walkMinutes > 0 {
			timelinePrefix = m.icons.wlk + " " + m.styles.text.Render(fmt.Sprintf("%d'", walkMinutes)) + "  "
		}
	}

	timelineFixedWidth := lipgloss.Width(timelinePrefix) +
		lipgloss.Width(departure) +
		lipgloss.Width(departureDelay) + 2 +
		2 +
		lipgloss.Width(arrival) +
		lipgloss.Width(arrivalDelay)
	stopsLineWidth := max(lineContentWidth-timelineFixedWidth, stopsLineMinWidth)
	stopsLineRaw := m.renderStopsLine(c, stopsLineWidth)
	timelineWidth := timelineFixedWidth + lipgloss.Width(stopsLineRaw)
	if overflow := timelineWidth - lineContentWidth; overflow > 0 {
		stopsLineWidth = max(stopsLineWidth-overflow, stopsLineMinWidth)
		stopsLineRaw = m.renderStopsLine(c, stopsLineWidth)
	}
	stopsLine := m.styles.bold.Render(stopsLineRaw)

	platformInfo := ""
	platform := c.Sections[firstVehicle].Departure.Platform
	if platform == "" {
		platform = c.FromData.Platform
	}
	if platform != "" {
		platformInfo = m.icons.plt + " " + m.styles.text.Render(platform)
	}

	duration := m.styles.text.Render(formatDuration(c.Duration))

	bottomLinePadding := max(lineContentWidth-lipgloss.Width(platformInfo)-lipgloss.Width(duration), 1)

	content := fmt.Sprintf("\n  %s %s %s  %s\n\n  %s%s%s  %s  %s%s\n\n  %s%s%v\n",
		vehicleIcon,
		vehicleModel,
		company,
		endStop,
		timelinePrefix,
		departure,
		departureDelay,
		stopsLine,
		arrival,
		arrivalDelay,
		platformInfo,
		strings.Repeat(" ", bottomLinePadding),
		duration,
	)

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
		return hours + " h " + minutes + " min"
	}
	return minutes + " min"
}

func (m model) formatDelay(delay int) string {
	if delay > 0 {
		return m.styles.warningBold.Render(fmt.Sprintf(" +%d", delay))
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
