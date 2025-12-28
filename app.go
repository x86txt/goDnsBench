package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

	// Server selection is managed by the GUI and used as a filter at run time.
	hasExplicitServerSelection bool
	selectedServerNames        []string
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
func (a *App) GetServers() []ServerDTO {
	out := make([]ServerDTO, 0, len(a.servers))
	for _, s := range a.servers {
		out = append(out, serverToDTO(s))
	}
	return out
}

// GetSettings returns the current settings
func (a *App) GetSettings() SettingsDTO {
	return settingsToDTO(a.settings)
}

// UpdateSettings updates and saves the settings
func (a *App) UpdateSettings(settings SettingsDTO) error {
	updated := mergeSettingsFromDTO(a.settings, settings)
	a.settings = updated
	return updated.Save()
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

// LoadServersFromContent loads servers from file content (for GUI file input)
func (a *App) LoadServersFromContent(content string, fileType string) error {
	var servers []config.Server
	var err error

	if fileType == "json" || fileType == ".json" {
		var serverList config.ServerList
		if err := json.Unmarshal([]byte(content), &serverList); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
		servers = serverList.Servers
	} else if fileType == "csv" || fileType == ".csv" {
		servers, err = config.ParseCSVFromString(content)
		if err != nil {
			return fmt.Errorf("failed to parse CSV: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported file format (must be .json or .csv)")
	}

	// Validate all servers
	for _, server := range servers {
		if err := server.Validate(); err != nil {
			return fmt.Errorf("invalid server %s: %w", server.Name, err)
		}
	}

	a.servers = servers
	return nil
}

// AddServer adds a custom server to the list
func (a *App) AddServer(server ServerDTO) error {
	srv := serverDTOToConfig(server)
	if err := srv.Validate(); err != nil {
		return err
	}

	a.servers = append(a.servers, srv)
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

// SetSelectedServers sets the active server selection (by server name).
// Passing an empty slice means the user explicitly selected none.
func (a *App) SetSelectedServers(names []string) error {
	a.hasExplicitServerSelection = true

	// Normalize and de-dupe.
	desired := make(map[string]struct{}, len(names))
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		desired[n] = struct{}{}
	}

	// Validate against current server list and preserve stable ordering.
	existing := make(map[string]struct{}, len(a.servers))
	for _, s := range a.servers {
		existing[s.Name] = struct{}{}
	}
	for n := range desired {
		if _, ok := existing[n]; !ok {
			return fmt.Errorf("unknown server: %s", n)
		}
	}

	selected := make([]string, 0, len(desired))
	for _, s := range a.servers {
		if _, ok := desired[s.Name]; ok {
			selected = append(selected, s.Name)
		}
	}

	a.selectedServerNames = selected
	return nil
}

// GetSelectedServers returns the currently selected server names.
func (a *App) GetSelectedServers() []string {
	return append([]string(nil), a.selectedServerNames...)
}

func filterServersByName(all []config.Server, names []string) []config.Server {
	if len(names) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(names))
	for _, n := range names {
		set[n] = struct{}{}
	}

	out := make([]config.Server, 0, len(names))
	for _, s := range all {
		if _, ok := set[s.Name]; ok {
			out = append(out, s)
		}
	}
	return out
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

	if len(protoList) == 0 {
		return fmt.Errorf("no protocols selected")
	}

	serversToRun := a.servers
	if a.hasExplicitServerSelection {
		if len(a.selectedServerNames) == 0 {
			return fmt.Errorf("no servers selected")
		}
		serversToRun = filterServersByName(a.servers, a.selectedServerNames)
	}
	if len(serversToRun) == 0 {
		return fmt.Errorf("no servers available to benchmark")
	}

	// Create benchmark config
	benchConfig := benchmark.BenchmarkConfig{
		Servers:       serversToRun,
		Protocols:     protoList,
		Timeout:       a.settings.QueryTimeout,
		MaxConcurrent: a.settings.MaxConcurrent,
		TestDomains: benchmark.TestDomains{
			A:      a.settings.SelectedDomains.A,
			MX:     a.settings.SelectedDomains.MX,
			TXT:    a.settings.SelectedDomains.TXT,
			DNSSEC: a.settings.SelectedDomains.DNSSEC,
		},
		ProgressCallback: func(p benchmark.Progress) {
			if a.ctx == nil {
				return
			}
			runtime.EventsEmit(a.ctx, EventBenchmarkProgress, BenchmarkProgressDTO{
				CurrentServer:   p.CurrentServer,
				CurrentProtocol: p.CurrentProtocol.String(),
				CompletedTests:  p.CompletedTests,
				TotalTests:      p.TotalTests,
				Percentage:      p.Percentage,
			})
		},
	}

	// Create runner
	runner := benchmark.NewRunner(benchConfig)

	// Run benchmark
	results, err := runner.Run(a.ctx)
	if err != nil {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, EventBenchmarkError, map[string]string{
				"error": err.Error(),
			})
		}
		return err
	}

	a.results = results

	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, EventBenchmarkDone, benchmarkResultsToDTO(results))
	}
	return nil
}

// GetResults returns the last benchmark results
func (a *App) GetResults() *BenchmarkResultsDTO {
	return benchmarkResultsToDTO(a.results)
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
func (a *App) CheckServerCapabilities(server ServerDTO) map[string]bool {
	timeout := 2 * time.Second
	srv := serverDTOToConfig(server)
	caps := dns.CheckCapabilities(srv.DNS, srv.DoH, srv.DoT, srv.DoQ, timeout)

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
