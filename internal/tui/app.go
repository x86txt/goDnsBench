package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/x86txt/goDnsBench/internal/benchmark"
	"github.com/x86txt/goDnsBench/internal/config"
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

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "s":
			if m.currentScreen == ScreenMain {
				m.currentScreen = ScreenSettings
			}
		case "r":
			if m.currentScreen == ScreenMain {
				m.currentScreen = ScreenRunning
				// TODO: Start benchmark
			}
		case "esc":
			m.currentScreen = ScreenMain
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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

	title := titleStyle.Render("Benchmark Results")

	// TODO: Format results in a table
	content := fmt.Sprintf(`
Total Duration: %v
Tests Completed: %d

[esc] Back to main menu
`, m.results.Duration(), len(m.results.Results))

	return title + "\n" + content
}

// viewRunning renders the running benchmark screen
func (m Model) viewRunning() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0066ff")).
		MarginBottom(1)

	title := titleStyle.Render("Running Benchmark...")

	content := `
Benchmarking in progress...

This feature will be implemented in the next phase.

[esc] Back to main menu
`

	return title + "\n" + content
}

// Run starts the TUI application
func Run(servers []config.Server, settings config.Settings) error {
	p := tea.NewProgram(NewModel(servers, settings))
	_, err := p.Run()
	return err
}
