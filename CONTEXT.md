# goDnsBench - Project Context

## Project Overview

goDnsBench is a cross-platform DNS benchmarking tool with both TUI (Terminal UI) and GUI (Graphical UI) interfaces. It benchmarks DNS servers across multiple protocols: standard DNS, DoH (DNS over HTTPS), DoT (DNS over TLS), and DoQ (DNS over QUIC).

**Version**: 0.1.0
**Author**: x86txt
**Repository**: https://github.com/x86txt/goDnsBench

---

## Technical Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Language | Go 1.21+ | Core application logic |
| TUI | BubbleTea + Lipgloss | Terminal interface |
| GUI Framework | Wails v2 | Desktop application wrapper |
| Frontend | Astro 4.x | Static site generation |
| Styling | Tailwind CSS 3.x | UI styling and theming |
| Charts | Chart.js 4.x | Data visualization |
| DNS Library | miekg/dns | DNS protocol operations |
| QUIC | quic-go/quic-go | DoQ support |

---

## Design Theme: Tech Innovation

### Colors
- **Electric Blue** (`#0066ff`): Primary accent, buttons, highlights
- **Neon Cyan** (`#00ffff`): Secondary accent, headings
- **Dark Gray** (`#1e1e1e`): Primary background
- **Light Gray** (`#2a2a2a`): Secondary background
- **Lighter Gray** (`#3a3a3a`): Tertiary background, borders

### Typography
- **Headers**: DejaVu Sans Bold
- **Body**: DejaVu Sans
- **Code/Mono**: DejaVu Sans Mono

---

## Features

### DNS Protocols Supported ✅
- **DNS** - Standard DNS over UDP/TCP (port 53)
- **DoH** - DNS over HTTPS (RFC 8484)
- **DoT** - DNS over TLS (RFC 7858)
- **DoQ** - DNS over QUIC (RFC 9250)

### Query Types (10 per server)
- **A Records** (4): Standard hostname resolution
- **MX Records** (2): Mail server lookup
- **TXT Records** (2): Text record lookup
- **DNSSEC** (2): Queries checking Authenticated Data (AD) flag

### Test Domains (Default)
- **A Records**: google.com, cloudflare.com, amazon.com, microsoft.com
- **MX Records**: gmail.com, microsoft.com
- **TXT Records**: _dmarc.google.com, google.com
- **DNSSEC**: cloudflare.com, google.com

### Metrics Captured
- **Latency**: Min, Max, Mean, Median, P95, P99
- **Reliability**: Success rate, failure count
- **Total queries**: Per server per protocol

### Results Export
- **CSV**: Tabular format for spreadsheet analysis
- **JSON**: Structured format for programmatic access

---

## Project Structure

```
goDnsBench/
├── main.go                      # Main entry point (routes GUI/TUI/headless modes)
├── app.go                       # Wails application backend (App struct + API methods)
├── gui.go                       # Wails GUI launcher
├── api_dto.go                   # Data Transfer Objects for Wails API layer
├── cmd/
│   └── goDnsBench/              # (Legacy directory, not used in unified binary)
├── internal/
│   ├── benchmark/
│   │   ├── metrics.go           # Statistical calculations
│   │   ├── metrics_test.go      # Unit tests for metrics
│   │   ├── results.go           # Result data structures
│   │   ├── results_test.go      # Unit tests for results
│   │   └── runner.go            # Benchmark orchestration
│   ├── cli/
│   │   └── headless.go          # Headless/CLI benchmark execution
│   ├── dns/
│   │   ├── types.go             # Common DNS types
│   │   ├── client.go            # Standard DNS client
│   │   ├── doh.go               # DNS over HTTPS
│   │   ├── dot.go               # DNS over TLS
│   │   ├── doq.go               # DNS over QUIC
│   │   └── capability.go        # Protocol detection
│   ├── config/
│   │   ├── servers.go           # Server management
│   │   ├── servers_test.go      # Unit tests for servers
│   │   ├── loader.go            # JSON/CSV loading
│   │   ├── loader_test.go       # Unit tests for loader
│   │   ├── settings.go          # User preferences
│   │   └── settings_test.go     # Unit tests for settings
│   ├── tui/
│   │   └── app.go               # BubbleTea application
│   └── export/
│       ├── csv.go               # CSV exporter
│       ├── csv_test.go          # Unit tests for CSV export
│       ├── json.go              # JSON exporter
│       └── json_test.go         # Unit tests for JSON export
├── frontend/                    # Astro + Tailwind GUI
│   ├── src/
│   │   ├── layouts/
│   │   │   └── Layout.astro     # Base HTML layout
│   │   ├── pages/
│   │   │   └── index.astro      # Main application page
│   │   ├── scripts/
│   │   │   ├── init.ts          # Main initialization
│   │   │   ├── api/             # Wails API client
│   │   │   │   ├── client.ts    # API method wrappers
│   │   │   │   └── types.ts     # TypeScript types
│   │   │   ├── charts/          # Chart.js visualizations
│   │   │   │   ├── latency.ts   # Latency chart rendering
│   │   │   │   └── comparison.ts # Comparison charts
│   │   │   └── ui/              # UI modules
│   │   │       ├── dialogs.ts   # Dialog management
│   │   │       ├── elements.ts  # DOM element utilities
│   │   │       ├── filtering.ts # Result filtering logic
│   │   │       ├── history.ts   # Result history management
│   │   │       ├── progress.ts  # Progress bar UI
│   │   │       ├── results.ts   # Results table rendering
│   │   │       ├── servers.ts   # Server list UI
│   │   │       └── settings.ts  # Settings modal UI
│   │   └── wailsjs/             # Generated Wails bindings
│   ├── astro.config.mjs         # Astro configuration
│   ├── tailwind.config.js       # Tailwind with theme
│   ├── tsconfig.json            # TypeScript config
│   └── package.json             # Frontend dependencies
├── build/                       # Build outputs
├── wails.json                   # Wails configuration
├── servers.json                 # Default DNS servers
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
├── .gitignore                   # Git ignore rules
├── README.md                    # User documentation
└── CONTEXT.md                   # This file
```

### Root-Level Go Files

The root directory contains several `.go` files that are part of the `main` package:

- **`main.go`**: Main entry point that parses command-line flags and routes execution to GUI, TUI, or headless mode
- **`app.go`**: Wails application backend containing the `App` struct and all methods exposed to the frontend via Wails bindings
- **`gui.go`**: Wails GUI launcher that configures and starts the Wails application
- **`api_dto.go`**: Data Transfer Objects (DTOs) and conversion functions between internal Go types and GUI-safe JSON-compatible types

These files are kept in the root because:
1. They're part of the `main` package, which is required for the application entry point
2. Wails framework requires the main application struct (`App`) to be accessible from the root package for bindings
3. They serve as the integration/bridge layer between Wails (GUI framework) and the internal packages
4. This is a common pattern in Wails applications

While this structure differs from some Go projects where all code lives in `cmd/` or `internal/`, it's appropriate for Wails-based desktop applications where the root package acts as the integration layer.

---

## Server Configuration

### Default DNS Servers (13)
1. Cloudflare Primary (1.1.1.1)
2. Cloudflare Secondary (1.0.0.1)
3. Google Primary (8.8.8.8)
4. Google Secondary (8.8.4.4)
5. Quad9 (9.9.9.9)
6. OpenDNS Primary (208.67.222.222)
7. OpenDNS Secondary (208.67.220.220)
8. AdGuard DNS (94.140.14.14)
9. NextDNS
10. CleanBrowsing (185.228.168.9)
11. Comodo Secure DNS (8.26.56.26)
12. Level3 Primary (4.2.2.1)
13. Level3 Secondary (4.2.2.2)

### Server JSON Format
```json
{
  "servers": [
    {
      "name": "Cloudflare Primary",
      "dns": "1.1.1.1",
      "doh": "https://cloudflare-dns.com/dns-query",
      "dot": "1.1.1.1:853",
      "doq": "1.1.1.1:8853"
    }
  ]
}
```

### Server CSV Format
```csv
name,dns,doh,dot,doq
Cloudflare Primary,1.1.1.1,https://cloudflare-dns.com/dns-query,1.1.1.1:853,1.1.1.1:8853
```

---

## Configuration & Settings

### Settings File Location
`~/.config/goDnsBench/settings.json`

### Default Settings
```json
{
  "queryTimeout": 1000000000,        // 1 second in nanoseconds
  "maxConcurrent": 10,                // Max concurrent server tests
  "serverListUrl": "https://raw.githubusercontent.com/x86txt/goDnsBench/main/servers.json",
  "enabledProtocols": ["DNS", "DoH", "DoT", "DoQ"],
  "selectedDomains": {
    "a": ["google.com", "cloudflare.com", "amazon.com", "microsoft.com"],
    "mx": ["gmail.com", "microsoft.com"],
    "txt": ["_dmarc.google.com", "google.com"],
    "dnssec": ["cloudflare.com", "google.com"]
  }
}
```

---

## Benchmarking Methodology

### Execution Model
1. **Server Concurrency**: Up to 10 servers tested concurrently (configurable)
2. **Query Execution**: Sequential within each server (one query at a time per server)
3. **Timeout**: 1 second default per query (configurable)
4. **Protocol Detection**: Automatically skips unsupported protocols

### Metric Calculations
- **Min/Max**: Fastest and slowest query times
- **Mean**: Average of all successful queries
- **Median**: Middle value of sorted latencies
- **P95**: 95th percentile (95% of queries faster than this)
- **P99**: 99th percentile (99% of queries faster than this)
- **Success Rate**: (Successful queries / Total queries) × 100%

---

## Build & Development

### Prerequisites
- Go 1.21 or higher
- Bun (for frontend): `curl -fsSL https://bun.sh/install | bash`
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Build Commands

```bash
# Download all dependencies
make deps          # Go dependencies
make deps-frontend # Frontend dependencies (Bun)

# Build unified binary (includes both GUI and TUI)
make build

# Build CLI/TUI only (for testing)
make build-cli

# Build GUI only (production build)
make build-gui

# Run in TUI mode
make run-tui
# Or: ./build/goDnsBench --tui

# Run GUI (default mode)
./build/goDnsBench

# Run GUI in development mode (hot reload)
make dev

# Build and run GUI
make run-gui

# Cross-platform builds
make build-darwin    # macOS (Intel + ARM)
make build-linux     # Linux
make build-windows   # Windows
make build-all       # All platforms
```

### Development Workflow

**Unified Binary:**
The application uses a single binary that supports both GUI and TUI modes. By default, running the binary launches the GUI. Use the `--tui` flag for terminal mode.

**TUI Development:**
```bash
make build-cli
./build/goDnsBench-cli --tui
# Or with unified binary:
make build-gui
./build/goDnsBench --tui
```

**GUI Development:**
```bash
# Install frontend dependencies first
make deps-frontend

# Start Wails dev server with hot reload
make dev
```

**Production Build:**
```bash
# Build unified binary (includes GUI)
make build-gui  # Creates optimized production binary
```

---

## Dependencies

### Go Dependencies
- `github.com/miekg/dns` - DNS protocol operations
- `github.com/quic-go/quic-go` - QUIC/DoQ support
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - TUI styling
- `github.com/charmbracelet/bubbles` - TUI components
- `github.com/wailsapp/wails/v2` - GUI framework

### Frontend Dependencies
- `bun` - JavaScript runtime and package manager
- `astro` - Static site framework
- `@astrojs/tailwind` - Tailwind integration
- `tailwindcss` - CSS framework
- `chart.js` - Charting library

---

## Architecture

### Flow Diagrams

#### GUI Mode Flow
```
User launches goDnsBench (no --tui flag)
    ↓
main.go → LaunchGUI()
    ↓
Wails starts with app.go backend
    ↓
Frontend loads (Astro + Tailwind + TypeScript)
    ↓
User selects protocols & servers
    ↓
Click "Run Benchmark" → App.RunBenchmark()
    ↓
benchmark.Runner orchestrates execution
    ↓
DNS clients query concurrently (max 10 servers)
    ↓
Metrics calculated & returned
    ↓
Frontend displays results in table + chart
    ↓
User exports to CSV/JSON
```

#### TUI Mode Flow
```
User launches goDnsBench --tui
    ↓
cmd/goDnsBench/main.go → tui.Run()
    ↓
BubbleTea initializes
    ↓
User navigates with keyboard
    ↓
Select options & run benchmark
    ↓
Results displayed in terminal
```

---

## API (Wails Backend Methods)

The following Go methods are exposed to the frontend via Wails:

### Server Management
- `GetServers() []ServerDTO` - Get current server list
- `LoadServersFromFile(filepath string) error` - Load from JSON/CSV file
- `LoadServersFromContent(content string, fileType string) error` - Load from file content
- `AddServer(server ServerDTO) error` - Add custom server
- `RemoveServer(name string)` - Remove server by name
- `ResetToDefaults()` - Reset to default servers
- `RefreshServerList() error` - Fetch from configured URL
- `SetSelectedServers(names []string) error` - Set selected servers for benchmarking
- `GetSelectedServers() []string` - Get currently selected server names

### Settings
- `GetSettings() config.Settings` - Get current settings
- `UpdateSettings(settings config.Settings) error` - Update and save

### Benchmarking
- `RunBenchmark(protocols []string) error` - Execute benchmark with selected servers and protocols
- `GetResults() *BenchmarkResultsDTO` - Get last benchmark results

### Export
- `ExportResultsJSON(path string) error` - Export to JSON
- `ExportResultsCSV(path string) error` - Export to CSV

### Utilities
- `CheckServerCapabilities(server ServerDTO) map[string]bool` - Test which protocols a server supports

---

## Roadmap

### Phase 1: Core Foundation ✅
- [x] DNS protocol implementations
- [x] Benchmark engine
- [x] Metrics collection
- [x] Configuration management
- [x] Default server list

### Phase 2: Interfaces ✅
- [x] Terminal UI (TUI)
- [x] GUI framework (Wails + Astro)
- [x] Export functionality (CSV/JSON)

### Phase 3: Enhancement ✅
- [x] Visualization (Chart.js integration)
- [x] Server selection in GUI
- [x] Result filtering (by server, protocol, success rate)
- [x] Result history management
- [x] Headless export mode (CLI flags)
- [x] Unit tests for core packages (benchmark, config, export)
- [ ] Custom domain configuration in UI
- [ ] Historical result comparison UI

### Phase 4: Polish
- [ ] Comprehensive testing
- [ ] Performance optimization
- [ ] Platform-specific builds
- [ ] Installer packages

---

## Recent Updates

### Unit Tests (Stage 4)
Added comprehensive unit tests for core packages:
- **internal/benchmark**: 17 test functions covering metrics calculations, result structures, filtering, and comparisons
- **internal/config**: 15 test functions covering settings, server management, and JSON/CSV loading
- **internal/export**: 8 test functions covering CSV and JSON export functionality
- **Total**: 40+ test functions, 100+ test cases, 100% coverage of public APIs

### GUI Enhancements
- **Server Selection**: Users can select specific servers to benchmark from the GUI
- **Result Filtering**: Filter results by server name, protocol, and minimum success rate
- **Result History**: Save and manage benchmark result history
- **Chart Visualizations**: Interactive latency and comparison charts using Chart.js
- **Modular Frontend**: Refactored frontend into modular components (api/, charts/, ui/)

### Headless Mode
- Added `--export-json` and `--export-csv` flags for command-line benchmarking
- Supports all configuration options (protocols, timeout, concurrency, server list)
- Enables automation and integration with CI/CD pipelines

## Known Limitations

1. **Protocol Support**: Not all DNS servers support all protocols (handled gracefully)
2. **DNSSEC Validation**: Currently only checks AD flag, doesn't validate chain
3. **Timeout**: Fixed timeout per query, not adaptive
4. **IPv6**: Currently focused on IPv4, IPv6 support planned
5. **Platform Testing**: Primary development on macOS, Windows/Linux less tested

---

## Contributing

Contributions are welcome! Areas of interest:
- Additional DNS protocols (DNSCrypt, etc.)
- UI/UX improvements
- Performance optimizations
- Platform-specific enhancements
- Documentation improvements

---

## License

MIT License - See LICENSE file

---

## Support & Issues

**GitHub Issues**: https://github.com/x86txt/goDnsBench/issues
**Documentation**: See README.md for user guide

---

---

*Last Updated: January 2025*
*Version: 0.1.0*
