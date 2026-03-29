package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	flag "github.com/spf13/pflag"

	"github.com/necrom4/sbb-tui/config"
	"github.com/necrom4/sbb-tui/ui"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	from := flag.String("from", "", "Pre-fill departure station")
	to := flag.String("to", "", "Pre-fill arrival station")
	date := flag.String("date", "", "Pre-fill date (DD.MM.YYYY)")
	timeStr := flag.String("time", "", "Pre-fill time (HH:MM)")
	arrival := flag.Bool("arrival", false, "Use arrival time instead of departure time")
	noNerdFont := flag.Bool("no-nerdfont", false, "Use Unicode fallback icons instead of Nerd Font icons")
	showVersion := flag.BoolP("version", "v", false, "Print version and exit")

	// --help
	flag.Usage = func() {
		fmt.Println("sbb-tui - Swiss SBB/CFF/FFS timetable app for the terminal")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  sbb-tui [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	theme, err := config.LoadTheme()
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load config:", err)
	}

	if *showVersion {
		fmt.Printf("sbb-tui %s\n", version)
		os.Exit(0)
	}

	cfg := config.Config{
		From:           *from,
		To:             *to,
		Date:           *date,
		Time:           *timeStr,
		IsArrivalTime:  *arrival,
		NoNerdFont:     *noNerdFont,
		Theme:          theme,
		CurrentVersion: version,
	}

	m := ui.NewModel(cfg)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "could not run program:", err)
		os.Exit(1)
	}
}
