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
├── cmd/
│   └── goDnsBench/
│       └── main.go              # CLI/TUI entry point
├── internal/
│   ├── benchmark/
│   │   ├── metrics.go           # Statistical calculations
│   │   ├── results.go           # Result data structures
│   │   └── runner.go            # Benchmark orchestration
│   ├── dns/
│   │   ├── types.go             # Common DNS types
│   │   ├── client.go            # Standard DNS client
│   │   ├── doh.go               # DNS over HTTPS
│   │   ├── dot.go               # DNS over TLS
│   │   ├── doq.go               # DNS over QUIC
│   │   └── capability.go        # Protocol detection
│   ├── config/
│   │   ├── servers.go           # Server management
│   │   ├── loader.go            # JSON/CSV loading
│   │   └── settings.go          # User preferences
│   ├── tui/
│   │   └── app.go               # BubbleTea application
│   └── export/
│       ├── csv.go               # CSV exporter
│       └── json.go              # JSON exporter
├── frontend/                    # Astro + Tailwind GUI
│   ├── src/
│   │   ├── layouts/
│   │   │   └── Layout.astro     # Base HTML layout
│   │   ├── pages/
│   │   │   └── index.astro      # Main application page
│   │   ├── scripts/
│   │   │   └── app.ts           # Frontend logic & Wails integration
│   │   └── wailsjs/             # Generated Wails bindings
│   ├── astro.config.mjs         # Astro configuration
│   ├── tailwind.config.js       # Tailwind with theme
│   ├── tsconfig.json            # TypeScript config
│   └── package.json             # Frontend dependencies
├── build/                       # Build outputs
├── app.go                       # Wails application backend
├── gui.go                       # GUI launcher
├── main.go                      # Main entry point (GUI/TUI router)
├── wails.json                   # Wails configuration
├── servers.json                 # Default DNS servers
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
├── .gitignore                   # Git ignore rules
├── README.md                    # User documentation
└── CONTEXT.md                   # This file

```

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
- Node.js 18+ and npm (for GUI)
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Build Commands

```bash
# Download all dependencies
make deps
make deps-frontend

# Build CLI/TUI only
make build-cli

# Build GUI only
make build-gui

# Build both
make build

# Run in TUI mode
make run-tui
./build/goDnsBench-cli --tui

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

**TUI Development:**
```bash
make build-cli
./build/goDnsBench-cli --tui
```

**GUI Development:**
```bash
make dev  # Starts Wails dev server with hot reload
```

**Production Build:**
```bash
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
- `GetServers() []config.Server` - Get current server list
- `LoadServersFromFile(path string) error` - Load from JSON/CSV
- `AddServer(server config.Server) error` - Add custom server
- `RemoveServer(name string)` - Remove server by name
- `ResetToDefaults()` - Reset to default servers
- `RefreshServerList() error` - Fetch from configured URL

### Settings
- `GetSettings() config.Settings` - Get current settings
- `UpdateSettings(settings config.Settings) error` - Update and save

### Benchmarking
- `RunBenchmark(protocols []string) error` - Execute benchmark
- `GetResults() *benchmark.BenchmarkResults` - Get last results

### Export
- `ExportResultsJSON(path string) error` - Export to JSON
- `ExportResultsCSV(path string) error` - Export to CSV

### Utilities
- `CheckServerCapabilities(server config.Server) map[string]bool` - Test protocols

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

### Phase 3: Enhancement (Current)
- [x] Visualization (Chart.js integration)
- [ ] Advanced filtering/sorting
- [ ] Result history
- [ ] Custom domain configuration in UI

### Phase 4: Polish
- [ ] Comprehensive testing
- [ ] Performance optimization
- [ ] Platform-specific builds
- [ ] Installer packages

---

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

*Last Updated: December 28, 2025*
*Version: 0.1.0*
