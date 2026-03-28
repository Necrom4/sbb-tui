// Package views implements the Bubbletea TUI for SBB timetable queries.
package views

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

func (m model) View() string {
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
		m.styles.dimmedBorder.
			Width(m.contentWidth()).
			Height(m.resultsHeight()).
			Render(results),
		helpBar,
	)
}

// Layout calculations

func (m model) contentWidth() int {
	return max(m.width-borderSize, 0)
}

func (m model) resultsHeight() int {
	return max(m.height-hdrHeight-borderSize-helpBarHeight, 0)
}

func (m model) maxVisibleConnections() int {
	return max(m.resultsHeight()/smplConnHeight, 1)
}

func (m model) resultBoxWidth() int {
	return max((m.width-smplConnMrgn)/2, rsltMrgn+stopsLineMinWidth+stopsLineFixedWidth)
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

// Header rendering

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

// Results layout

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
