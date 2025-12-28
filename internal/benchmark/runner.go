package benchmark

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/x86txt/goDnsBench/internal/config"
	"github.com/x86txt/goDnsBench/internal/dns"
)

// BenchmarkConfig holds the configuration for a benchmark run
type BenchmarkConfig struct {
	Servers          []config.Server
	Protocols        []dns.Protocol
	Timeout          time.Duration
	MaxConcurrent    int
	TestDomains      TestDomains
	ProgressCallback func(progress Progress)
}

// TestDomains defines the domains to test for each query type
type TestDomains struct {
	A      []string
	MX     []string
	TXT    []string
	DNSSEC []string
}

// DefaultTestDomains returns the default set of test domains
func DefaultTestDomains() TestDomains {
	return TestDomains{
		A:      []string{"google.com", "cloudflare.com", "amazon.com", "microsoft.com"},
		MX:     []string{"gmail.com", "microsoft.com"},
		TXT:    []string{"_dmarc.google.com", "google.com"},
		DNSSEC: []string{"cloudflare.com", "google.com"},
	}
}

// Progress represents the current progress of a benchmark run
type Progress struct {
	CurrentServer   string
	CurrentProtocol dns.Protocol
	CompletedTests  int
	TotalTests      int
	Percentage      float64
}

// Runner orchestrates benchmark execution
type Runner struct {
	config BenchmarkConfig
}

// NewRunner creates a new benchmark runner
func NewRunner(config BenchmarkConfig) *Runner {
	if config.Timeout == 0 {
		config.Timeout = 1 * time.Second
	}
	if config.MaxConcurrent == 0 {
		config.MaxConcurrent = 10
	}
	if len(config.TestDomains.A) == 0 {
		config.TestDomains = DefaultTestDomains()
	}
	return &Runner{
		config: config,
	}
}

// Run executes the benchmark
func (r *Runner) Run(ctx context.Context) (*BenchmarkResults, error) {
	results := &BenchmarkResults{
		StartTime: time.Now(),
		Config:    r.config,
	}

	// Calculate total tests
	totalTests := len(r.config.Servers) * len(r.config.Protocols)

	// Create a semaphore for concurrency control
	sem := make(chan struct{}, r.config.MaxConcurrent)

	// WaitGroup to wait for all goroutines
	var wg sync.WaitGroup

	// Mutex to protect results
	var mu sync.Mutex

	var completedTests atomic.Int64

	// Iterate through servers and protocols
	for _, server := range r.config.Servers {
		for _, protocol := range r.config.Protocols {
			// Check if server supports this protocol
			if !r.serverSupportsProtocol(server, protocol) {
				newCompleted := completedTests.Add(1)
				if r.config.ProgressCallback != nil {
					r.config.ProgressCallback(Progress{
						CurrentServer:   server.Name,
						CurrentProtocol: protocol,
						CompletedTests:  int(newCompleted),
						TotalTests:      totalTests,
						Percentage:      float64(newCompleted) / float64(totalTests) * 100,
					})
				}
				continue
			}

			wg.Add(1)

			// Acquire semaphore
			sem <- struct{}{}

			go func(srv config.Server, proto dns.Protocol) {
				defer wg.Done()
				defer func() { <-sem }()

				// Check context
				select {
				case <-ctx.Done():
					return
				default:
				}

				// Run benchmark for this server and protocol
				result := r.runServerBenchmark(srv, proto)

				// Store result
				mu.Lock()
				results.Results = append(results.Results, result)
				mu.Unlock()

				newCompleted := completedTests.Add(1)

				// Report progress
				if r.config.ProgressCallback != nil {
					r.config.ProgressCallback(Progress{
						CurrentServer:   srv.Name,
						CurrentProtocol: proto,
						CompletedTests:  int(newCompleted),
						TotalTests:      totalTests,
						Percentage:      float64(newCompleted) / float64(totalTests) * 100,
					})
				}
			}(server, protocol)
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()

	results.EndTime = time.Now()
	return results, nil
}

// runServerBenchmark runs a benchmark for a single server and protocol
func (r *Runner) runServerBenchmark(server config.Server, protocol dns.Protocol) ServerResult {
	result := ServerResult{
		ServerName: server.Name,
		Protocol:   protocol,
		StartTime:  time.Now(),
	}

	// Create client for this protocol
	client, err := r.createClient(server, protocol)
	if err != nil {
		result.EndTime = time.Now()
		return result
	}
	defer client.Close()

	// Build query list
	queries := r.buildQueries(protocol)

	// Execute queries sequentially (per requirements)
	var latencies []time.Duration
	successCount := 0
	failedCount := 0

	for _, query := range queries {
		queryResult, err := client.Query(query)
		if err != nil {
			failedCount++
			continue
		}

		result.Queries = append(result.Queries, *queryResult)

		if queryResult.Success {
			successCount++
			latencies = append(latencies, queryResult.Latency)
		} else {
			failedCount++
		}
	}

	// Calculate metrics
	result.Metrics = CalculateMetrics(latencies, successCount, failedCount)
	result.EndTime = time.Now()

	return result
}

// createClient creates a DNS client for the specified protocol
func (r *Runner) createClient(server config.Server, protocol dns.Protocol) (dns.Client, error) {
	switch protocol {
	case dns.ProtocolDNS:
		if server.DNS == "" {
			return nil, fmt.Errorf("server does not support DNS")
		}
		return dns.NewStandardClient(server.DNS, r.config.Timeout), nil

	case dns.ProtocolDoH:
		if server.DoH == "" {
			return nil, fmt.Errorf("server does not support DoH")
		}
		return dns.NewDoHClient(server.DoH, r.config.Timeout), nil

	case dns.ProtocolDoT:
		if server.DoT == "" {
			return nil, fmt.Errorf("server does not support DoT")
		}
		return dns.NewDoTClient(server.DoT, r.config.Timeout), nil

	case dns.ProtocolDoQ:
		if server.DoQ == "" {
			return nil, fmt.Errorf("server does not support DoQ")
		}
		return dns.NewDoQClient(server.DoQ, r.config.Timeout), nil

	default:
		return nil, fmt.Errorf("unsupported protocol: %v", protocol)
	}
}

// buildQueries creates the list of queries based on test domains
func (r *Runner) buildQueries(protocol dns.Protocol) []dns.Query {
	var queries []dns.Query

	// A records (4 queries)
	for _, domain := range r.config.TestDomains.A {
		queries = append(queries, dns.Query{
			Domain:   domain,
			Type:     dns.QueryTypeA,
			Protocol: protocol,
			Timeout:  r.config.Timeout,
		})
	}

	// MX records (2 queries)
	for _, domain := range r.config.TestDomains.MX {
		queries = append(queries, dns.Query{
			Domain:   domain,
			Type:     dns.QueryTypeMX,
			Protocol: protocol,
			Timeout:  r.config.Timeout,
		})
	}

	// TXT records (2 queries)
	for _, domain := range r.config.TestDomains.TXT {
		queries = append(queries, dns.Query{
			Domain:   domain,
			Type:     dns.QueryTypeTXT,
			Protocol: protocol,
			Timeout:  r.config.Timeout,
		})
	}

	// DNSSEC queries (2 queries)
	for _, domain := range r.config.TestDomains.DNSSEC {
		queries = append(queries, dns.Query{
			Domain:   domain,
			Type:     dns.QueryTypeDNSSEC,
			Protocol: protocol,
			Timeout:  r.config.Timeout,
		})
	}

	return queries
}

// serverSupportsProtocol checks if a server supports a given protocol
func (r *Runner) serverSupportsProtocol(server config.Server, protocol dns.Protocol) bool {
	switch protocol {
	case dns.ProtocolDNS:
		return server.DNS != ""
	case dns.ProtocolDoH:
		return server.DoH != ""
	case dns.ProtocolDoT:
		return server.DoT != ""
	case dns.ProtocolDoQ:
		return server.DoQ != ""
	default:
		return false
	}
}
