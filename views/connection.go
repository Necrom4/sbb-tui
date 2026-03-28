package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/necrom4/sbb-tui/models"
	"github.com/necrom4/sbb-tui/utils"

	"github.com/charmbracelet/lipgloss"
)

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

		// TODO: add `` icon and set that as clickable url link instead of the time
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
