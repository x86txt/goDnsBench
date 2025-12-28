package benchmark

import (
	"testing"
	"time"

	"github.com/x86txt/goDnsBench/internal/dns"
)

func TestServerResult_Duration(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Second)
	result := ServerResult{
		StartTime: start,
		EndTime:   end,
	}

	duration := result.Duration()
	expected := 5 * time.Second

	if duration != expected {
		t.Errorf("Expected Duration=%v, got %v", expected, duration)
	}
}

func TestBenchmarkResults_Duration(t *testing.T) {
	start := time.Now()
	end := start.Add(10 * time.Second)
	results := BenchmarkResults{
		StartTime: start,
		EndTime:   end,
	}

	duration := results.Duration()
	expected := 10 * time.Second

	if duration != expected {
		t.Errorf("Expected Duration=%v, got %v", expected, duration)
	}
}

func TestBenchmarkResults_GetResultsByProtocol(t *testing.T) {
	results := BenchmarkResults{
		Results: []ServerResult{
			{ServerName: "Server1", Protocol: dns.ProtocolDNS},
			{ServerName: "Server2", Protocol: dns.ProtocolDoH},
			{ServerName: "Server3", Protocol: dns.ProtocolDNS},
			{ServerName: "Server4", Protocol: dns.ProtocolDoT},
		},
	}

	dnsResults := results.GetResultsByProtocol(dns.ProtocolDNS)
	if len(dnsResults) != 2 {
		t.Errorf("Expected 2 DNS results, got %d", len(dnsResults))
	}
	for _, r := range dnsResults {
		if r.Protocol != dns.ProtocolDNS {
			t.Errorf("Expected Protocol=DNS, got %v", r.Protocol)
		}
	}

	dohResults := results.GetResultsByProtocol(dns.ProtocolDoH)
	if len(dohResults) != 1 {
		t.Errorf("Expected 1 DoH result, got %d", len(dohResults))
	}

	emptyResults := results.GetResultsByProtocol(dns.ProtocolDoQ)
	if len(emptyResults) != 0 {
		t.Errorf("Expected 0 DoQ results, got %d", len(emptyResults))
	}
}

func TestBenchmarkResults_GetResultsByServer(t *testing.T) {
	results := BenchmarkResults{
		Results: []ServerResult{
			{ServerName: "Server1", Protocol: dns.ProtocolDNS},
			{ServerName: "Server1", Protocol: dns.ProtocolDoH},
			{ServerName: "Server2", Protocol: dns.ProtocolDNS},
			{ServerName: "Server3", Protocol: dns.ProtocolDoT},
		},
	}

	server1Results := results.GetResultsByServer("Server1")
	if len(server1Results) != 2 {
		t.Errorf("Expected 2 results for Server1, got %d", len(server1Results))
	}
	for _, r := range server1Results {
		if r.ServerName != "Server1" {
			t.Errorf("Expected ServerName=Server1, got %s", r.ServerName)
		}
	}

	server2Results := results.GetResultsByServer("Server2")
	if len(server2Results) != 1 {
		t.Errorf("Expected 1 result for Server2, got %d", len(server2Results))
	}

	emptyResults := results.GetResultsByServer("NonExistent")
	if len(emptyResults) != 0 {
		t.Errorf("Expected 0 results for NonExistent, got %d", len(emptyResults))
	}
}

func TestBenchmarkResults_GetFastestServer(t *testing.T) {
	tests := []struct {
		name     string
		results  BenchmarkResults
		expected string
	}{
		{
			name: "Single result",
			results: BenchmarkResults{
				Results: []ServerResult{
					{ServerName: "Server1", Metrics: Metrics{Mean: 100 * time.Millisecond}},
				},
			},
			expected: "Server1",
		},
		{
			name: "Multiple results",
			results: BenchmarkResults{
				Results: []ServerResult{
					{ServerName: "Server1", Metrics: Metrics{Mean: 150 * time.Millisecond}},
					{ServerName: "Server2", Metrics: Metrics{Mean: 50 * time.Millisecond}},
					{ServerName: "Server3", Metrics: Metrics{Mean: 200 * time.Millisecond}},
				},
			},
			expected: "Server2",
		},
		{
			name: "Tie (first wins)",
			results: BenchmarkResults{
				Results: []ServerResult{
					{ServerName: "Server1", Metrics: Metrics{Mean: 100 * time.Millisecond}},
					{ServerName: "Server2", Metrics: Metrics{Mean: 100 * time.Millisecond}},
				},
			},
			expected: "Server1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fastest := tt.results.GetFastestServer()
			if fastest == nil {
				t.Fatal("Expected non-nil result")
			}
			if fastest.ServerName != tt.expected {
				t.Errorf("Expected fastest server=%s, got %s", tt.expected, fastest.ServerName)
			}
		})
	}

	// Test empty results
	emptyResults := BenchmarkResults{Results: []ServerResult{}}
	if emptyResults.GetFastestServer() != nil {
		t.Error("Expected nil for empty results")
	}
}

func TestBenchmarkResults_GetSlowestServer(t *testing.T) {
	results := BenchmarkResults{
		Results: []ServerResult{
			{ServerName: "Server1", Metrics: Metrics{Mean: 50 * time.Millisecond}},
			{ServerName: "Server2", Metrics: Metrics{Mean: 200 * time.Millisecond}},
			{ServerName: "Server3", Metrics: Metrics{Mean: 100 * time.Millisecond}},
		},
	}

	slowest := results.GetSlowestServer()
	if slowest == nil {
		t.Fatal("Expected non-nil result")
	}
	if slowest.ServerName != "Server2" {
		t.Errorf("Expected slowest server=Server2, got %s", slowest.ServerName)
	}

	// Test empty results
	emptyResults := BenchmarkResults{Results: []ServerResult{}}
	if emptyResults.GetSlowestServer() != nil {
		t.Error("Expected nil for empty results")
	}
}

func TestBenchmarkResults_GetMostReliable(t *testing.T) {
	results := BenchmarkResults{
		Results: []ServerResult{
			{ServerName: "Server1", Metrics: Metrics{Success: 8, Total: 10}}, // 80%
			{ServerName: "Server2", Metrics: Metrics{Success: 10, Total: 10}}, // 100%
			{ServerName: "Server3", Metrics: Metrics{Success: 5, Total: 10}}, // 50%
		},
	}

	mostReliable := results.GetMostReliable()
	if mostReliable == nil {
		t.Fatal("Expected non-nil result")
	}
	if mostReliable.ServerName != "Server2" {
		t.Errorf("Expected most reliable server=Server2, got %s", mostReliable.ServerName)
	}

	// Test empty results
	emptyResults := BenchmarkResults{Results: []ServerResult{}}
	if emptyResults.GetMostReliable() != nil {
		t.Error("Expected nil for empty results")
	}
}
