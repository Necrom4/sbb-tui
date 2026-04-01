package main

import (
	"fmt"
	"os"
	_ "time/tzdata" // embed timezone database so Europe/Zurich always resolves

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
	flag.Bool("nerdfont", true, "Use Nerd Font icons (disable with --nerdfont=false)")
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

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load config: %v\n", err)
	}

	if *showVersion {
		fmt.Printf("sbb-tui %s\n", version)
		os.Exit(0)
	}

	// CLI flags override config file values.
	cfg.From = *from
	cfg.To = *to
	cfg.Date = *date
	cfg.Time = *timeStr
	cfg.IsArrivalTime = *arrival
	cfg.CurrentVersion = version

	if flag.CommandLine.Changed("nerdfont") {
		nf, _ := flag.CommandLine.GetBool("nerdfont")
		cfg.NerdFont = nf
	}

	m := ui.NewModel(cfg)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
}
