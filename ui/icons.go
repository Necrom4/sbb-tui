package ui

type iconSet struct {
	// Mode-dependent (Nerd Font vs Unicode fallback)
	arrival   string
	departure string
	platform  string
	stop      string
	search    string
	swap      string
	vehicle   string
	walk      string
	prompt    string

	// Mode-invariant
	towards    string
	filledDot  string
	hollowDot  string
	horizLine  string
	vertLine   string
	scrollUp   string
	scrollDown string
	keyTab     string
	keyEnter   string
	keySpace   string
	keyUpDw    string
	keyUPDW    string
	keyRight   string
	keyEsc     string
}

func newIconSet(nerdFont bool) iconSet {
	icons := iconSet{
		// Shared symbols
		platform: "Pl.",
		stop:     "Stop",
		towards:  "→",

		filledDot:  "●",
		hollowDot:  "○",
		horizLine:  "─",
		vertLine:   "│",
		scrollUp:   "▲",
		scrollDown: "▼",

		keyTab:   "⇥",
		keyEnter: "↵",
		keySpace: "␣",
		keyUpDw:  "↕",
		keyUPDW:  "⇧↕",
		keyRight: "→",
		keyEsc:   "⎋",
	}

	if nerdFont {
		icons.arrival = "󰗔"
		icons.departure = ""
		icons.search = ""
		icons.swap = ""
		icons.vehicle = ""
		icons.walk = ""
		icons.prompt = " "
		icons.scrollUp = "▴"
		icons.scrollDown = "▾"
	} else {
		icons.arrival = "⤙"
		icons.departure = "⤚"
		icons.search = "⌕"
		icons.swap = "↔"
		icons.vehicle = "◇"
		icons.walk = "walk:"
		icons.prompt = "⏵ "
		icons.scrollUp = "▲"
		icons.scrollDown = "▼"
	}

	return icons
}

func (ic iconSet) platformLabel(platform string) string {
	if len(platform) > 0 && platform[0] >= 'A' && platform[0] <= 'Z' {
		return ic.stop
	}
	return ic.platform
}
