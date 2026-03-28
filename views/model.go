package views

import (
	"time"

	"github.com/necrom4/sbb-tui/config"
	"github.com/necrom4/sbb-tui/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	// Focusable item kinds
	KindInput int = iota
	KindButton
)

type focusable struct {
	kind  int
	id    string
	index int
}

// DataMsg is sent when the API returns connection results.
type DataMsg struct {
	connections []models.Connection
	err         error
}

// SuggestionsMsg is sent when station name suggestions are fetched.
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
