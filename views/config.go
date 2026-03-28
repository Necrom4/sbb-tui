package views

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Theme struct {
	Primary    string `yaml:"primary"`
	Text       string `yaml:"text"`
	Border     string `yaml:"border"`
	Muted      string `yaml:"muted"`
	Vehicle    string `yaml:"vehicle"`
	Company    string `yaml:"company"`
	Warning    string `yaml:"warning"`
	KeysBg     string `yaml:"keysBg"`
}

type FileConfig struct {
	Theme Theme `yaml:"theme"`
}

func DefaultTheme() Theme {
	return Theme{
		Primary:    "#D82E20",
		Text:       "#FFFFFF",
		Border:     "#862010",
		Muted:      "#888888",
		Vehicle:    "#315086",
		Company:    "#484848",
		Warning:    "#3A7446",
		KeysBg:     "#484848",
	}
}

func configFilePath() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, "sbb-tui", "config.yaml"), nil
}

func LoadThemeConfig() (Theme, error) {
	theme := DefaultTheme()

	path, err := configFilePath()
	if err != nil {
		return theme, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return theme, nil
		}
		return theme, err
	}

	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return theme, err
	}

	theme = mergeTheme(theme, cfg.Theme)
	return theme, nil
}

func mergeTheme(base Theme, override Theme) Theme {
	if override.Primary != "" {
		base.Primary = override.Primary
	}
	if override.Text != "" {
		base.Text = override.Text
	}
	if override.Border != "" {
		base.Border = override.Border
	}
	if override.Muted != "" {
		base.Muted = override.Muted
	}
	if override.Vehicle != "" {
		base.Vehicle = override.Vehicle
	}
	if override.Company != "" {
		base.Company = override.Company
	}
	if override.Warning != "" {
		base.Warning = override.Warning
	}
	if override.KeysBg != "" {
		base.KeysBg = override.KeysBg
	}

	return base
}
