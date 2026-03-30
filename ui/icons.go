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
	towards   string
	filledDot string
	hollowDot string
	horizLine string
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
		platform: "Pl.",
		stop:     "Stop",
		towards:  "→",

		filledDot: "●",
		hollowDot: "○",
		horizLine: "─",
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
		icons.arrival = "⤙"
		icons.departure = "⤚"
		icons.search = "⌕"
		icons.swap = "↔"
		icons.vehicle = "◇"
		icons.walk = "walk:"
		icons.prompt = "⏵ "
	} else {
		icons.arrival = "󰗔"
		icons.departure = ""
		icons.search = ""
		icons.swap = ""
		icons.vehicle = ""
		icons.walk = ""
		icons.prompt = " "
	}

	return icons
}

func (ic iconSet) platformLabel(platform string) string {
	if len(platform) > 0 && platform[0] >= 'A' && platform[0] <= 'Z' {
		return ic.stop
	}
	return ic.platform
}
