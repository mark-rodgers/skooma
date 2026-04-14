// Package config provides functions for managing the Skooma configuration file.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Template represents a project template
type Template struct {
	Description string `json:"description"`
	Repo        string `json:"repo"`
	Author      string `json:"author"`
}

// Config represents the Skooma configuration
type Config struct {
	Templates map[string]Template `json:"templates"`
}

// GetConfig retrieves the config object, creating default config if it doesn't exist
func GetConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// Ensure the skooma config directory exists
	skoomaDir := filepath.Join(configDir, "skooma")
	if _, err := os.Stat(skoomaDir); os.IsNotExist(err) {
		err = os.MkdirAll(skoomaDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// If config file exists, read and unmarshal it
	configPath := filepath.Join(skoomaDir, "config.json")
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		var config Config
		err = json.Unmarshal(data, &config)
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	// File doesn't exist, create it with default config
	defaultConfig := &Config{
		Templates: map[string]Template{
			"default": {
				Description: "A default template with Go, React, Tailwind, and Vite",
				Repo:        "github.com/mark-rodgers/skooma-default-template",
				Author:      "Mark Rodgers <mark@marknrodgers.com>",
			},
		},
	}

	// Write default config to file
	file, err := os.Create(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	encoder.SetEscapeHTML(false)

	err = encoder.Encode(defaultConfig)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}
