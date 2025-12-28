package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/x86txt/goDnsBench/internal/benchmark"
	"github.com/x86txt/goDnsBench/internal/config"
	"github.com/x86txt/goDnsBench/internal/dns"
)

// Screen represents different screens in the TUI
type Screen int

const (
	ScreenMain Screen = iota
	ScreenSettings
	ScreenResults
	ScreenRunning
)

// Model represents the TUI application state
type Model struct {
	currentScreen Screen
	servers       []config.Server
	settings      config.Settings
	results       *benchmark.BenchmarkResults
	err           error
	width         int
	height        int
	
	// Progress tracking
	progress      benchmark.Progress
	isRunning     bool
	
	// Benchmark state
	progressChan  chan benchmark.Progress
	doneChan      chan benchmarkDoneMsg
}

// NewModel creates a new TUI model
func NewModel(servers []config.Server, settings config.Settings) Model {
	return Model{
		currentScreen: ScreenMain,
		servers:       servers,
		settings:      settings,
	}
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return nil
}

// benchmarkProgressMsg is sent when benchmark progress updates
type benchmarkProgressMsg benchmark.Progress

// benchmarkDoneMsg is sent when benchmark completes
type benchmarkDoneMsg struct {
	results *benchmark.BenchmarkResults
	err     error
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.isRunning {
				// Allow quit even when running
				return m, tea.Quit
			}
			return m, tea.Quit
		case "s":
			if m.currentScreen == ScreenMain && !m.isRunning {
				m.currentScreen = ScreenSettings
			}
		case "r":
			if m.currentScreen == ScreenMain && !m.isRunning {
				m.currentScreen = ScreenRunning
				m.isRunning = true
				m.progress = benchmark.Progress{}
				// Initialize channels
				m.progressChan = make(chan benchmark.Progress, 10)
				m.doneChan = make(chan benchmarkDoneMsg, 1)
				// Start benchmark and progress polling
				return m, tea.Batch(
					runBenchmark(m.servers, m.settings, m.progressChan, m.doneChan),
					waitForProgress(m.progressChan, m.doneChan),
				)
			}
		case "esc":
			if !m.isRunning {
				m.currentScreen = ScreenMain
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case progressTickMsg:
		// Ticker fired, check for updates
		if m.isRunning && m.progressChan != nil && m.doneChan != nil {
			return m, tea.Batch(
				checkProgress(m.progressChan, m.doneChan),
				waitForProgress(m.progressChan, m.doneChan),
			)
		}
		return m, nil

	case benchmarkProgressMsg:
		m.progress = benchmark.Progress(msg)
		// Continue polling for more updates if still running
		if m.isRunning && m.progressChan != nil && m.doneChan != nil {
			return m, waitForProgress(m.progressChan, m.doneChan)
		}
		return m, nil

	case benchmarkDoneMsg:
		m.isRunning = false
		if msg.err != nil {
			m.err = msg.err
			m.currentScreen = ScreenMain
		} else {
			m.results = msg.results
			m.currentScreen = ScreenResults
		}
		// Clean up channels
		if m.progressChan != nil {
			close(m.progressChan)
			m.progressChan = nil
		}
		if m.doneChan != nil {
			close(m.doneChan)
			m.doneChan = nil
		}
		return m, nil
	}

	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	switch m.currentScreen {
	case ScreenMain:
		return m.viewMain()
	case ScreenSettings:
		return m.viewSettings()
	case ScreenResults:
		return m.viewResults()
	case ScreenRunning:
		return m.viewRunning()
	default:
		return "Unknown screen"
	}
}

// viewMain renders the main menu
func (m Model) viewMain() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0066ff")).
		MarginBottom(1)

	menuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ffff"))

	title := titleStyle.Render("goDnsBench - DNS Benchmarking Tool")

	menu := menuStyle.Render(fmt.Sprintf(`
Loaded %d DNS servers

Commands:
  [r] Run benchmark
  [s] Settings
  [q] Quit

Press a key to continue...
`, len(m.servers)))

	return title + "\n" + menu
}

// viewSettings renders the settings screen
func (m Model) viewSettings() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0066ff")).
		MarginBottom(1)

	title := titleStyle.Render("Settings")

	content := fmt.Sprintf(`
Query Timeout: %v
Max Concurrent: %d
Enabled Protocols: %v

[esc] Back to main menu
`, m.settings.QueryTimeout, m.settings.MaxConcurrent, m.settings.EnabledProtocols)

	return title + "\n" + content
}

// viewResults renders the results screen
func (m Model) viewResults() string {
	if m.results == nil {
		return "No results available"
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0066ff")).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ffff")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	title := titleStyle.Render("Benchmark Results")

	// Summary
	summary := fmt.Sprintf(`
Duration: %v
Total Results: %d

`, m.results.Duration(), len(m.results.Results))

	// Table header
	header := headerStyle.Render(fmt.Sprintf("%-20s %-6s %-8s %-8s %-8s %-8s %-8s",
		"Server", "Proto", "Min(ms)", "Mean(ms)", "P95(ms)", "P99(ms)", "Success%"))

	// Table rows
	var rows []string
	for _, result := range m.results.Results {
		serverName := result.ServerName
		if len(serverName) > 20 {
			serverName = serverName[:17] + "..."
		}
		
		protocol := result.Protocol.String()
		minMs := result.Metrics.Min.Seconds() * 1000
		meanMs := result.Metrics.Mean.Seconds() * 1000
		p95Ms := result.Metrics.P95.Seconds() * 1000
		p99Ms := result.Metrics.P99.Seconds() * 1000
		successRate := result.Metrics.SuccessRate()

		row := cellStyle.Render(fmt.Sprintf("%-20s %-6s %8.2f %8.2f %8.2f %8.2f %7.1f%%",
			serverName, protocol, minMs, meanMs, p95Ms, p99Ms, successRate))
		rows = append(rows, row)
	}

	table := header + "\n" + strings.Repeat("-", 80) + "\n" + strings.Join(rows, "\n")

	content := summary + table + "\n\n[esc] Back to main menu"

	return title + "\n" + content
}

// viewRunning renders the running benchmark screen
func (m Model) viewRunning() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0066ff")).
		MarginBottom(1)

	title := titleStyle.Render("Running Benchmark...")

	// Progress bar
	progressBarWidth := 50
	if m.width > 0 && m.width < 80 {
		progressBarWidth = m.width - 10
	}
	
	progressPercent := int(m.progress.Percentage)
	if progressPercent > 100 {
		progressPercent = 100
	}
	if progressPercent < 0 {
		progressPercent = 0
	}
	
	filled := progressPercent * progressBarWidth / 100
	empty := progressBarWidth - filled
	
	progressBar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ffff"))
	
	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))
	
	content := fmt.Sprintf(`
%s

Current: %s (%s)
Progress: %d/%d tests (%.1f%%)

%s

[ctrl+c] Quit
`, 
		infoStyle.Render("Benchmarking in progress..."),
		m.progress.CurrentServer,
		m.progress.CurrentProtocol.String(),
		m.progress.CompletedTests,
		m.progress.TotalTests,
		m.progress.Percentage,
		progressStyle.Render(progressBar),
	)

	return title + "\n" + content
}

// runBenchmark starts a benchmark run asynchronously
func runBenchmark(servers []config.Server, settings config.Settings, progressChan chan<- benchmark.Progress, doneChan chan<- benchmarkDoneMsg) tea.Cmd {
	return func() tea.Msg {
		// Convert protocol strings to Protocol types
		var protoList []dns.Protocol
		for _, p := range settings.EnabledProtocols {
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
			doneChan <- benchmarkDoneMsg{
				err: fmt.Errorf("no protocols enabled"),
			}
			return nil
		}

		// Create benchmark config with progress callback
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
				select {
				case progressChan <- p:
				default:
					// Channel full, skip update
				}
			},
		}

		// Start benchmark in a goroutine
		go func() {
			// Create runner
			runner := benchmark.NewRunner(benchConfig)

			// Create context for cancellation
			ctx := context.Background()

			// Run benchmark
			results, err := runner.Run(ctx)
			
			// Send done message
			doneChan <- benchmarkDoneMsg{
				results: results,
				err:     err,
			}
		}()

		// Return initial progress message
		return benchmarkProgressMsg{
			CurrentServer:   "Starting...",
			CurrentProtocol: dns.ProtocolDNS,
			CompletedTests:  0,
			TotalTests:      len(servers) * len(protoList),
			Percentage:      0,
		}
	}
}

// progressTickMsg is a message sent by the ticker to check for updates
type progressTickMsg struct{}

// waitForProgress waits for progress updates and sends them as messages
// Uses a ticker to poll channels periodically
func waitForProgress(progressChan <-chan benchmark.Progress, doneChan <-chan benchmarkDoneMsg) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return progressTickMsg{}
	})
}

// checkProgress checks channels for updates
func checkProgress(progressChan <-chan benchmark.Progress, doneChan <-chan benchmarkDoneMsg) tea.Cmd {
	return func() tea.Msg {
		select {
		case progress, ok := <-progressChan:
			if ok {
				return benchmarkProgressMsg(progress)
			}
			// Channel closed, check for done
		case done := <-doneChan:
			return done
		default:
			// No update yet
		}
		return nil
	}
}

// Run starts the TUI application
func Run(servers []config.Server, settings config.Settings) error {
	p := tea.NewProgram(NewModel(servers, settings))
	_, err := p.Run()
	return err
}
