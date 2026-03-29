# SBB-TUI

TUI client for Switzerland's public transports timetables, inspidered by the SBB/CFF/FFS [app](https://www.sbb.ch/).

<img width="1382" height="1054" alt="Bildschirmfoto 2026-03-01 um 11 43 00" src="https://github.com/user-attachments/assets/f3600847-50ce-418b-b682-5249ee00ab6f" />

## 🚀 Features

- 📍 **Any station**: search connections between any Swiss (and neighbouring) public transport stations
- 🚆 **Any transport**: train, tram, bus, boat, cable cars (and even funiculars!)
- 🧾 **Detailed journey view** with: transfers, platforms, delays, vehicle information and walk sections
- ➡️ **Autocompletion** of station names
- 🚩 **CLI flags** to pre-fill fields for quick lookups
- 🧭 Clickable **Google Maps links** for walking sections
- **<+>** SBB brand styling with Nerd Font **icons** ([Nerd Font](https://www.nerdfonts.com/) recommended, Unicode fallback via `--no-nerdfont`)

## 📦 Install

```sh
# homebrew
brew install necrom4/tap/sbb-tui
# go
go install github.com/necrom4/sbb-tui
# aur
yay -S sbb-tui
```

> or download your preferred executable from [Releases](https://github.com/Necrom4/sbb-tui/releases)

## Build from source

```sh
git clone https://github.com/necrom4/sbb-tui.git
cd sbb-tui
go build
```

## 🚀 Usage

1. Run `sbb-tui`
2. Input **departure** and **arrival** locations (navigate with `tab`).
3. Add optional information such as **date**, **time**, and **whether** those are for departure or arrival.
4. Press `Enter` to view the results (navigate with arrows).

## ⚙️ Configuration

Add your optional config at `$HOME/.config/sbb-tui/config.yaml`

> OS default config paths are also supported (such as `~/Library/Application Support/sbb-tui/config.yaml` in macOS)

```yaml
# default configuration
theme:
  text:           "#FFFFFF"
  errorText:      "#D82E20"
  ghostText:      "#888888"
  activeBorder:   "#D82E20"
  inactiveBorder: "#484848"
  dimmedBorder:   "#862010"
  warningFlag:    "#D82E20"
  keysFg:         "#FFFFFF"
  keysBg:         "#484848"
  vehicleFg:      "#FFFFFF"
  vehicleBg:      "#2E3279"
  modelFg:        "#FFFFFF"
  modelBg:        "#D82E20"
  companyFg:      "#484848"
  companyBg:      "#FFFFFF"
  logo:           "#FFFFFF"
```

More **themes** can be found at [docs/themes.md](https://github.com/Necrom4/sbb-tui/blob/master/docs/themes.md)

## 🚩 Options

```sh
# sbb-tui --help
sbb-tui - Swiss SBB/CFF/FFS timetable app for the terminal

Usage:
  sbb-tui [flags]

Flags:
      --arrival       Use arrival time instead of departure time
      --date string   Pre-fill date (DD.MM.YYYY)
      --from string   Pre-fill departure station
      --no-nerdfont   Use Unicode fallback icons instead of Nerd Font icons
      --time string   Pre-fill time (HH:MM)
      --to string     Pre-fill arrival station
  -v, --version       Print version and exit
```

## ❓ Why

> I often work in the train, passing through remote regions of Switzerland where I'll have to wait up to an entire minute to finally be able to load the SBB website/app and get the much needed information about my next connection (I have a cheap cellular data subscription). Someday I fell onto the incredible [Swiss public transport API](https://transport.opendata.ch/docs.html) and decided it was the perfect occasion to learn how to create TUIs.

## 📝 TODO

- [ ] **Stationboard** mode, returns a list of the next departures at a specific station.
- [ ] **Via** input field, allows adding an intermediary station.
- [ ] Connection warnings
- [ ] Better keymaps/navigation logic
- [ ] Better keymap help
- [ ] Suggestions when writing strings without accent (writing "zurich", "Zürich" isn't suggested)
- [ ] Revise UI for not-so-wide terminals
- [ ] Scroll icons as hint in border of scrollable detailedRender window
- [ ] Only autocomplete with cursor at last character, otherwise move cursor right
- [ ] Shorten date/time fields by one character length by either extending cursor placement to character before right border, or by removing cursor when finished at end of input CharLimit
- [ ] Better error messages when date/time flags are given non-conforming strings
- [ ] Protect date/time inputs from delete->reinserts in the middle of the string, breaks

## Star History

<a href="https://www.star-history.com/?repos=necrom4%2Fsbb-tui&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/image?repos=necrom4/sbb-tui&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/image?repos=necrom4/sbb-tui&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/image?repos=necrom4/sbb-tui&type=date&legend=top-left" />
 </picture>
</a>
