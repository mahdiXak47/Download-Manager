package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mahdiXak47/Download-Manager/internal/downloader"
)

type Config struct {
	DefaultQueue  string                `json:"default_queue"`
	MaxConcurrent int                   `json:"max_concurrent"`
	SavePath      string                `json:"save_path"`
	Downloads     []downloader.Download `json:"downloads"`
}

var defaultConfig = Config{
	DefaultQueue:  "default",
	MaxConcurrent: 3,
	SavePath:      "downloads",
}

const configFileName = "download-manager.json"

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".config", "download-manager", configFileName)
}

// LoadConfig loads the configuration from file or creates default if not exists
func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	// Try to read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config
			config := defaultConfig
			if err := SaveConfig(&config); err != nil {
				return nil, err
			}
			return &config, nil
		}
		return nil, err
	}

	// Parse existing config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetConfigPath(), data, 0644)
}
