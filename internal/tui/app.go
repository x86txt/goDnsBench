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

// Theme colors - Deep space aesthetic
var (
	colorPrimary   = lipgloss.Color("#58a6ff")
	colorSecondary = lipgloss.Color("#a371f7")
	colorSuccess   = lipgloss.Color("#3fb950")
	colorWarning   = lipgloss.Color("#d29922")
	colorError     = lipgloss.Color("#f85149")
	colorMuted     = lipgloss.Color("#8b949e")
	colorDim       = lipgloss.Color("#484f58")
	colorBg        = lipgloss.Color("#0d1117")
	colorBgAlt     = lipgloss.Color("#161b22")
	colorBorder    = lipgloss.Color("#30363d")
	colorText      = lipgloss.Color("#c9d1d9")
	colorTextBright = lipgloss.Color("#f0f6fc")
)

// Styles
var (
	// Base styles
	styleTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextBright).
		Background(colorPrimary).
		Padding(0, 2)

	styleSubtitle = lipgloss.NewStyle().
		Foreground(colorMuted).
		Italic(true)

	styleBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 2)

	styleHighlight = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true)

	styleSuccess = lipgloss.NewStyle().
		Foreground(colorSuccess)

	styleWarning = lipgloss.NewStyle().
		Foreground(colorWarning)

	styleError = lipgloss.NewStyle().
		Foreground(colorError)

	styleMuted = lipgloss.NewStyle().
		Foreground(colorMuted)

	styleDim = lipgloss.NewStyle().
		Foreground(colorDim)

	styleKey = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color("#21262d"))

	styleValue = lipgloss.NewStyle().
		Foreground(colorText)

	styleTableHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextBright).
		Background(lipgloss.Color("#21262d")).
		Padding(0, 1)

	styleTableCell = lipgloss.NewStyle().
		Foreground(colorText).
		Padding(0, 1)

	styleProgressBar = lipgloss.NewStyle().
		Foreground(colorPrimary)

	styleProgressTrack = lipgloss.NewStyle().
		Foreground(colorDim)
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
	progress  benchmark.Progress
	isRunning bool
	startTime time.Time

	// Benchmark state
	progressChan chan benchmark.Progress
	doneChan     chan benchmarkDoneMsg
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
			return m, tea.Quit
		case "s":
			if m.currentScreen == ScreenMain && !m.isRunning {
				m.currentScreen = ScreenSettings
			}
		case "r":
			if m.currentScreen == ScreenMain && !m.isRunning {
				m.currentScreen = ScreenRunning
				m.isRunning = true
				m.startTime = time.Now()
				m.progress = benchmark.Progress{}
				m.progressChan = make(chan benchmark.Progress, 10)
				m.doneChan = make(chan benchmarkDoneMsg, 1)
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
		if m.isRunning && m.progressChan != nil && m.doneChan != nil {
			return m, tea.Batch(
				checkProgress(m.progressChan, m.doneChan),
				waitForProgress(m.progressChan, m.doneChan),
			)
		}
		return m, nil

	case benchmarkProgressMsg:
		m.progress = benchmark.Progress(msg)
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
	width := m.width
	if width < 60 {
		width = 60
	}
	if width > 100 {
		width = 100
	}

	// Logo/Header
	logo := `
   ██████╗  ██████╗ ██████╗ ███╗   ██╗███████╗██████╗ ███████╗███╗   ██╗ ██████╗██╗  ██╗
  ██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║██╔════╝██╔══██╗██╔════╝████╗  ██║██╔════╝██║  ██║
  ██║  ███╗██║   ██║██║  ██║██╔██╗ ██║███████╗██████╔╝█████╗  ██╔██╗ ██║██║     ███████║
  ██║   ██║██║   ██║██║  ██║██║╚██╗██║╚════██║██╔══██╗██╔══╝  ██║╚██╗██║██║     ██╔══██║
  ╚██████╔╝╚██████╔╝██████╔╝██║ ╚████║███████║██████╔╝███████╗██║ ╚████║╚██████╗██║  ██║
   ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝╚═════╝ ╚══════╝╚═╝  ╚═══╝ ╚═════╝╚═╝  ╚═╝`

	logoStyled := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Render(logo)

	subtitle := styleMuted.Render("                           DNS Performance Analyzer • v0.1.0")

	// Divider
	divider := styleDim.Render(strings.Repeat("─", width))

	// Stats box
	statsBox := m.renderStatsBox(width)

	// Menu
	menuItems := []struct {
		key  string
		desc string
	}{
		{"r", "Run Benchmark"},
		{"s", "Settings"},
		{"q", "Quit"},
	}

	var menuLines []string
	menuLines = append(menuLines, "")
	menuLines = append(menuLines, styleHighlight.Render("  ⚡ Commands"))
	menuLines = append(menuLines, "")
	for _, item := range menuItems {
		line := fmt.Sprintf("     %s  %s",
			styleKey.Render(item.key),
			styleValue.Render(item.desc))
		menuLines = append(menuLines, line)
	}
	menuLines = append(menuLines, "")

	menu := strings.Join(menuLines, "\n")

	// Error display
	errorMsg := ""
	if m.err != nil {
		errorMsg = "\n" + styleError.Render("  ⚠ Error: "+m.err.Error()) + "\n"
	}

	// Combine all
	content := strings.Join([]string{
		"",
		logoStyled,
		subtitle,
		"",
		divider,
		statsBox,
		divider,
		menu,
		errorMsg,
	}, "\n")

	return content
}

// renderStatsBox renders a statistics box
func (m Model) renderStatsBox(width int) string {
	serverCount := len(m.servers)
	protocolCount := len(m.settings.EnabledProtocols)
	queriesPerServer := 10 // Fixed in the app

	// Create stat items
	stat1 := fmt.Sprintf("  %s %s",
		styleHighlight.Render(fmt.Sprintf("%d", serverCount)),
		styleMuted.Render("Servers"))

	stat2 := fmt.Sprintf("  %s %s",
		styleSuccess.Render(fmt.Sprintf("%d", protocolCount)),
		styleMuted.Render("Protocols"))

	stat3 := fmt.Sprintf("  %s %s",
		lipgloss.NewStyle().Foreground(colorSecondary).Render(fmt.Sprintf("%d", queriesPerServer)),
		styleMuted.Render("Queries/Server"))

	totalTests := serverCount * protocolCount * queriesPerServer
	stat4 := fmt.Sprintf("  %s %s",
		styleWarning.Render(fmt.Sprintf("%d", totalTests)),
		styleMuted.Render("Total Queries"))

	return fmt.Sprintf("\n%s    │    %s    │    %s    │    %s\n", stat1, stat2, stat3, stat4)
}

// viewSettings renders the settings screen
func (m Model) viewSettings() string {
	title := styleTitle.Render(" ⚙ Settings ")

	var lines []string
	lines = append(lines, "")
	lines = append(lines, title)
	lines = append(lines, "")

	// Settings display
	settings := []struct {
		label string
		value string
	}{
		{"Query Timeout", fmt.Sprintf("%v", m.settings.QueryTimeout)},
		{"Max Concurrent", fmt.Sprintf("%d servers", m.settings.MaxConcurrent)},
		{"Protocols", strings.Join(m.settings.EnabledProtocols, ", ")},
	}

	for _, s := range settings {
		line := fmt.Sprintf("  %s: %s",
			styleMuted.Render(s.label),
			styleValue.Render(s.value))
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, styleDim.Render("  ───────────────────────────────────"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s Back to main menu", styleKey.Render("esc")))
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

// viewResults renders the results screen
func (m Model) viewResults() string {
	if m.results == nil {
		return styleMuted.Render("No results available")
	}

	title := styleTitle.Render(" 📊 Benchmark Results ")

	var lines []string
	lines = append(lines, "")
	lines = append(lines, title)
	lines = append(lines, "")

	// Summary stats
	duration := m.results.Duration()
	resultCount := len(m.results.Results)

	summaryLine := fmt.Sprintf("  Completed %s tests in %s",
		styleHighlight.Render(fmt.Sprintf("%d", resultCount)),
		styleSuccess.Render(duration.Round(time.Millisecond).String()))
	lines = append(lines, summaryLine)
	lines = append(lines, "")

	// Table header
	headerFormat := "  %-22s %-6s %10s %10s %10s %10s"
	header := fmt.Sprintf(headerFormat, "SERVER", "PROTO", "MIN", "MEAN", "P95", "SUCCESS")
	lines = append(lines, styleTableHeader.Render(header))
	lines = append(lines, styleDim.Render("  "+strings.Repeat("─", 74)))

	// Table rows
	for _, result := range m.results.Results {
		serverName := result.ServerName
		if len(serverName) > 22 {
			serverName = serverName[:19] + "..."
		}

		protocol := result.Protocol.String()
		minMs := result.Metrics.Min.Seconds() * 1000
		meanMs := result.Metrics.Mean.Seconds() * 1000
		p95Ms := result.Metrics.P95.Seconds() * 1000
		successRate := result.Metrics.SuccessRate()

		// Color code the latency
		meanStyle := styleValue
		if meanMs < 50 {
			meanStyle = styleSuccess
		} else if meanMs < 100 {
			meanStyle = styleWarning
		} else {
			meanStyle = styleError
		}

		// Color code success rate
		successStyle := styleSuccess
		if successRate < 100 {
			successStyle = styleWarning
		}
		if successRate < 50 {
			successStyle = styleError
		}

		// Protocol badge color
		protoStyle := styleMuted
		switch result.Protocol {
		case dns.ProtocolDNS:
			protoStyle = styleSuccess
		case dns.ProtocolDoH:
			protoStyle = lipgloss.NewStyle().Foreground(colorPrimary)
		case dns.ProtocolDoT:
			protoStyle = lipgloss.NewStyle().Foreground(colorSecondary)
		case dns.ProtocolDoQ:
			protoStyle = styleWarning
		}

		row := fmt.Sprintf("  %-22s %s %10s %s %10s %s",
			styleValue.Render(serverName),
			protoStyle.Render(fmt.Sprintf("%-6s", protocol)),
			styleMuted.Render(fmt.Sprintf("%.1fms", minMs)),
			meanStyle.Render(fmt.Sprintf("%10s", fmt.Sprintf("%.1fms", meanMs))),
			styleMuted.Render(fmt.Sprintf("%.1fms", p95Ms)),
			successStyle.Render(fmt.Sprintf("%9.0f%%", successRate)))
		lines = append(lines, row)
	}

	lines = append(lines, "")
	lines = append(lines, styleDim.Render("  "+strings.Repeat("─", 74)))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s Back to main menu", styleKey.Render("esc")))
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

// viewRunning renders the running benchmark screen
func (m Model) viewRunning() string {
	title := styleTitle.Render(" ⚡ Running Benchmark ")

	var lines []string
	lines = append(lines, "")
	lines = append(lines, title)
	lines = append(lines, "")

	// Current status
	currentServer := m.progress.CurrentServer
	if currentServer == "" {
		currentServer = "Initializing..."
	}

	statusLine := fmt.Sprintf("  Testing: %s",
		styleHighlight.Render(currentServer))
	lines = append(lines, statusLine)

	protoLine := fmt.Sprintf("  Protocol: %s",
		lipgloss.NewStyle().Foreground(colorSecondary).Render(m.progress.CurrentProtocol.String()))
	lines = append(lines, protoLine)
	lines = append(lines, "")

	// Progress bar
	barWidth := 50
	if m.width > 0 && m.width < 70 {
		barWidth = m.width - 20
	}

	percent := m.progress.Percentage
	if percent > 100 {
		percent = 100
	}
	if percent < 0 {
		percent = 0
	}

	filled := int(percent) * barWidth / 100
	empty := barWidth - filled

	progressBar := fmt.Sprintf("  %s%s %s",
		styleProgressBar.Render(strings.Repeat("█", filled)),
		styleProgressTrack.Render(strings.Repeat("░", empty)),
		styleHighlight.Render(fmt.Sprintf("%.0f%%", percent)))
	lines = append(lines, progressBar)
	lines = append(lines, "")

	// Stats
	statsLine := fmt.Sprintf("  %s / %s tests completed",
		styleHighlight.Render(fmt.Sprintf("%d", m.progress.CompletedTests)),
		styleMuted.Render(fmt.Sprintf("%d", m.progress.TotalTests)))
	lines = append(lines, statsLine)

	// Elapsed time
	elapsed := time.Since(m.startTime).Round(time.Second)
	timeLine := fmt.Sprintf("  Elapsed: %s", styleMuted.Render(elapsed.String()))
	lines = append(lines, timeLine)

	lines = append(lines, "")
	lines = append(lines, styleDim.Render("  ───────────────────────────────────"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s Force quit", styleKey.Render("ctrl+c")))
	lines = append(lines, "")

	return strings.Join(lines, "\n")
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
			runner := benchmark.NewRunner(benchConfig)
			ctx := context.Background()
			results, err := runner.Run(ctx)

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
		case done := <-doneChan:
			return done
		default:
		}
		return nil
	}
}

// Run starts the TUI application
func Run(servers []config.Server, settings config.Settings) error {
	p := tea.NewProgram(NewModel(servers, settings), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
