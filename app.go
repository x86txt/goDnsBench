package main

import (
	"context"
	"fmt"
	"time"

	"github.com/x86txt/goDnsBench/internal/benchmark"
	"github.com/x86txt/goDnsBench/internal/config"
	"github.com/x86txt/goDnsBench/internal/dns"
	"github.com/x86txt/goDnsBench/internal/export"
)

// App struct
type App struct {
	ctx      context.Context
	settings config.Settings
	servers  []config.Server
	results  *benchmark.BenchmarkResults
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load settings
	settings, err := config.LoadSettings()
	if err != nil {
		settings = config.DefaultSettings()
	}
	a.settings = settings

	// Load default servers
	a.servers = config.DefaultServers()
}

// GetServers returns the current server list
func (a *App) GetServers() []config.Server {
	return a.servers
}

// GetSettings returns the current settings
func (a *App) GetSettings() config.Settings {
	return a.settings
}

// UpdateSettings updates and saves the settings
func (a *App) UpdateSettings(settings config.Settings) error {
	a.settings = settings
	return settings.Save()
}

// LoadServersFromFile loads servers from a JSON or CSV file
func (a *App) LoadServersFromFile(filepath string) error {
	var servers []config.Server
	var err error

	// Determine file type by extension
	if len(filepath) >= 5 && filepath[len(filepath)-5:] == ".json" {
		servers, err = config.LoadServersFromJSON(filepath)
	} else if len(filepath) >= 4 && filepath[len(filepath)-4:] == ".csv" {
		servers, err = config.LoadServersFromCSV(filepath)
	} else {
		return fmt.Errorf("unsupported file format (must be .json or .csv)")
	}

	if err != nil {
		return err
	}

	a.servers = servers
	return nil
}

// AddServer adds a custom server to the list
func (a *App) AddServer(server config.Server) error {
	if err := server.Validate(); err != nil {
		return err
	}

	a.servers = append(a.servers, server)
	return nil
}

// RemoveServer removes a server from the list by name
func (a *App) RemoveServer(name string) {
	filtered := make([]config.Server, 0)
	for _, server := range a.servers {
		if server.Name != name {
			filtered = append(filtered, server)
		}
	}
	a.servers = filtered
}

// ResetToDefaults resets the server list to defaults
func (a *App) ResetToDefaults() {
	a.servers = config.DefaultServers()
}

// BenchmarkProgress represents progress updates during benchmarking
type BenchmarkProgress struct {
	CurrentServer  string  `json:"currentServer"`
	CurrentProtocol string  `json:"currentProtocol"`
	CompletedTests int     `json:"completedTests"`
	TotalTests     int     `json:"totalTests"`
	Percentage     float64 `json:"percentage"`
}

// RunBenchmark executes the benchmark with the selected servers and protocols
func (a *App) RunBenchmark(protocols []string) error {
	// Convert protocol strings to Protocol types
	var protoList []dns.Protocol
	for _, p := range protocols {
		switch p {
		case "DNS":
			protoList = append(protoList, dns.ProtocolDNS)
		case "DoH":
			protoList = append(protoList, dns.ProtocolDoH)
		case "DoT":
			protoList = append(protoList, dns.ProtocolDoT)
		case "DoQ":
			protoList = append(protoList, dns.ProtocolDoQ)
		}
	}

	// Create benchmark config
	benchConfig := benchmark.BenchmarkConfig{
		Servers:       a.servers,
		Protocols:     protoList,
		Timeout:       a.settings.QueryTimeout,
		MaxConcurrent: a.settings.MaxConcurrent,
		TestDomains: benchmark.TestDomains{
			A:      a.settings.SelectedDomains.A,
			MX:     a.settings.SelectedDomains.MX,
			TXT:    a.settings.SelectedDomains.TXT,
			DNSSEC: a.settings.SelectedDomains.DNSSEC,
		},
	}

	// Create runner
	runner := benchmark.NewRunner(benchConfig)

	// Run benchmark
	results, err := runner.Run(a.ctx)
	if err != nil {
		return err
	}

	a.results = results
	return nil
}

// GetResults returns the last benchmark results
func (a *App) GetResults() *benchmark.BenchmarkResults {
	return a.results
}

// ExportResultsJSON exports results to JSON file
func (a *App) ExportResultsJSON(filepath string) error {
	if a.results == nil {
		return fmt.Errorf("no results to export")
	}
	return export.ExportResultsJSON(a.results, filepath)
}

// ExportResultsCSV exports results to CSV file
func (a *App) ExportResultsCSV(filepath string) error {
	if a.results == nil {
		return fmt.Errorf("no results to export")
	}
	return export.ExportResultsCSV(a.results, filepath)
}

// CheckServerCapabilities checks which protocols a server supports
func (a *App) CheckServerCapabilities(server config.Server) map[string]bool {
	timeout := 2 * time.Second
	caps := dns.CheckCapabilities(server.DNS, server.DoH, server.DoT, server.DoQ, timeout)

	return map[string]bool{
		"DNS": caps.DNS,
		"DoH": caps.DoH,
		"DoT": caps.DoT,
		"DoQ": caps.DoQ,
	}
}

// RefreshServerList fetches the latest server list from the configured URL
func (a *App) RefreshServerList() error {
	if a.settings.ServerListURL == "" {
		return fmt.Errorf("no server list URL configured")
	}

	servers, err := config.FetchServersFromURL(a.settings.ServerListURL, 10*time.Second)
	if err != nil {
		return err
	}

	a.servers = servers
	a.settings.LastServerUpdate = time.Now()
	return a.settings.Save()
}
