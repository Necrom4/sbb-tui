// Package config handles application configuration and theming.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds CLI flag values to pre-fill the TUI form.
type Config struct {
	From           string
	To             string
	Date           string
	Time           string
	IsArrivalTime  bool
	NerdFont       bool
	Theme          Theme
	CurrentVersion string
}

// UIConfig groups all UI-related settings.
type UIConfig struct {
	NerdFont *bool `yaml:"nerdfont"`
	Theme    Theme `yaml:"theme"`
}

type fileConfig struct {
	UI UIConfig `yaml:"ui"`
}

// Theme defines color values for the TUI appearance.
type Theme struct {
	Text              string `yaml:"text"`
	TextMuted         string `yaml:"textMuted"`
	Error             string `yaml:"error"`
	Warning           string `yaml:"warning"`
	BorderFocused     string `yaml:"borderFocused"`
	BorderUnfocused   string `yaml:"borderUnfocused"`
	BadgeKeyFg        string `yaml:"badgeKeyFg"`
	BadgeKeyBg        string `yaml:"badgeKeyBg"`
	BadgeVehicleFg    string `yaml:"badgeVehicleFg"`
	BadgeVehicleBg    string `yaml:"badgeVehicleBg"`
	BadgeBadgeModelFg string `yaml:"badgeModelFg"`
	BadgeModelBg      string `yaml:"badgeModelBg"`
	BadgeCompanyFg    string `yaml:"badgeCompanyFg"`
	BadgeCompanyBg    string `yaml:"badgeCompanyBg"`
	Logo              string `yaml:"logo"`
}

// DefaultTheme returns the SBB brand color scheme.
func DefaultTheme() Theme {
	return Theme{
		Text:              "#FFFFFF",
		TextMuted:         "#888888",
		Error:             "#D82E20",
		Warning:           "#D82E20",
		BorderFocused:     "#D82E20",
		BorderUnfocused:   "#484848",
		BadgeKeyFg:        "#FFFFFF",
		BadgeKeyBg:        "#484848",
		BadgeVehicleFg:    "#FFFFFF",
		BadgeVehicleBg:    "#2E3279",
		BadgeBadgeModelFg: "#FFFFFF",
		BadgeModelBg:      "#D82E20",
		BadgeCompanyFg:    "#484848",
		BadgeCompanyBg:    "#FFFFFF",
		Logo:              "#FFFFFF",
	}
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving config path: %w", err)
	}

	// Prefer $HOME/.config/
	primary := filepath.Join(home, ".config", "sbb-tui", "config.yaml")
	if _, err := os.Stat(primary); err == nil {
		return primary, nil
	}

	// Fall back to OS default config path
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return primary, nil
	}
	return filepath.Join(cfgDir, "sbb-tui", "config.yaml"), nil
}

// loadFile reads and parses the config file, returning a raw fileConfig.
func loadFile() (fileConfig, error) {
	path, err := configFilePath()
	if err != nil {
		return fileConfig{}, fmt.Errorf("loading config: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fileConfig{}, nil
		}
		return fileConfig{}, fmt.Errorf("loading config: reading %s: %w", path, err)
	}

	var cfg fileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("loading config: parsing %s: %w", path, err)
	}

	return cfg, nil
}

// LoadConfig reads the config file and returns a Config with defaults merged.
func LoadConfig() (Config, error) {
	result := Config{
		NerdFont: true,
		Theme:    DefaultTheme(),
	}

	fc, err := loadFile()
	if err != nil {
		return result, err
	}

	if fc.UI.NerdFont != nil {
		result.NerdFont = *fc.UI.NerdFont
	}

	// NOTE: update mergeTheme when adding new Theme fields.
	result.Theme = mergeTheme(result.Theme, fc.UI.Theme)
	return result, nil
}

func mergeTheme(base Theme, override Theme) Theme {
	if override.Text != "" {
		base.Text = override.Text
	}
	if override.TextMuted != "" {
		base.TextMuted = override.TextMuted
	}
	if override.BorderFocused != "" {
		base.BorderFocused = override.BorderFocused
	}
	if override.BorderUnfocused != "" {
		base.BorderUnfocused = override.BorderUnfocused
	}
	if override.Warning != "" {
		base.Warning = override.Warning
	}
	if override.BadgeKeyFg != "" {
		base.BadgeKeyFg = override.BadgeKeyFg
	}
	if override.BadgeKeyBg != "" {
		base.BadgeKeyBg = override.BadgeKeyBg
	}
	if override.BadgeVehicleFg != "" {
		base.BadgeVehicleFg = override.BadgeVehicleFg
	}
	if override.BadgeVehicleBg != "" {
		base.BadgeVehicleBg = override.BadgeVehicleBg
	}
	if override.BadgeBadgeModelFg != "" {
		base.BadgeBadgeModelFg = override.BadgeBadgeModelFg
	}
	if override.BadgeModelBg != "" {
		base.BadgeModelBg = override.BadgeModelBg
	}
	if override.BadgeCompanyFg != "" {
		base.BadgeCompanyFg = override.BadgeCompanyFg
	}
	if override.BadgeCompanyBg != "" {
		base.BadgeCompanyBg = override.BadgeCompanyBg
	}
	if override.Logo != "" {
		base.Logo = override.Logo
	}

	return base
}
