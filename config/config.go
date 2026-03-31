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
	NerdFont *bool `yaml:"nerdFont"`
	Theme    Theme `yaml:"theme"`
}

type fileConfig struct {
	UI UIConfig `yaml:"ui"`
}

// Theme defines color values for the TUI appearance.
type Theme struct {
	Text           string `yaml:"text"`
	ErrorText      string `yaml:"errorText"`
	GhostText      string `yaml:"ghostText"`
	ActiveBorder   string `yaml:"activeBorder"`
	InactiveBorder string `yaml:"inactiveBorder"`
	WarningFlag    string `yaml:"warningFlag"`
	KeysFg         string `yaml:"keysFg"`
	KeysBg         string `yaml:"keysBg"`
	VehicleFg      string `yaml:"vehicleFg"`
	VehicleBg      string `yaml:"vehicleBg"`
	ModelFg        string `yaml:"modelFg"`
	ModelBg        string `yaml:"modelBg"`
	CompanyFg      string `yaml:"companyFg"`
	CompanyBg      string `yaml:"companyBg"`
	Logo           string `yaml:"logo"`
}

// DefaultTheme returns the SBB brand color scheme.
func DefaultTheme() Theme {
	return Theme{
		Text:           "#FFFFFF",
		ErrorText:      "#D82E20",
		GhostText:      "#888888",
		ActiveBorder:   "#D82E20",
		InactiveBorder: "#484848",
		WarningFlag:    "#D82E20",
		KeysFg:         "#FFFFFF",
		KeysBg:         "#484848",
		VehicleFg:      "#FFFFFF",
		VehicleBg:      "#2E3279",
		ModelFg:        "#FFFFFF",
		ModelBg:        "#D82E20",
		CompanyFg:      "#484848",
		CompanyBg:      "#FFFFFF",
		Logo:           "#FFFFFF",
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
	if override.GhostText != "" {
		base.GhostText = override.GhostText
	}
	if override.ActiveBorder != "" {
		base.ActiveBorder = override.ActiveBorder
	}
	if override.InactiveBorder != "" {
		base.InactiveBorder = override.InactiveBorder
	}
	if override.WarningFlag != "" {
		base.WarningFlag = override.WarningFlag
	}
	if override.KeysFg != "" {
		base.KeysFg = override.KeysFg
	}
	if override.KeysBg != "" {
		base.KeysBg = override.KeysBg
	}
	if override.VehicleFg != "" {
		base.VehicleFg = override.VehicleFg
	}
	if override.VehicleBg != "" {
		base.VehicleBg = override.VehicleBg
	}
	if override.ModelFg != "" {
		base.ModelFg = override.ModelFg
	}
	if override.ModelBg != "" {
		base.ModelBg = override.ModelBg
	}
	if override.CompanyFg != "" {
		base.CompanyFg = override.CompanyFg
	}
	if override.CompanyBg != "" {
		base.CompanyBg = override.CompanyBg
	}
	if override.Logo != "" {
		base.Logo = override.Logo
	}

	return base
}
