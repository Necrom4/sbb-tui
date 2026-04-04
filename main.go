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
	date := flag.String("date", "", "Pre-fill date [DD.MM.YYYY]")
	time := flag.String("time", "", "Pre-fill time [HH:MM]")
	arrival := flag.Bool("arrival", false, "Set date/time as arrival instead of departure time")
	flag.Bool("nerdfont", true, "Use Nerd Font icons")
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

	flag.CommandLine.SortFlags = false
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load config: %v\n", err)
	}

	if *showVersion {
		fmt.Printf("sbb-tui %s\n", version)
		os.Exit(0)
	}

	if *date != "" {
		if err := validateDate(*date); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	if *time != "" {
		if err := validateTime(*time); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	// CLI flags override config file values.
	cfg.From = *from
	cfg.To = *to
	cfg.Date = *date
	cfg.Time = *time
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

func validateDate(s string) error {
	if len(s) != 10 || s[2] != '.' || s[5] != '.' {
		return fmt.Errorf("invalid date format %q, expected DD.MM.YYYY", s)
	}
	for _, i := range []int{0, 1, 3, 4, 6, 7, 8, 9} {
		if s[i] < '0' || s[i] > '9' {
			return fmt.Errorf("invalid date format %q, expected DD.MM.YYYY", s)
		}
	}
	day := int(s[0]-'0')*10 + int(s[1]-'0')
	month := int(s[3]-'0')*10 + int(s[4]-'0')
	if day < 1 || day > 31 || month < 1 || month > 12 {
		return fmt.Errorf("invalid date %q: day must be 01-31, month must be 01-12", s)
	}
	return nil
}

func validateTime(s string) error {
	if len(s) != 5 || s[2] != ':' {
		return fmt.Errorf("invalid time format %q, expected HH:MM", s)
	}
	for _, i := range []int{0, 1, 3, 4} {
		if s[i] < '0' || s[i] > '9' {
			return fmt.Errorf("invalid time format %q, expected HH:MM", s)
		}
	}
	hour := int(s[0]-'0')*10 + int(s[1]-'0')
	minute := int(s[3]-'0')*10 + int(s[4]-'0')
	if hour > 23 || minute > 59 {
		return fmt.Errorf("invalid time %q: hours must be 00-23, minutes must be 00-59", s)
	}
	return nil
}
