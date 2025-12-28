package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/x86txt/goDnsBench/internal/cli"
	"github.com/x86txt/goDnsBench/internal/config"
	"github.com/x86txt/goDnsBench/internal/export"
	"github.com/x86txt/goDnsBench/internal/tui"
)

var (
	version = "0.1.0"
)

func main() {
	// Parse command line flags
	tuiMode := flag.Bool("tui", false, "Run in terminal UI mode")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	serverFile := flag.String("servers", "", "Path to servers JSON/CSV file")
	exportJSON := flag.String("export-json", "", "Export results to JSON file (runs headless benchmark)")
	exportCSV := flag.String("export-csv", "", "Export results to CSV file (runs headless benchmark)")
	protocolsFlag := flag.String("protocols", "DNS,DoH,DoT,DoQ", "Comma-separated list of protocols to test (DNS,DoH,DoT,DoQ)")
	timeout := flag.Int("timeout", 1000, "Query timeout in milliseconds")
	concurrent := flag.Int("concurrent", 10, "Maximum concurrent servers to test")

	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("goDnsBench version %s\n", version)
		os.Exit(0)
	}

	// Load or initialize settings
	settings, err := config.LoadSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load settings: %v\n", err)
		settings = config.DefaultSettings()
	}

	// Override settings with command line flags if provided
	if *timeout > 0 {
		settings.QueryTimeout = time.Duration(*timeout) * time.Millisecond
	}
	if *concurrent > 0 {
		settings.MaxConcurrent = *concurrent
	}

	// Load servers
	var servers []config.Server
	if *serverFile != "" {
		// Load from specified file
		servers, err = loadServersFromFile(*serverFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading servers: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Use default servers
		servers = config.DefaultServers()
	}

	// Handle export flags (headless mode)
	if *exportJSON != "" || *exportCSV != "" {
		// Parse protocols
		protocols := parseProtocols(*protocolsFlag)

		// Run headless benchmark
		results, err := cli.RunHeadlessBenchmark(servers, settings, protocols)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running benchmark: %v\n", err)
			os.Exit(1)
		}

		// Export results
		if *exportJSON != "" {
			if err := export.ExportResultsJSON(results, *exportJSON); err != nil {
				fmt.Fprintf(os.Stderr, "Error exporting JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Results exported to %s\n", *exportJSON)
		}

		if *exportCSV != "" {
			if err := export.ExportResultsCSV(results, *exportCSV); err != nil {
				fmt.Fprintf(os.Stderr, "Error exporting CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Results exported to %s\n", *exportCSV)
		}

		return
	}

	// Launch appropriate interface
	if *tuiMode {
		// Launch TUI
		if err := tui.Run(servers, settings); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Launch GUI with Wails
		if err := LaunchGUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running GUI: %v\n", err)
			os.Exit(1)
		}
	}
}

// parseProtocols parses a comma-separated list of protocol names
func parseProtocols(protocolsStr string) []string {
	if protocolsStr == "" {
		return []string{"DNS", "DoH", "DoT", "DoQ"}
	}

	parts := strings.Split(protocolsStr, ",")
	var protocols []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			protocols = append(protocols, p)
		}
	}

	if len(protocols) == 0 {
		return []string{"DNS", "DoH", "DoT", "DoQ"}
	}

	return protocols
}

// loadServersFromFile loads servers from a file (JSON or CSV based on extension)
func loadServersFromFile(filepath string) ([]config.Server, error) {
	// Check file extension
	if len(filepath) >= 5 && filepath[len(filepath)-5:] == ".json" {
		return config.LoadServersFromJSON(filepath)
	} else if len(filepath) >= 4 && filepath[len(filepath)-4:] == ".csv" {
		return config.LoadServersFromCSV(filepath)
	}

	return nil, fmt.Errorf("unsupported file format (must be .json or .csv)")
}
