package main

import (
	"time"

	"github.com/x86txt/goDnsBench/internal/benchmark"
	"github.com/x86txt/goDnsBench/internal/config"
)

const (
	EventBenchmarkProgress = "benchmark:progress"
	EventBenchmarkDone     = "benchmark:done"
	EventBenchmarkError    = "benchmark:error"
)

// ServerDTO is the GUI-safe representation of a DNS server configuration.
type ServerDTO struct {
	Name string `json:"name"`
	DNS  string `json:"dns,omitempty"`
	DoH  string `json:"doh,omitempty"`
	DoT  string `json:"dot,omitempty"`
	DoQ  string `json:"doq,omitempty"`
}

// TestDomainsDTO is the GUI-safe representation of domain configuration.
type TestDomainsDTO struct {
	A      []string `json:"a"`
	MX     []string `json:"mx"`
	TXT    []string `json:"txt"`
	DNSSEC []string `json:"dnssec"`
}

// SettingsDTO is the GUI-safe representation of user settings.
// Durations are expressed in milliseconds; times are ISO-8601 strings.
type SettingsDTO struct {
	QueryTimeoutMs   int64          `json:"queryTimeoutMs"`
	MaxConcurrent    int            `json:"maxConcurrent"`
	ServerListURL    string         `json:"serverListUrl"`
	LastServerUpdate string         `json:"lastServerUpdate,omitempty"`
	EnabledProtocols []string       `json:"enabledProtocols"`
	SelectedDomains  TestDomainsDTO `json:"selectedDomains"`
}

// BenchmarkProgressDTO is emitted during a benchmark run.
type BenchmarkProgressDTO struct {
	CurrentServer   string  `json:"currentServer"`
	CurrentProtocol string  `json:"currentProtocol"`
	CompletedTests  int     `json:"completedTests"`
	TotalTests      int     `json:"totalTests"`
	Percentage      float64 `json:"percentage"`
}

// MetricsDTO is a GUI-safe representation of benchmark metrics.
// Durations are expressed in milliseconds.
type MetricsDTO struct {
	MinMs       float64 `json:"minMs"`
	MaxMs       float64 `json:"maxMs"`
	MeanMs      float64 `json:"meanMs"`
	MedianMs    float64 `json:"medianMs"`
	P95Ms       float64 `json:"p95Ms"`
	P99Ms       float64 `json:"p99Ms"`
	Success     int     `json:"success"`
	Failed      int     `json:"failed"`
	Total       int     `json:"total"`
	SuccessRate float64 `json:"successRate"`
}

// ServerResultDTO is the GUI-safe representation of a server+protocol benchmark result.
type ServerResultDTO struct {
	ServerName string     `json:"serverName"`
	Protocol   string     `json:"protocol"`
	Metrics    MetricsDTO `json:"metrics"`
	DurationMs float64    `json:"durationMs"`
}

// BenchmarkResultsDTO is the GUI-safe representation of a benchmark run.
type BenchmarkResultsDTO struct {
	StartTime  string            `json:"startTime"`
	EndTime    string            `json:"endTime"`
	DurationMs float64           `json:"durationMs"`
	Results    []ServerResultDTO `json:"results"`
}

func serverToDTO(s config.Server) ServerDTO {
	return ServerDTO{
		Name: s.Name,
		DNS:  s.DNS,
		DoH:  s.DoH,
		DoT:  s.DoT,
		DoQ:  s.DoQ,
	}
}

func serverDTOToConfig(s ServerDTO) config.Server {
	return config.Server{
		Name: s.Name,
		DNS:  s.DNS,
		DoH:  s.DoH,
		DoT:  s.DoT,
		DoQ:  s.DoQ,
	}
}

func formatTimeISO(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func settingsToDTO(s config.Settings) SettingsDTO {
	return SettingsDTO{
		QueryTimeoutMs:   int64(s.QueryTimeout / time.Millisecond),
		MaxConcurrent:    s.MaxConcurrent,
		ServerListURL:    s.ServerListURL,
		LastServerUpdate: formatTimeISO(s.LastServerUpdate),
		EnabledProtocols: append([]string(nil), s.EnabledProtocols...),
		SelectedDomains: TestDomainsDTO{
			A:      append([]string(nil), s.SelectedDomains.A...),
			MX:     append([]string(nil), s.SelectedDomains.MX...),
			TXT:    append([]string(nil), s.SelectedDomains.TXT...),
			DNSSEC: append([]string(nil), s.SelectedDomains.DNSSEC...),
		},
	}
}

func mergeSettingsFromDTO(existing config.Settings, dto SettingsDTO) config.Settings {
	// Preserve values that are managed by the backend (like LastServerUpdate).
	existing.QueryTimeout = time.Duration(dto.QueryTimeoutMs) * time.Millisecond
	existing.MaxConcurrent = dto.MaxConcurrent
	existing.ServerListURL = dto.ServerListURL
	existing.EnabledProtocols = append([]string(nil), dto.EnabledProtocols...)
	existing.SelectedDomains = config.TestDomains{
		A:      append([]string(nil), dto.SelectedDomains.A...),
		MX:     append([]string(nil), dto.SelectedDomains.MX...),
		TXT:    append([]string(nil), dto.SelectedDomains.TXT...),
		DNSSEC: append([]string(nil), dto.SelectedDomains.DNSSEC...),
	}
	return existing
}

func metricsToDTO(m benchmark.Metrics) MetricsDTO {
	return MetricsDTO{
		MinMs:       float64(m.Min) / float64(time.Millisecond),
		MaxMs:       float64(m.Max) / float64(time.Millisecond),
		MeanMs:      float64(m.Mean) / float64(time.Millisecond),
		MedianMs:    float64(m.Median) / float64(time.Millisecond),
		P95Ms:       float64(m.P95) / float64(time.Millisecond),
		P99Ms:       float64(m.P99) / float64(time.Millisecond),
		Success:     m.Success,
		Failed:      m.Failed,
		Total:       m.Total,
		SuccessRate: (&m).SuccessRate(),
	}
}

func benchmarkResultsToDTO(results *benchmark.BenchmarkResults) *BenchmarkResultsDTO {
	if results == nil {
		return nil
	}

	out := &BenchmarkResultsDTO{
		StartTime:  formatTimeISO(results.StartTime),
		EndTime:    formatTimeISO(results.EndTime),
		DurationMs: float64(results.Duration()) / float64(time.Millisecond),
		Results:    make([]ServerResultDTO, 0, len(results.Results)),
	}

	for _, r := range results.Results {
		out.Results = append(out.Results, ServerResultDTO{
			ServerName: r.ServerName,
			Protocol:   r.Protocol.String(),
			Metrics:    metricsToDTO(r.Metrics),
			DurationMs: float64(r.Duration()) / float64(time.Millisecond),
		})
	}

	return out
}
