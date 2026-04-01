// Package ui implements the Bubbletea TUI for SBB timetable queries.
package ui

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	//go:embed sbb-logo.txt
	sbbLogo string

	//go:embed sbb-logo-nerdfont.txt
	sbbLogoNerdFont string
)

// View implements tea.Model.
func (m appModel) View() string {
	if m.width < minTermWidth || m.height < minTermHeight {
		msg := fmt.Sprintf("Terminal too small (%dx%d)\nMinimum size: %dx%d", m.width, m.height, minTermWidth, minTermHeight)
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
		results,
		helpBar,
	)
}

// Layout calculations

func (m appModel) contentWidth() int {
	return max(m.width, 0)
}

func (m appModel) resultsHeight() int {
	return max(m.height-headerHeight-helpBarHeight, 0)
}

func (m appModel) maxVisibleConnections() int {
	return max(m.resultsHeight()/simpleConnHeight, 1)
}

func (m appModel) resultBoxWidth() int {
	return max((m.width-simpleConnMargin)/2, resultMargin+stopsLineMinWidth+stopsLineFixedWidth)
}

func (m appModel) headerFixedWidth() int {
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

// Header rendering

func (m appModel) renderHeader() string {
	var headerItems []string
	for i := range m.headerOrder {
		headerItems = append(headerItems, m.renderHeaderItem(i))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, headerItems...)
}

func (m appModel) renderHelpBar() string {
	var keys []struct{ key, desc string }

	if m.inputMode {
		keys = []struct{ key, desc string }{
			{m.icons.keyTab, "complete"},
			{m.icons.keyEnter, "search"},
			{m.icons.keyEsc, "back"},
		}
	} else {
		keys = []struct{ key, desc string }{
			{m.icons.keyLeftRight, "navigate"},
			{m.icons.keyEnter, "select"},
			{m.icons.keyTab, "results"},
			{m.icons.keyUpDown, "scroll"},
			{m.icons.keyQuit, "quit"},
		}
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, m.styles.helpKey.Render(k.key)+" "+m.styles.helpDesc.Render(k.desc))
	}

	return " " + strings.Join(parts, "   ")
}

func (m appModel) renderHeaderItem(idx int) string {
	item := m.headerOrder[idx]
	style := m.styles.inactive
	if m.tabIndex == idx {
		style = m.styles.active
	}

	if item.kind == kindInput {
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
		icon = m.icons.swap
	case "isArrivalTime":
		if m.isArrivalTime {
			icon = m.icons.arrival
		} else {
			icon = m.icons.departure
		}
	case "search":
		icon = m.icons.search
	}
	return style.Render(icon)
}

// Results layout

func (m appModel) renderResults() string {
	if m.loading {
		return "\n  Searching connections..."
	}

	if m.errorMsg != nil {
		return "\n  " + m.styles.warning.Render(userError(m.errorMsg))
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

func (m appModel) renderStartScreen() string {
	logo := sbbLogo
	if m.nerdFont {
		logo = sbbLogoNerdFont
	}
	logo = strings.TrimRight(logo, "\n")

	coloredLogo := m.styles.logo.Render(logo)

	text := m.styles.ghostText.Render("Enter stations above to see timetables")

	block := lipgloss.JoinVertical(lipgloss.Center, text, "", coloredLogo)

	if m.newerVersion != "" {
		url := "https://github.com/Necrom4/sbb-tui/releases/latest"
		link := renderLink(m.newerVersion, url)
		label := fmt.Sprintf("Update available: %s", link)
		block = lipgloss.JoinVertical(lipgloss.Center, block, "", m.styles.active.Render(label))
	}

	width := max(m.contentWidth(), 0)
	height := m.resultsHeight()

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, block)
}

func (m appModel) renderDetailedResult() string {
	if len(m.connections) == 0 {
		return ""
	}

	boxWidth := max(m.width-borderSize*2-m.resultBoxWidth(), 0)
	return m.renderFullConnection(m.connections[m.resultIndex], boxWidth)
}
