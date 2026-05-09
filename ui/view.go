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

	latestReleaseURL = "https://github.com/Necrom4/sbb-tui/releases/latest"
)

// View implements tea.Model.
func (m appModel) View() string {
	if m.width < minTermWidth || m.height < minTermHeight {
		msg := fmt.Sprintf("Terminal too small (%dx%d)\nMinimum size: %dx%d", m.width, m.height, minTermWidth, minTermHeight)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			m.styles.warningBold.Render(msg))
	}

	header := m.renderHeader()
	results := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.styles.text.
			Height(m.resultsHeight()).
			Render(m.renderResults()),
		m.styles.text.
			Height(m.resultsHeight()).
			Render(m.renderDetailedResult()),
	)

	footer := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		results,
		footer,
	)
}

// ----- Layout calculations -----

// contentWidth returns the usable horizontal width of the TUI.
func (m appModel) contentWidth() int {
	return max(m.width, 0)
}

// resultsHeight returns the vertical space available for the results pane.
func (m appModel) resultsHeight() int {
	return max(m.height-headerHeight-helpBarHeight, 0)
}

// maxVisibleConnections is the count of result rows that fit in the results pane.
func (m appModel) maxVisibleConnections() int {
	return max(m.resultsHeight()/simpleConnHeight, 1)
}

// resultBoxWidth returns the width of one result column (simple list or detail).
func (m appModel) resultBoxWidth() int {
	return max((m.width-simpleConnMargin)/2, resultMargin+stopsLineMinWidth+stopsLineFixedWidth)
}

// headerFixedWidth returns the total horizontal space taken by the
// header items, treating the From/To inputs as their per-item overhead
// only so the leftover width can be split between them.
func (m appModel) headerFixedWidth() int {
	width := 0
	for i, item := range m.headerOrder {
		if item.id == "from" || item.id == "to" {
			width += borderSize + 2 + lipgloss.Width(m.inputs[item.index].Prompt)
			continue
		}
		width += lipgloss.Width(m.renderHeaderItem(i))
	}
	return width
}

// ----- Header rendering -----

// renderHeader joins every header item horizontally.
func (m appModel) renderHeader() string {
	var headerItems []string
	for i := range m.headerOrder {
		headerItems = append(headerItems, m.renderHeaderItem(i))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, headerItems...)
}

// ----- Footer rendering -----

// renderHelpBar returns the bottom-bar key hints.
func (m appModel) renderHelpBar() string {
	bindings := []struct{ key, desc string }{
		{m.icons.keyTab, "navigate"},
		{m.icons.keyEnter, "search"},
		{m.icons.keySpace, "toggle"},
		{m.icons.keyUpDw, "results"},
		{m.icons.keyUPDW, "scroll"},
		{m.icons.keyRight, "complete"},
		{m.icons.keyEsc, "quit"},
	}

	parts := make([]string, len(bindings))
	for i, b := range bindings {
		parts[i] = m.styles.helpKey.Render(b.key) + " " + m.styles.helpDesc.Render(b.desc)
	}

	return " " + strings.Join(parts, "   ")
}

// renderVersionBadge returns the bottom-right "SBB-TUI vX.Y.Z" badge,
// shrinking or hiding itself when availableWidth is too narrow.
func (m appModel) renderVersionBadge(availableWidth int) string {
	const (
		appName = "SBB-TUI"
		minGap  = 2
	)

	if availableWidth <= minGap {
		return ""
	}

	if m.newerVersion != "" {
		full := fmt.Sprintf(
			"%s %s %s%s%s",
			m.styles.text.Render(appName),
			m.styles.textMuted.Render(m.currentVersion),
			m.styles.warning.Render("(latest: "),
			m.styles.warning.Render(renderLink(m.newerVersion, latestReleaseURL)),
			m.styles.warning.Render(")"),
		)
		if lipgloss.Width(full)+minGap <= availableWidth {
			return full
		}
	}

	short := fmt.Sprintf(
		"%s %s",
		m.styles.text.Render(appName),
		m.styles.textMuted.Render(m.currentVersion),
	)

	if lipgloss.Width(short)+minGap <= availableWidth {
		return short
	}

	return ""
}

// renderFooter combines the help bar and the version badge with stretching whitespace.
func (m appModel) renderFooter() string {
	helpBar := m.renderHelpBar()
	versionBadge := m.renderVersionBadge(m.width - lipgloss.Width(helpBar))

	if versionBadge == "" {
		return helpBar
	}

	gap := m.width - lipgloss.Width(helpBar) - lipgloss.Width(versionBadge)
	return helpBar + strings.Repeat(" ", gap) + versionBadge
}

// renderHeaderItem renders a single header item (input or button) styled
// per its focus state.
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
			// Clip text so the rendered suggestion never overflows the input box.
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

// ----- Results layout -----

// renderResults dispatches to the loading/error/start/list views.
func (m appModel) renderResults() string {
	if m.loading {
		return m.renderLoading()
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

// onStartScreen reports whether the start screen (logo + tagline) is currently shown.
func (m appModel) onStartScreen() bool {
	return m.errorMsg == nil && !m.loading && len(m.connections) == 0 && !m.searched
}

// renderStartScreen renders the centered logo, tagline and optional update notice.
func (m appModel) renderStartScreen() string {
	logo := sbbLogo
	if m.nerdFont {
		logo = sbbLogoNerdFont
	}
	logo = strings.TrimRight(logo, "\n")

	coloredLogo := m.renderLogo(logo)
	text := m.renderStartTagline("Enter stations above to see timetables")

	block := lipgloss.JoinVertical(lipgloss.Center, text, "", coloredLogo)

	if m.newerVersion != "" {
		latestVersion := renderLink(m.newerVersion, latestReleaseURL)
		label := fmt.Sprintf("Update available: %s", latestVersion)
		block = lipgloss.JoinVertical(lipgloss.Center, block, "", m.styles.active.Render(label))
	}

	width := max(m.contentWidth(), 0)
	height := m.resultsHeight()

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, block)
}

// renderDetailedResult renders the right-hand detail box for the currently selected connection.
func (m appModel) renderDetailedResult() string {
	if len(m.connections) == 0 {
		return ""
	}

	boxWidth := max(m.width-borderSize*2-m.resultBoxWidth(), 0)
	return m.renderFullConnection(m.connections[m.resultIndex], boxWidth)
}
