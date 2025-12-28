package benchmark

import (
	"time"

	"github.com/x86txt/goDnsBench/internal/dns"
)

// ServerResult represents the benchmark results for a single server
type ServerResult struct {
	ServerName string
	Protocol   dns.Protocol
	Metrics    Metrics
	Queries    []dns.QueryResult
	StartTime  time.Time
	EndTime    time.Time
}

// Duration returns the total time taken for this server's benchmark
func (sr *ServerResult) Duration() time.Duration {
	return sr.EndTime.Sub(sr.StartTime)
}

// BenchmarkResults represents the complete results of a benchmark run
type BenchmarkResults struct {
	Results   []ServerResult
	StartTime time.Time
	EndTime   time.Time
	Config    BenchmarkConfig
}

// Duration returns the total time taken for the entire benchmark
func (br *BenchmarkResults) Duration() time.Duration {
	return br.EndTime.Sub(br.StartTime)
}

// GetResultsByProtocol filters results by protocol
func (br *BenchmarkResults) GetResultsByProtocol(protocol dns.Protocol) []ServerResult {
	var filtered []ServerResult
	for _, result := range br.Results {
		if result.Protocol == protocol {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// GetResultsByServer filters results by server name
func (br *BenchmarkResults) GetResultsByServer(serverName string) []ServerResult {
	var filtered []ServerResult
	for _, result := range br.Results {
		if result.ServerName == serverName {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// GetFastestServer returns the server with the lowest mean latency
func (br *BenchmarkResults) GetFastestServer() *ServerResult {
	if len(br.Results) == 0 {
		return nil
	}

	fastest := &br.Results[0]
	for i := range br.Results {
		if br.Results[i].Metrics.Mean < fastest.Metrics.Mean {
			fastest = &br.Results[i]
		}
	}
	return fastest
}

// GetSlowestServer returns the server with the highest mean latency
func (br *BenchmarkResults) GetSlowestServer() *ServerResult {
	if len(br.Results) == 0 {
		return nil
	}

	slowest := &br.Results[0]
	for i := range br.Results {
		if br.Results[i].Metrics.Mean > slowest.Metrics.Mean {
			slowest = &br.Results[i]
		}
	}
	return slowest
}

// GetMostReliable returns the server with the highest success rate
func (br *BenchmarkResults) GetMostReliable() *ServerResult {
	if len(br.Results) == 0 {
		return nil
	}

	mostReliable := &br.Results[0]
	for i := range br.Results {
		if br.Results[i].Metrics.SuccessRate() > mostReliable.Metrics.SuccessRate() {
			mostReliable = &br.Results[i]
		}
	}
	return mostReliable
}
