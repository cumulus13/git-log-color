package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type UIConfig struct {
	ColorEnabled        bool `json:"color_enabled"`
	IconsEnabled        bool `json:"icons_enabled"`
	AuthorColorsEnabled bool `json:"author_colors_enabled"`
}

type Config struct {
	Styles        map[string]string `json:"styles"`
	Icons         map[string]string `json:"icons"`
	UI            UIConfig          `json:"ui"`
	AuthorPalette []string          `json:"author_palette,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		Styles:        defaultStyles(),
		Icons:         defaultIcons(),
		AuthorPalette: defaultAuthorPalette(),
		UI: UIConfig{
			ColorEnabled:        true,
			IconsEnabled:        true,
			AuthorColorsEnabled: true,
		},
	}
}

func LoadConfig(configPath string, repoPath string) (Config, error) {
	cfg := DefaultConfig()

	path, err := findConfigPath(configPath, repoPath)
	if err != nil {
		return cfg, err
	}
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	var fileCfg Config
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}

	if len(fileCfg.Styles) > 0 {
		for key, value := range fileCfg.Styles {
			cfg.Styles[key] = value
		}
	}
	if len(fileCfg.Icons) > 0 {
		for key, value := range fileCfg.Icons {
			cfg.Icons[key] = value
		}
	}
	if len(fileCfg.AuthorPalette) > 0 {
		cfg.AuthorPalette = fileCfg.AuthorPalette
	}

	if fileCfg.UI.ColorEnabled || containsJSONKey(data, "color_enabled") {
		cfg.UI.ColorEnabled = fileCfg.UI.ColorEnabled
	}
	if fileCfg.UI.IconsEnabled || containsJSONKey(data, "icons_enabled") {
		cfg.UI.IconsEnabled = fileCfg.UI.IconsEnabled
	}
	if fileCfg.UI.AuthorColorsEnabled || containsJSONKey(data, "author_colors_enabled") {
		cfg.UI.AuthorColorsEnabled = fileCfg.UI.AuthorColorsEnabled
	}

	return cfg, nil
}

func findConfigPath(configPath string, repoPath string) (string, error) {
	if strings.TrimSpace(configPath) != "" {
		if filepath.IsAbs(configPath) {
			return configPath, nil
		}
		return filepath.Abs(filepath.Join(repoPath, configPath))
	}

	candidates := []string{
		filepath.Join(repoPath, ".git-log-color.json"),
		filepath.Join(repoPath, "git-log-color.json"),
		filepath.Join(repoPath, ".config", "git-log-color.json"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", nil
}

func containsJSONKey(data []byte, key string) bool {
	return strings.Contains(string(data), `"`+key+`"`)
}
