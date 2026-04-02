# SBB-TUI

TUI client for Switzerland's public transports timetables, inspired by the SBB/CFF/FFS [app](https://www.sbb.ch/).

<img width="1382" height="1054" alt="Bildschirmfoto 2026-03-01 um 11 43 00" src="https://github.com/user-attachments/assets/f3600847-50ce-418b-b682-5249ee00ab6f" />

## ✨ Features

- 📍 **Any station**: search connections between any Swiss (and neighbouring) public transport stations
- 🚆 **Any transport**: train, tram, bus, boat, cable cars (and even funiculars!)
- 🧾 **Detailed journey view** with: transfers, platforms, delays, vehicle information and walk sections
- ➡️ **Autocompletion** of station names
- 🚩 **CLI flags** to pre-fill fields for quick lookups
- 🧭 Clickable **Google Maps links** for walking sections
- ⇹ SBB brand **styling** with Nerd Font icons ([Nerd Font](https://www.nerdfonts.com/) recommended, Unicode fallback via `--nerdfont=false`)

## 📚 Documentation

- [Install](#📦-install)
- [Usage](#🚀-usage)
- [Configuration](#⚙️-configuration)
- [Options](#🚩-options)
- [Why?](#❓-why)
- [HELP WANTED!](#🍻-help-wanted)

## 📦 Install

### Homebrew

```sh
brew install necrom4/tap/sbb-tui
```

### Go

```sh
go install github.com/necrom4/sbb-tui
```

### AUR

```sh
yay -S sbb-tui # or equivalent AUR package manager
```

> [!NOTE]
> Or download your preferred executable from [Releases](https://github.com/Necrom4/sbb-tui/releases) directly.

### Build from source

```sh
git clone https://github.com/necrom4/sbb-tui.git
cd sbb-tui
go build
```

> [!WARNING]
> Due to the first release being [v1.0.0](https://github.com/Necrom4/sbb-tui/releases/tag/v1.0.0), `v1.X.X` is considered experimental/unstable. What that mostly means is that breaking changes won't trigger a MAJOR version bump, but the [`CHANGELOG.md`](https://github.com/Necrom4/sbb-tui/blob/master/CHANGELOG.md) clearly documents them.

## 🚀 Usage

1. Run `sbb-tui`
2. Input **departure** and **arrival** locations (navigate with `tab`).
3. Add optional information such as **date**, **time**, and **whether** those are for departure or arrival.
4. Press `Enter` to view the results (navigate with arrows).

## ⚙️ Configuration

Add your optional config at `$HOME/.config/sbb-tui/config.yaml`

> [!TIP]
> OS default config paths are also supported (such as `~/Library/Application Support/sbb-tui/config.yaml` in macOS)

```yaml
# default configuration
ui:
  nerdfont: true
  theme:
    text:            "#FFFFFF"
    textMuted:       "#888888"
    error:           "#D82E20"
    warning:         "#D82E20"
    borderFocused:   "#D82E20"
    borderUnfocused: "#484848"
    keysFg:          "#FFFFFF"
    keysBg:          "#484848"
    vehicleFg:       "#FFFFFF"
    vehicleBg:       "#2E3279"
    badgeModelFg:    "#FFFFFF"
    badgeModelBg:    "#D82E20"
    badgeCompanyFg:  "#484848"
    badgeCompanyBg:  "#FFFFFF"
    logo:            "#FFFFFF"
```

> More **themes** can be found at [`docs/themes.md`](https://github.com/Necrom4/sbb-tui/blob/master/docs/themes.md)

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
      --nerdfont      Use Nerd Font icons (`=false` fallback to Unicode)
      --time string   Pre-fill time (HH:MM)
      --to string     Pre-fill arrival station
  -v, --version       Print version and exit
```

## ❓ Why

> I travel 4 hours a day and often work from the train, passing through remote regions of Switzerland where loading the SBB website alone can take up to an entire minute before I've even had a chance to search for my next connection (back when I started building this, my cellular data plan provided me with 500kb/s at most in the city). On top of that, I love open source and disliked handing data over to corporations unnecessarily, so I loved the idea of a faster solution that only fetched the necessary data, the one I requested. One day, while exploring the idea of building my first TUI, I stumbled upon the incredible [Swiss public transport API](https://transport.opendata.ch/docs.html) and I knew I'd found the perfect occasion to start!

## 🍻 HELP WANTED!

**SBB-TUI** is in constant improvement thanks to the work of many volunteers passionate about this little tool.
This project is a very good playground for those who want to learn to collaborate in open source projects and improve the way they architecture their code and commit it. Whether you're a novice or a professional, don't be scared to hop on the train!

There's plenty of [Issues](https://github.com/Necrom4/sbb-tui/issues) of different complexity levels, many tagged with `help wanted` or `good first issue`.
So if you're up for the challenge, read [`CONTRIBUTING.md`](https://github.com/Necrom4/sbb-tui/blob/master/CONTRIBUTING.md), create your PR, and start coding!

## Star History

<a href="https://www.star-history.com/?repos=necrom4%2Fsbb-tui&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/image?repos=necrom4/sbb-tui&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/image?repos=necrom4/sbb-tui&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/image?repos=necrom4/sbb-tui&type=date&legend=top-left" />
 </picture>
</a>
