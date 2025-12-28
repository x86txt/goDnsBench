package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Settings represents user preferences
type Settings struct {
	QueryTimeout      time.Duration `json:"queryTimeout"`
	MaxConcurrent     int           `json:"maxConcurrent"`
	ServerListURL     string        `json:"serverListUrl"`
	LastServerUpdate  time.Time     `json:"lastServerUpdate"`
	CustomServers     []Server      `json:"customServers"`
	EnabledProtocols  []string      `json:"enabledProtocols"`
	SelectedDomains   TestDomains   `json:"selectedDomains"`
}

// TestDomains represents the test domains configuration
type TestDomains struct {
	A      []string `json:"a"`
	MX     []string `json:"mx"`
	TXT    []string `json:"txt"`
	DNSSEC []string `json:"dnssec"`
}

// DefaultSettings returns default settings
func DefaultSettings() Settings {
	return Settings{
		QueryTimeout:  1 * time.Second,
		MaxConcurrent: 10,
		ServerListURL: "https://raw.githubusercontent.com/x86txt/goDnsBench/main/servers.json",
		EnabledProtocols: []string{"DNS", "DoH", "DoT", "DoQ"},
		SelectedDomains: TestDomains{
			A:      []string{"google.com", "cloudflare.com", "amazon.com", "microsoft.com"},
			MX:     []string{"gmail.com", "microsoft.com"},
			TXT:    []string{"_dmarc.google.com", "google.com"},
			DNSSEC: []string{"cloudflare.com", "google.com"},
		},
	}
}

// GetConfigDir returns the application configuration directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "goDnsBench")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

// LoadSettings loads settings from the config file
func LoadSettings() (Settings, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return DefaultSettings(), err
	}

	settingsPath := filepath.Join(configDir, "settings.json")

	// If file doesn't exist, return defaults
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return DefaultSettings(), nil
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return DefaultSettings(), fmt.Errorf("failed to read settings: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return DefaultSettings(), fmt.Errorf("failed to parse settings: %w", err)
	}

	return settings, nil
}

// SaveSettings saves settings to the config file
func (s Settings) Save() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(configDir, "settings.json")

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

// NeedsServerUpdate checks if the server list needs to be updated
func (s Settings) NeedsServerUpdate() bool {
	// Update if never updated or older than 7 days
	if s.LastServerUpdate.IsZero() {
		return true
	}

	return time.Since(s.LastServerUpdate) > 7*24*time.Hour
}
