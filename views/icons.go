package views

type iconSet struct {
	// Mode-dependent (Nerd Font vs Unicode fallback)
	arr    string
	dpt    string
	plt    string
	srch   string
	swp    string
	vhc    string
	wlk    string
	prompt string

	// Mode-invariant
	twrds     string
	filledDot string
	hollowDot string
	horzLine  string
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
		plt:   "Pl.",
		twrds: "→",

		filledDot: "●",
		hollowDot: "○",
		horzLine:  "─",
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
		icons.arr = "↘"
		icons.dpt = "↗"
		icons.srch = "⌕"
		icons.swp = "⇋"
		icons.vhc = "×"
		icons.wlk = "Walk:"
		icons.prompt = "> "
	} else {
		icons.arr = "󰗔"
		icons.dpt = ""
		icons.srch = ""
		icons.swp = ""
		icons.vhc = ""
		icons.wlk = ""
		icons.prompt = " "
	}

	return icons
}
