# SBB-TUI

TUI client for Switzerland's public transports timetables, inspidered by the SBB/CFF/FFS [app](https://www.sbb.ch/).

<img width="1382" height="1054" alt="Bildschirmfoto 2026-03-01 um 11 43 00" src="https://github.com/user-attachments/assets/f3600847-50ce-418b-b682-5249ee00ab6f" />

## ❓Why

> I often work in the train, passing through remote regions of Switzerland where I'll have to wait up to an entire minute to finally be able to load the SBB website/app and get the much needed information about my next connection (I have a cheap cellular data subscription). Someday I fell onto the incredible [Swiss public transport API](https://transport.opendata.ch/docs.html) and decided it was the perfect occasion to learn how to create TUIs.

## 📦 Install

```sh
# homebrew
brew tap necrom4/homebrew-tap && brew install sbb-tui
# or go
go install github.com/necrom4/sbb-tui
```

## 🚀 Usage

1. Run `sbb-tui`
2. Input **departure** and **arrival** locations (navigate with `tab`).
3. Add optional information such as **date**, **time**, and **whether** those are for departure or arrival.
4. Press `Enter` to view the results (navigate with arrows).

## ❓ Options

```sh
# sbb-tui --help
sbb-tui - Swiss SBB/CFF/FFS timetable app for the terminal

Usage:
  sbb-tui [flags]

Flags:
      --arrival       Use arrival time instead of departure time
      --date string   Pre-fill date (YYYY-MM-DD)
      --from string   Pre-fill departure station
      --no-nerdfont   Use Unicode fallback icons instead of Nerd Font icons
      --time string   Pre-fill time (HH:MM)
      --to string     Pre-fill arrival station
  -v, --version       Print version and exit
```

## 📝 TODO

- [ ] **Stationboard** mode, returns a list of the next departures at a specific station.
- [ ] Connection warnings
- [ ] Better keymaps/navigation logic
- [ ] Better keymap help
- [ ] Suggestions when writing strings without accent (writing "zurich", "Zürich" isn't suggested)
- [ ] Revise UI for not-so-wide terminals
- [ ] Scroll icons as hint in border of scrollable detailedRender window
