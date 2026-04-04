package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/necrom4/sbb-tui/model"
)

// buildDetailLines builds the raw content lines for a connection's detail view.
func (m appModel) buildDetailLines(c model.Connection, innerWidth int) []string {
	var lines []string

	// Pre-compute widest label and value widths so platform values align.
	labelCol := 0
	valueCol := 0
	for _, section := range c.Sections {
		if section.Journey == nil {
			continue
		}
		for _, p := range []string{section.Departure.Platform, section.Arrival.Platform} {
			if p != "" {
				label := m.icons.platformLabel(p)
				if lw := len([]rune(label)); lw > labelCol {
					labelCol = lw
				}
				if vw := len([]rune(p)); vw > valueCol {
					valueCol = vw
				}
			}
		}
	}
	// platformCol is the total visible width: padded label + space + widest value
	platformCol := 0
	if labelCol > 0 {
		platformCol = labelCol + 1 + valueCol
	}

	for i, section := range c.Sections {
		isFirst := i == 0
		isLast := i == len(c.Sections)-1

		if section.Walk != nil {
			lines = append(lines, m.renderWalkSection(section)...)
		} else if section.Journey != nil {
			lines = append(lines, m.renderJourneySection(section, innerWidth, labelCol, platformCol, isFirst, isLast)...)
		}

		if !isLast {
			nextIsWalk := c.Sections[i+1].Walk != nil
			currentIsWalk := section.Walk != nil
			hasArrDelay := section.Journey != nil && section.Arrival.Delay > 0
			if currentIsWalk {
				lines = append(lines, "")
			} else if hasArrDelay {
				if nextIsWalk {
				} else {
					lines = append(lines, "")
				}
			} else {
				if nextIsWalk {
					lines = append(lines, "")
				} else {
					lines = append(lines, "", "")
				}
			}
		}
	}

	return lines
}

// maxDetailScroll returns the maximum useful scroll offset for the current
// connection's detail view. Returns 0 when there's nothing to scroll.
func (m appModel) maxDetailScroll() int {
	if len(m.connections) == 0 || m.resultIndex >= len(m.connections) {
		return 0
	}
	c := m.connections[m.resultIndex]
	boxWidth := max(m.width-borderSize*2-m.resultBoxWidth(), 0)
	innerWidth := max(boxWidth-borderSize-(detailPaddingH*2), 0)

	lines := m.buildDetailLines(c, innerWidth)
	content := strings.Join(lines, "\n")
	wrapped := m.styles.text.Width(innerWidth).Render(content)
	visLines := strings.Split(wrapped, "\n")

	detailFrame := m.styles.detailedResult.GetVerticalFrameSize()
	boxHeight := max(m.resultsHeight()-detailFrame, 0)

	if len(visLines) <= boxHeight {
		return 0
	}
	return len(visLines) - boxHeight
}

func (m appModel) renderFullConnection(c model.Connection, width int) string {
	innerWidth := max(width-borderSize-(detailPaddingH*2), 0)
	lines := m.buildDetailLines(c, innerWidth)

	detailFrame := m.styles.detailedResult.GetVerticalFrameSize()
	boxHeight := max(m.resultsHeight()-detailFrame, 0)

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

func (m appModel) renderJourneySection(section model.Section, width, labelCol, platformCol int, isFirst, isLast bool) []string {
	var lines []string

	const timeCol = 5
	const symbolCol = 5

	depTime := section.Departure.Scheduled.Local().Format("15:04")
	depDelay := section.Departure.Delay
	depStation := section.Departure.Station.Name
	depPlatform := section.Departure.Platform

	depDot := m.icons.hollowDot
	if isFirst {
		depDot = m.icons.filledDot
	}

	depLine := m.formatStationLine(depTime, depDot, depStation, depPlatform, width, timeCol, symbolCol, labelCol, platformCol, true)
	lines = append(lines, depLine)

	indent := strings.Repeat(" ", timeCol)
	spacingLine := fmt.Sprintf("%s  %s", indent, m.icons.vertLine)

	if depDelay > 0 {
		delayStr := m.styles.warningBold.Render(fmt.Sprintf("%*s'", timeCol, fmt.Sprintf("+%d", depDelay)))
		lines = append(lines, fmt.Sprintf("%s %s", delayStr, m.styles.bold.Render(m.icons.vertLine)))
	} else {
		lines = append(lines, spacingLine)
	}

	vehicleIcon := m.styles.vehicleIcon.Render(" " + m.icons.vehicle + " ")
	vehicleModel := m.styles.vehicleModel.Render(section.Journey.Category + " " + section.Journey.Number)
	company := m.styles.company.Render(section.Journey.Operator)
	vehicleLine := fmt.Sprintf("%s  %s  %s %s %s", indent, m.icons.vertLine, vehicleIcon, vehicleModel, company)
	lines = append(lines, vehicleLine)

	destLine := fmt.Sprintf("%s  %s   %s", indent, m.icons.vertLine, m.styles.textMuted.Render(m.icons.towards+" "+section.Journey.To))
	lines = append(lines, destLine)

	lines = append(lines, spacingLine)

	arrTime := section.Arrival.Scheduled.Local().Format("15:04")
	arrDelay := section.Arrival.Delay
	arrStation := section.Arrival.Station.Name
	arrPlatform := section.Arrival.Platform

	arrSymbol := m.icons.vertLine
	if isLast {
		arrSymbol = m.icons.filledDot
	}

	arrLine := m.formatStationLine(arrTime, arrSymbol, arrStation, arrPlatform, width, timeCol, symbolCol, labelCol, platformCol, false)
	lines = append(lines, arrLine)

	if arrDelay > 0 {
		delayStr := m.styles.warningBold.Render(fmt.Sprintf("%*s'", timeCol, fmt.Sprintf("+%d", arrDelay)))
		lines = append(lines, delayStr)
	}

	return lines
}

func googleMapsURL(s model.Section) string {
	dep := s.Departure.Station.Coordinate
	arr := s.Arrival.Station.Coordinate
	return fmt.Sprintf("https://www.google.com/maps/dir/?api=1&origin=%f,%f&destination=%f,%f&travelmode=walking",
		dep.X, dep.Y, arr.X, arr.Y)
}

// renderLink generates an OSC 8 terminal hyperlink.
func renderLink(text, url string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

func (m appModel) renderWalkSection(section model.Section) []string {
	var lines []string

	walkDuration := ""
	if section.Walk != nil {
		dur := section.Walk.Duration
		if dur > 0 {
			walkDuration = fmt.Sprintf("%d", dur/60)
		} else {
			depTime := section.Departure.Scheduled.Time
			arrTime := section.Arrival.Scheduled.Time
			if !depTime.IsZero() && !arrTime.IsZero() {
				walkDuration = fmt.Sprintf("%d", int(arrTime.Sub(depTime).Minutes()))
			}
		}
		url := googleMapsURL(section)

		// TODO: add `` icon and set that as clickable url link instead of the time
		walkDuration = renderLink(walkDuration, url)
	}

	walkLine := fmt.Sprintf("%s  %s %s'", strings.Repeat(" ", 5), m.icons.walk, walkDuration)
	lines = append(lines, walkLine)

	return lines
}

func (m appModel) formatStationLine(timeStr, symbol, station, platform string, width, timeCol, symbolCol, labelCol, platformCol int, bold bool) string {
	textStyle := m.styles.text
	if bold {
		textStyle = m.styles.bold
	}

	timePart := textStyle.Render(timeStr)

	symbolPart := fmt.Sprintf("  %s  ", textStyle.Render(symbol))

	platformPart := ""
	if platform != "" {
		label := m.icons.platformLabel(platform)
		leadingPad := strings.Repeat(" ", max(labelCol-len([]rune(label)), 0))
		labelPart := leadingPad + m.styles.textMuted.Render(label)
		valuePart := textStyle.Render(platform)
		platformPart = labelPart + " " + valuePart
	}

	fixedWidth := timeCol + symbolCol
	if platformCol > 0 {
		fixedWidth += platformCol
	}
	availableForStation := max(width-fixedWidth-1, 5)

	truncatedStation := truncateString(station, availableForStation)
	stationPart := textStyle.Render(truncatedStation)

	stationLen := len([]rune(truncatedStation))
	padding := max(availableForStation-stationLen, 1)

	if platformPart != "" {
		return fmt.Sprintf("%s%s%s%s%s",
			timePart, symbolPart, stationPart, strings.Repeat(" ", padding), platformPart)
	}
	return fmt.Sprintf("%s%s%s", timePart, symbolPart, stationPart)
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

func (m appModel) renderSimpleConnection(c model.Connection, index int, width int) string {
	firstVehicle := -1
	lastVehicle := -1
	for i := range c.Sections {
		if c.Sections[i].Journey != nil {
			if firstVehicle == -1 {
				firstVehicle = i
			}
			lastVehicle = i
		}
	}

	style := m.styles.inactive.Width(width)
	if index == m.resultIndex {
		style = m.styles.active.Width(width)
	}

	if firstVehicle == -1 {
		return m.styles.error.Width(width).Padding(1, 2).Render(userError(errConnectionMalformed))
	}

	lineContentWidth := max(width-style.GetHorizontalFrameSize()-2, 0)

	vehicleIcon := m.styles.vehicleIcon.Render(" " + m.icons.vehicle + " ")
	vehicleModel := m.styles.vehicleModel.Render(c.Sections[firstVehicle].Journey.Category + " " + c.Sections[firstVehicle].Journey.Number)
	company := m.styles.company.Render(c.Sections[firstVehicle].Journey.Operator)
	endStop := m.styles.text.Render(c.Sections[firstVehicle].Journey.To)

	dep := c.Sections[firstVehicle].Departure.Scheduled.Local().Format("15:04")
	arr := c.To.Arrival.Local().Format("15:04")
	departure := m.styles.bold.Render(dep)
	arrival := m.styles.bold.Render(arr)

	departureDelay := m.formatDelay(c.Sections[firstVehicle].Departure.Delay)
	arrivalDelay := m.formatDelay(c.Sections[lastVehicle].Arrival.Delay)

	timelinePrefix := ""
	if c.Sections[0].Walk != nil {
		walkMinutes := int(c.Sections[0].Arrival.Scheduled.Sub(c.Sections[0].Departure.Scheduled).Minutes())
		if walkMinutes > 0 {
			timelinePrefix = m.icons.walk + " " + m.styles.text.Render(fmt.Sprintf("%d'", walkMinutes)) + "  "
		}
	}

	depGap := "  "
	if departureDelay != "" {
		depGap = " "
	}
	arrGap := "  "
	if arrivalDelay != "" {
		arrGap = " "
	}

	timelineFixedWidth := lipgloss.Width(timelinePrefix) +
		lipgloss.Width(departure) +
		lipgloss.Width(departureDelay) + len(depGap) +
		len(arrGap) +
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
		platform = c.From.Platform
	}
	if platform != "" {
		label := m.icons.platformLabel(platform)
		platformInfo = m.styles.textMuted.Render(label) + " " + m.styles.text.Render(platform)
	}

	duration := m.styles.text.Render(formatDuration(c.Duration))

	bottomLinePadding := max(lineContentWidth-lipgloss.Width(platformInfo)-lipgloss.Width(duration), 1)

	content := fmt.Sprintf("\n  %s %s %s  %s\n\n  %s%s%s%s%s  %s%s\n\n  %s%s%v\n",
		vehicleIcon,
		vehicleModel,
		company,
		endStop,
		timelinePrefix,
		departure,
		departureDelay,
		depGap,
		stopsLine,
		arrival,
		arrivalDelay,
		platformInfo,
		strings.Repeat(" ", bottomLinePadding),
		duration,
	)

	return style.Render(content)
}

// formatDuration converts the API duration format (e.g. "00d01:15:00") to a
// human-readable string like "1 h 15 min" or "15 min".
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

func (m appModel) formatDelay(delay int) string {
	if delay > 0 {
		return m.styles.warningBold.Render(fmt.Sprintf(" +%d'", delay))
	}
	return ""
}

func (m appModel) renderStopsLine(c model.Connection, totalWidth int) string {
	if len(c.Sections) == 0 {
		return m.icons.filledDot + m.icons.horizLine + m.icons.horizLine + m.icons.filledDot
	}

	var sectionDurations []time.Duration
	var totalSectionDuration time.Duration
	for _, s := range c.Sections {
		// Skip walking sections
		if s.Journey == nil {
			continue
		}
		dep := s.Departure.Scheduled.Time
		arr := s.Arrival.Scheduled.Time
		if !dep.IsZero() && !arr.IsZero() {
			dur := arr.Sub(dep)
			sectionDurations = append(sectionDurations, dur)
			totalSectionDuration += dur
		}
	}

	if totalSectionDuration == 0 || len(sectionDurations) == 0 {
		// Fallback to equal distribution
		return m.icons.filledDot + strings.Repeat(m.icons.horizLine+m.icons.horizLine+m.icons.hollowDot, c.Transfers) + m.icons.horizLine + m.icons.horizLine + m.icons.filledDot
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

		sb.WriteString(strings.Repeat(m.icons.horizLine, lineChars))
		if i < len(sectionDurations)-1 {
			sb.WriteString(m.icons.hollowDot)
		} else {
			sb.WriteString(m.icons.filledDot)
		}
	}

	return sb.String()
}
