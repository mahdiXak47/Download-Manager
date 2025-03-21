package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/mahdiXak47/Download-Manager/internal/downloader"
)

type QueueConfig struct {
	Name          string `json:"name"`
	MaxConcurrent int    `json:"max_concurrent"`
	StartTime     string `json:"start_time"`  // Format: "HH:MM"
	EndTime       string `json:"end_time"`    // Format: "HH:MM"
	SpeedLimit    int64  `json:"speed_limit"` // Bytes per second, 0 for unlimited
	Enabled       bool   `json:"enabled"`
	Path          string `json:"path"` // Download directory path for this queue
}

type Config struct {
	DefaultQueue string                `json:"default_queue"`
	SavePath     string                `json:"save_path"`
	Downloads    []downloader.Download `json:"downloads"`
	Queues       []QueueConfig         `json:"queues"`
}

var defaultConfig = Config{
	DefaultQueue: "default",
	SavePath:     "downloads",
	Queues: []QueueConfig{
		{
			Name:          "default",
			MaxConcurrent: 3,
			StartTime:     "00:00",
			EndTime:       "23:59",
			SpeedLimit:    0,
			Enabled:       true,
			Path:          "downloads/default",
		},
		{
			Name:          "night",
			MaxConcurrent: 5,
			StartTime:     "23:00",
			EndTime:       "06:00",
			SpeedLimit:    0,
			Enabled:       true,
			Path:          "downloads/night",
		},
	},
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

// IsTimeAllowed checks if downloads are allowed for a queue at the current time
func (q *QueueConfig) IsTimeAllowed() bool {
	if !q.Enabled {
		return false
	}

	now := time.Now()
	currentTime := now.Format("15:04")

	// Handle overnight windows (e.g., 23:00-06:00)
	if q.StartTime > q.EndTime {
		// If current time is after start OR before end, it's allowed
		return currentTime >= q.StartTime || currentTime <= q.EndTime
	}

	// Normal time window (e.g., 09:00-17:00)
	return currentTime >= q.StartTime && currentTime <= q.EndTime
}

// GetQueue returns a queue configuration by name
func (c *Config) GetQueue(name string) *QueueConfig {
	for i := range c.Queues {
		if c.Queues[i].Name == name {
			return &c.Queues[i]
		}
	}
	return nil
}
