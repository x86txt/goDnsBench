package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultSettings(t *testing.T) {
	settings := DefaultSettings()

	if settings.QueryTimeout != 1*time.Second {
		t.Errorf("Expected QueryTimeout=1s, got %v", settings.QueryTimeout)
	}
	if settings.MaxConcurrent != 10 {
		t.Errorf("Expected MaxConcurrent=10, got %d", settings.MaxConcurrent)
	}
	if len(settings.EnabledProtocols) != 4 {
		t.Errorf("Expected 4 enabled protocols, got %d", len(settings.EnabledProtocols))
	}
	if len(settings.SelectedDomains.A) == 0 {
		t.Error("Expected default A record domains")
	}
	if len(settings.SelectedDomains.MX) == 0 {
		t.Error("Expected default MX record domains")
	}
}

func TestSettings_SaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	originalSettings := Settings{
		QueryTimeout:  2 * time.Second,
		MaxConcurrent: 5,
		ServerListURL: "https://example.com/servers.json",
		EnabledProtocols: []string{"DNS", "DoH"},
		SelectedDomains: TestDomains{
			A:      []string{"test.com"},
			MX:     []string{"mail.test.com"},
			TXT:    []string{"txt.test.com"},
			DNSSEC: []string{"dnssec.test.com"},
		},
	}

	// Manually save to temp file
	data, err := json.MarshalIndent(originalSettings, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal settings: %v", err)
	}
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatalf("Failed to write settings file: %v", err)
	}

	// Load settings from temp file
	data, err = os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("Failed to read settings file: %v", err)
	}

	var loadedSettings Settings
	if err := json.Unmarshal(data, &loadedSettings); err != nil {
		t.Fatalf("Failed to unmarshal settings: %v", err)
	}

	// Verify loaded settings match
	if loadedSettings.QueryTimeout != originalSettings.QueryTimeout {
		t.Errorf("QueryTimeout mismatch: expected %v, got %v", originalSettings.QueryTimeout, loadedSettings.QueryTimeout)
	}
	if loadedSettings.MaxConcurrent != originalSettings.MaxConcurrent {
		t.Errorf("MaxConcurrent mismatch: expected %d, got %d", originalSettings.MaxConcurrent, loadedSettings.MaxConcurrent)
	}
	if loadedSettings.ServerListURL != originalSettings.ServerListURL {
		t.Errorf("ServerListURL mismatch: expected %s, got %s", originalSettings.ServerListURL, loadedSettings.ServerListURL)
	}
	if len(loadedSettings.EnabledProtocols) != len(originalSettings.EnabledProtocols) {
		t.Errorf("EnabledProtocols length mismatch: expected %d, got %d", len(originalSettings.EnabledProtocols), len(loadedSettings.EnabledProtocols))
	}
}

func TestSettings_NeedsServerUpdate(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		expected bool
	}{
		{
			name:     "Never updated",
			settings: Settings{LastServerUpdate: time.Time{}},
			expected: true,
		},
		{
			name:     "Recently updated",
			settings: Settings{LastServerUpdate: time.Now().Add(-1 * time.Hour)},
			expected: false,
		},
		{
			name:     "Updated 7 days ago (just under threshold)",
			settings: Settings{LastServerUpdate: time.Now().Add(-7*24*time.Hour + time.Hour)},
			expected: false,
		},
		{
			name:     "Updated 8 days ago",
			settings: Settings{LastServerUpdate: time.Now().Add(-8 * 24 * time.Hour)},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.settings.NeedsServerUpdate()
			if result != tt.expected {
				t.Errorf("Expected NeedsServerUpdate=%v, got %v", tt.expected, result)
			}
		})
	}
}
