package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCSVFromString(t *testing.T) {
	csvContent := `name,dns,doh,dot,doq
Server 1,1.1.1.1,https://example.com/dns-query,1.1.1.1:853,1.1.1.1:8853
Server 2,8.8.8.8,https://dns.google/dns-query,8.8.8.8:853,
Server 3,9.9.9.9,,,`

	servers, err := ParseCSVFromString(csvContent)
	if err != nil {
		t.Fatalf("ParseCSVFromString() error = %v", err)
	}

	if len(servers) != 3 {
		t.Errorf("Expected 3 servers, got %d", len(servers))
	}

	if servers[0].Name != "Server 1" {
		t.Errorf("Expected first server name='Server 1', got %s", servers[0].Name)
	}
	if servers[0].DNS != "1.1.1.1" {
		t.Errorf("Expected first server DNS='1.1.1.1', got %s", servers[0].DNS)
	}

	if servers[1].Name != "Server 2" {
		t.Errorf("Expected second server name='Server 2', got %s", servers[1].Name)
	}

	if servers[2].Name != "Server 3" {
		t.Errorf("Expected third server name='Server 3', got %s", servers[2].Name)
	}
}

func TestParseCSVFromString_MinimalColumns(t *testing.T) {
	csvContent := `name,dns
Server 1,1.1.1.1`

	servers, err := ParseCSVFromString(csvContent)
	if err != nil {
		t.Fatalf("ParseCSVFromString() error = %v", err)
	}

	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}
}

func TestParseCSVFromString_InvalidServer(t *testing.T) {
	csvContent := `name,dns
,1.1.1.1`

	_, err := ParseCSVFromString(csvContent)
	if err == nil {
		t.Error("Expected error for invalid server (empty name)")
	}
}

func TestParseCSVFromString_NoProtocols(t *testing.T) {
	csvContent := `name,dns
Server 1,`

	_, err := ParseCSVFromString(csvContent)
	if err == nil {
		t.Error("Expected error for server with no protocols")
	}
}

func TestLoadServersFromCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "servers.csv")

	csvContent := `name,dns,doh
Server 1,1.1.1.1,https://example.com/dns-query
Server 2,8.8.8.8,https://dns.google/dns-query`

	if err := os.WriteFile(csvFile, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	servers, err := LoadServersFromCSV(csvFile)
	if err != nil {
		t.Fatalf("LoadServersFromCSV() error = %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(servers))
	}
}

func TestLoadServersFromCSV_InvalidFile(t *testing.T) {
	_, err := LoadServersFromCSV("/nonexistent/file.csv")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestSaveServersToCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "output.csv")

	servers := []Server{
		{Name: "Server 1", DNS: "1.1.1.1", DoH: "https://example.com/dns-query"},
		{Name: "Server 2", DNS: "8.8.8.8"},
	}

	if err := SaveServersToCSV(csvFile, servers); err != nil {
		t.Fatalf("SaveServersToCSV() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Load and verify
	loadedServers, err := LoadServersFromCSV(csvFile)
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

func TestSaveServersToCSV_EmptyList(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "empty.csv")

	if err := SaveServersToCSV(csvFile, []Server{}); err != nil {
		t.Fatalf("SaveServersToCSV() error = %v", err)
	}

	// File should exist with just header
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}
