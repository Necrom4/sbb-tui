package ui

import (
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/text/unicode/norm"

	"github.com/necrom4/sbb-tui/api"
	"github.com/necrom4/sbb-tui/model"
)

// Update implements tea.Model.
func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if active.kind == kindButton {
				return m, tea.Quit
			}

		case "enter":
			if err := m.validateInputs(); err != nil {
				m.errorMsg = err
				m.connections = nil
				m.searched = false
				m.resultIndex = 0
				return m, nil
			}
			m.loading = true
			m.connections = nil
			m.errorMsg = nil
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
				if err := m.validateInputs(); err != nil {
					m.errorMsg = err
					m.connections = nil
					m.searched = false
					m.resultIndex = 0
					return m, nil
				}
				m.loading = true
				m.connections = nil
				m.errorMsg = nil
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
				if item.kind == kindInput {
					if item.index == m.headerOrder[m.tabIndex].index {
						cmds = append(cmds, m.inputs[item.index].Focus())
					} else {
						m.inputs[item.index].Blur()
					}
				}
			}
			return m, tea.Batch(cmds...)

		// Disable autocompletion if cursor is not at the end of the stirng.
		case "right":
			active := m.headerOrder[m.tabIndex]

			if active.kind == kindInput {
				input := m.inputs[active.index]

				if input.Position() < len(input.Value()) {
					original := input.KeyMap.AcceptSuggestion
					input.KeyMap.AcceptSuggestion = key.NewBinding() // empty binding

					var cmd tea.Cmd
					m.inputs[active.index], cmd = input.Update(msg)
					m.inputs[active.index].KeyMap.AcceptSuggestion = original

					return m, cmd
				}
			}

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

	case suggestionsMsg:
		if msg.err == nil {
			userInput := m.inputs[msg.inputIndex].Value()
			m.inputs[msg.inputIndex].SetSuggestions(adaptSuggestions(userInput, msg.names))
		}
		return m, nil

	case dataMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = fmt.Errorf("failed to fetch connections: %w", msg.err)
			return m, nil
		}
		m.connections = msg.connections
		m.resultIndex = 0
		m.detailScrollY = 0
		if len(m.connections) == 0 {
			m.errorMsg = errNoConnections
		}
		return m, nil

	case versionCheckMsg:
		m.newerVersion = msg.newerVersion
		return m, nil
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *appModel) updateInputs(msg tea.Msg) tea.Cmd {
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

func (m appModel) validateInputs() error {
	if m.inputs[0].Value() == "" {
		return errMissingDeparture
	}
	if m.inputs[1].Value() == "" {
		return errMissingArrival
	}
	return nil
}

func fetchSuggestionsCmd(inputIndex int, query string) tea.Cmd {
	return func() tea.Msg {
		names, err := api.FetchLocations(query)
		return suggestionsMsg{inputIndex: inputIndex, names: names, err: err}
	}
}

func completeDate(partial string) string {
	now := time.Now().In(model.SwissLocation)
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

// toAPIDate converts Swiss date format (DD.MM.YYYY) to API format (YYYY-MM-DD).
func toAPIDate(swiss string) string {
	parts := strings.Split(swiss, ".")
	if len(parts) != 3 {
		return swiss
	}
	return parts[2] + "-" + parts[1] + "-" + parts[0]
}

func (m appModel) searchCmd() tea.Cmd {
	return func() tea.Msg {
		res, err := api.FetchConnections(
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			toAPIDate(completeDate(m.inputs[2].Value())),
			completeTime(m.inputs[3].Value()),
			m.isArrivalTime,
			m.maxVisibleConnections(),
		)
		return dataMsg{connections: res, err: err}
	}
}

// adaptSuggestions grafts the user's input onto the front of each suggestion
// so that it satisfies the textinput widget's HasPrefix matching.
// (e.g. "zur" (input) + "Zürich HB" (suggestion) = "zurich HB")
func adaptSuggestions(userInput string, suggestions []string) []string {
	if userInput == "" {
		return suggestions
	}
	lower := strings.ToLower(userInput)
	out := make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		idx := prefixMatchLen(strings.ToLower(s), lower)
		if idx > 0 {
			out = append(out, userInput+s[idx:])
		}
	}
	return out
}

// prefixMatchLen returns the offset into `suggestion`
// that has been fuzzy matched against `input`.
func prefixMatchLen(suggestion, input string) int {
	si, ii := 0, 0
	for si < len(suggestion) && ii < len(input) {
		sr, sw := utf8.DecodeRuneInString(suggestion[si:])
		ir, iw := utf8.DecodeRuneInString(input[ii:])

		if sr == ir {
			si += sw
			ii += iw
			continue
		}

		if !unicode.IsLetter(sr) && !unicode.IsDigit(sr) {
			si += sw
			continue
		}

		if foldRune(sr) == foldRune(ir) {
			si += sw
			ii += iw
			continue
		}

		// No match
		return 0
	}

	if ii < len(input) {
		return 0 // didn't consume the entire user input
	}
	return si
}

// foldRune strips common diacritics by decomposing to NFD and
// returning only the first (base) rune.
// (e.g. ü → u + combining diaeresis)
func foldRune(r rune) rune {
	decomposed := norm.NFD.String(string(r))
	base, _ := utf8.DecodeRuneInString(decomposed)
	return base
}
