package ui

import (
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/necrom4/sbb-tui/config"
	"github.com/necrom4/sbb-tui/model"
	"github.com/necrom4/sbb-tui/util"
)

const (
	kindInput int = iota
	kindButton
)

type focusable struct {
	kind  int
	id    string
	index int
}

// User-facing error messages shown in the TUI.
var (
	errNoConnections       = errors.New("no connections found for the specified route")
	errMissingDeparture    = errors.New("please enter a departure station")
	errMissingArrival      = errors.New("please enter an arrival station")
	errConnectionMalformed = errors.New("connection details unavailable")
)

type dataMsg struct {
	connections []model.Connection
	err         error
}

type suggestionsMsg struct {
	inputIndex int
	names      []string
	err        error
}

const suggestDebounce = 300 * time.Millisecond

type suggestTickMsg struct {
	inputIndex int
	seq        int
}

type versionCheckMsg struct {
	newerVersion string
}

type appModel struct {
	width          int
	height         int
	tabIndex       int
	resultIndex    int
	detailScrollY  int
	headerOrder    []focusable
	inputs         []textinput.Model
	icons          iconSet
	styles         styles
	nerdFont       bool
	isArrivalTime  bool
	connections    []model.Connection
	loading        bool
	errorMsg       error
	searched       bool
	lastFromQuery  string
	lastToQuery    string
	suggestSeq     [2]int
	currentVersion string
	newerVersion   string
}

// NewModel creates the initial Bubbletea model from the application config.
func NewModel(cfg config.Config) appModel {
	m := appModel{
		headerOrder: []focusable{
			{kindInput, "from", 0},
			{kindInput, "to", 1},
			{kindButton, "swap", -1},
			{kindButton, "isArrivalTime", -1},
			{kindInput, "date", 2},
			{kindInput, "time", 3},
			{kindButton, "search", -1},
		},
		inputs:         make([]textinput.Model, 4),
		icons:          newIconSet(cfg.NerdFont),
		styles:         newStyles(cfg.Theme),
		nerdFont:       cfg.NerdFont,
		isArrivalTime:  cfg.IsArrivalTime,
		currentVersion: cfg.CurrentVersion,
	}

	now := time.Now()

	for i := range m.inputs {
		t := textinput.New()
		t.CharLimit = 32

		t.TextStyle = m.styles.text
		t.PromptStyle = m.styles.text
		t.PlaceholderStyle = m.styles.textMuted
		t.Cursor.Style = m.styles.active
		t.CompletionStyle = m.styles.textMuted
		t.Prompt = m.icons.prompt
		t.ShowSuggestions = true

		switch i {
		case 0:
			t.Placeholder = "From"
			t.KeyMap.AcceptSuggestion = key.NewBinding(key.WithKeys("right"))
			if cfg.From != "" {
				t.SetValue(cfg.From)
			}
			t.Focus()
		case 1:
			t.Placeholder = "To"
			t.KeyMap.AcceptSuggestion = key.NewBinding(key.WithKeys("right"))
			if cfg.To != "" {
				t.SetValue(cfg.To)
			}
		case 2:
			t.CharLimit = 10
			t.Width = t.CharLimit
			t.KeyMap.AcceptSuggestion = key.NewBinding(key.WithKeys("right"))
			if cfg.Date != "" {
				t.SetValue(cfg.Date)
			} else {
				t.SetValue(now.Format("02.01.2006"))
			}
		case 3:
			t.CharLimit = 5
			t.Width = t.CharLimit
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

// Init implements tea.Model.
func (m appModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, checkVersionCmd(m.currentVersion))
}

func checkVersionCmd(current string) tea.Cmd {
	return func() tea.Msg {
		newer, _ := util.NewerVersion(current)
		return versionCheckMsg{newerVersion: newer}
	}
}

// userError formats an error for display in the TUI.
func userError(err error) string {
	if errors.Is(err, errNoConnections) ||
		errors.Is(err, errMissingDeparture) ||
		errors.Is(err, errMissingArrival) ||
		errors.Is(err, errConnectionMalformed) {
		return capitalise(err.Error())
	}
	return fmt.Sprintf("Something went wrong: %v", err)
}

// capitalise uppercases the first rune.
func capitalise(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] -= 'a' - 'A'
	}
	return string(r)
}
