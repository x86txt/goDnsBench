package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/x86txt/goDnsBench/internal/benchmark"
	"github.com/x86txt/goDnsBench/internal/config"
	"github.com/x86txt/goDnsBench/internal/dns"
)

// RunHeadlessBenchmark runs a benchmark in headless mode and returns the results
func RunHeadlessBenchmark(servers []config.Server, settings config.Settings, protocols []string) (*benchmark.BenchmarkResults, error) {
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
		default:
			fmt.Fprintf(os.Stderr, "Warning: unknown protocol '%s', skipping\n", p)
		}
	}

	if len(protoList) == 0 {
		return nil, fmt.Errorf("no valid protocols specified")
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers specified")
	}

	// Create benchmark config
	benchConfig := benchmark.BenchmarkConfig{
		Servers:       servers,
		Protocols:     protoList,
		Timeout:       settings.QueryTimeout,
		MaxConcurrent: settings.MaxConcurrent,
		TestDomains: benchmark.TestDomains{
			A:      settings.SelectedDomains.A,
			MX:     settings.SelectedDomains.MX,
			TXT:    settings.SelectedDomains.TXT,
			DNSSEC: settings.SelectedDomains.DNSSEC,
		},
		ProgressCallback: func(p benchmark.Progress) {
			fmt.Fprintf(os.Stderr, "Progress: %s (%s) - %.1f%% (%d/%d)\n",
				p.CurrentServer,
				p.CurrentProtocol.String(),
				p.Percentage,
				p.CompletedTests,
				p.TotalTests,
			)
		},
	}

	// Create runner
	runner := benchmark.NewRunner(benchConfig)

	// Run benchmark
	fmt.Fprintf(os.Stderr, "Running benchmark with %d server(s) and %d protocol(s)...\n", len(servers), len(protoList))
	ctx := context.Background()
	results, err := runner.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("benchmark failed: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Benchmark completed in %s\n", results.Duration())
	return results, nil
}
