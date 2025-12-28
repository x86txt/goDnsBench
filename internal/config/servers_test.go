package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestServer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		server  Server
		wantErr bool
	}{
		{
			name:    "Valid server with DNS",
			server:  Server{Name: "Test Server", DNS: "1.1.1.1"},
			wantErr: false,
		},
		{
			name:    "Valid server with DoH",
			server:  Server{Name: "Test Server", DoH: "https://example.com/dns-query"},
			wantErr: false,
		},
		{
			name:    "Valid server with multiple protocols",
			server:  Server{Name: "Test Server", DNS: "1.1.1.1", DoH: "https://example.com/dns-query"},
			wantErr: false,
		},
		{
			name:    "Empty name",
			server:  Server{Name: "", DNS: "1.1.1.1"},
			wantErr: true,
		},
		{
			name:    "No protocols",
			server:  Server{Name: "Test Server"},
			wantErr: true,
		},
		{
			name:    "Empty name and no protocols",
			server:  Server{Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadServersFromJSON(t *testing.T) {
	// Create a temporary JSON file
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "servers.json")

	jsonContent := `{
  "servers": [
    {
      "name": "Test Server 1",
      "dns": "1.1.1.1",
      "doh": "https://example.com/dns-query"
    },
    {
      "name": "Test Server 2",
      "dns": "8.8.8.8"
    }
  ]
}`

	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	servers, err := LoadServersFromJSON(jsonFile)
	if err != nil {
		t.Fatalf("LoadServersFromJSON() error = %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(servers))
	}

	if servers[0].Name != "Test Server 1" {
		t.Errorf("Expected first server name='Test Server 1', got %s", servers[0].Name)
	}
	if servers[0].DNS != "1.1.1.1" {
		t.Errorf("Expected first server DNS='1.1.1.1', got %s", servers[0].DNS)
	}

	if servers[1].Name != "Test Server 2" {
		t.Errorf("Expected second server name='Test Server 2', got %s", servers[1].Name)
	}
}

func TestLoadServersFromJSON_InvalidFile(t *testing.T) {
	_, err := LoadServersFromJSON("/nonexistent/file.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestLoadServersFromJSON_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(jsonFile, []byte("invalid json content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := LoadServersFromJSON(jsonFile)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestLoadServersFromJSON_InvalidServer(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "invalid_servers.json")

	jsonContent := `{
  "servers": [
    {
      "name": "",
      "dns": "1.1.1.1"
    }
  ]
}`

	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	_, err := LoadServersFromJSON(jsonFile)
	if err == nil {
		t.Error("Expected error for invalid server (empty name)")
	}
}

func TestSaveServersToJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "output.json")

	servers := []Server{
		{Name: "Server 1", DNS: "1.1.1.1", DoH: "https://example.com/dns-query"},
		{Name: "Server 2", DNS: "8.8.8.8"},
	}

	if err := SaveServersToJSON(jsonFile, servers); err != nil {
		t.Fatalf("SaveServersToJSON() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Load and verify
	loadedServers, err := LoadServersFromJSON(jsonFile)
	if err != nil {
		t.Fatalf("Failed to load saved servers: %v", err)
	}

	if len(loadedServers) != len(servers) {
		t.Errorf("Expected %d servers, got %d", len(servers), len(loadedServers))
	}

	if loadedServers[0].Name != servers[0].Name {
		t.Errorf("Server name mismatch: expected %s, got %s", servers[0].Name, loadedServers[0].Name)
	}
}

func TestDefaultServers(t *testing.T) {
	servers := DefaultServers()

	if len(servers) == 0 {
		t.Error("Expected at least one default server")
	}

	// Verify all default servers are valid
	for _, server := range servers {
		if err := server.Validate(); err != nil {
			t.Errorf("Default server '%s' is invalid: %v", server.Name, err)
		}
	}
}
