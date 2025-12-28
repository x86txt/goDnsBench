# goDnsBench

A cross-platform DNS benchmarking tool with both TUI and GUI interfaces. Benchmark your DNS servers across multiple protocols including standard DNS, DNS over HTTPS (DoH), DNS over TLS (DoT), and DNS over QUIC (DoQ).

## Features

- 🚀 **Multiple Protocol Support**: DNS, DNS over HTTPS (DoH), DNS over TLS (DoT), and DNS over QUIC (DoQ)
- 📊 **Comprehensive Metrics**: Min/Max/Mean/Median/P95/P99 latency statistics
- 🎨 **Dual Interface**: Terminal UI (TUI) and graphical UI (GUI)
- 📈 **Export Results**: Save results in CSV or JSON format (interactive or headless)
- ⚙️ **Configurable**: Customize timeout, concurrency, and test domains
- 🌐 **Pre-configured Servers**: Ships with popular DNS providers
- 🔍 **Filtering & Sorting**: Filter results by server, protocol, and success rate
- 📉 **Visualization**: Interactive charts for latency comparison and analysis
- 🎯 **Server Selection**: Select specific servers to benchmark from the GUI
- 📚 **Result History**: Save and compare benchmark results over time
- 🧪 **Headless Mode**: Run benchmarks from command line with export flags

## Supported Protocols

- **DNS** - Standard DNS over UDP/TCP (port 53)
- **DNS over HTTPS (DoH)** - DNS over HTTPS (RFC 8484)
- **DNS over TLS (DoT)** - DNS over TLS (RFC 7858)
- **DNS over QUIC (DoQ)** - DNS over QUIC (RFC 9250)

## Query Types

Each server is tested with 10 queries:
- 4x A Record lookups
- 2x MX Record lookups
- 2x TXT Record lookups
- 2x DNSSEC queries (checking AD flag)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/x86txt/goDnsBench.git
cd goDnsBench

# Download dependencies (Go and frontend)
make deps
make deps-frontend

# Build the unified binary (includes GUI)
make build

# Or build and install to $GOPATH/bin
make install
```

**Note**: The unified binary includes both GUI and TUI interfaces. GUI mode launches by default; use `--tui` flag for terminal mode.

### Build for Specific Platforms

```bash
# macOS
make build-darwin

# Linux
make build-linux

# Windows
make build-windows

# All platforms
make build-all
```

## Usage

goDnsBench uses a unified binary that supports both GUI and TUI modes. By default, it launches the graphical interface. Use the `--tui` flag to run in terminal mode.

### GUI Mode (Default)

```bash
# Launch GUI (default behavior)
./build/goDnsBench

# Use custom server list
./build/goDnsBench --servers servers.json

# Customize timeout and concurrency
./build/goDnsBench --timeout 2000 --concurrent 5
```

### Terminal UI (TUI) Mode

```bash
# Run with default settings
./build/goDnsBench --tui

# Use custom server list
./build/goDnsBench --tui --servers servers.json

# Customize timeout and concurrency
./build/goDnsBench --tui --timeout 2000 --concurrent 5
```

### Command Line Options

```
Usage of goDnsBench:
  --tui                  Run in terminal UI mode
  --version              Print version and exit
  --servers string       Path to servers JSON/CSV file
  --export-json string   Export results to JSON file (runs headless benchmark)
  --export-csv string    Export results to CSV file (runs headless benchmark)
  --protocols string     Comma-separated list of protocols to test (default: "DNS,DoH,DoT,DoQ")
  --timeout int          Query timeout in milliseconds (default 1000)
  --concurrent int       Maximum concurrent servers to test (default 10)
```

### Headless Export Mode

You can run benchmarks directly from the command line and export results without opening a UI:

```bash
# Export to JSON
./build/goDnsBench --export-json results.json

# Export to CSV with specific protocols
./build/goDnsBench --export-csv results.csv --protocols "DNS,DoH"

# Use custom server list and timeout
./build/goDnsBench --export-json results.json --servers myservers.json --timeout 2000
```

## Server Configuration

### JSON Format

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

### CSV Format

```csv
name,dns,doh,dot,doq
Cloudflare Primary,1.1.1.1,https://cloudflare-dns.com/dns-query,1.1.1.1:853,1.1.1.1:8853
Google Primary,8.8.8.8,https://dns.google/dns-query,8.8.8.8:853,
```

### Default DNS Servers

goDnsBench comes pre-configured with popular DNS providers:

- Cloudflare (1.1.1.1, 1.0.0.1)
- Google (8.8.8.8, 8.8.4.4)
- Quad9 (9.9.9.9)
- OpenDNS (208.67.222.222, 208.67.220.220)
- AdGuard DNS
- NextDNS
- CleanBrowsing
- Comodo Secure DNS
- Level3 (4.2.2.1, 4.2.2.2)

## Configuration

Configuration is stored in `~/.config/goDnsBench/settings.json`

```json
{
  "queryTimeout": 1000000000,
  "maxConcurrent": 10,
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

## Development

### Project Structure

```
goDnsBench/
├── main.go                # Main entry point (routes GUI/TUI/headless)
├── app.go                 # Wails application backend (GUI API)
├── gui.go                 # Wails GUI launcher
├── api_dto.go             # Data Transfer Objects for Wails API
├── cmd/goDnsBench/        # (Legacy, not used in unified binary)
├── internal/
│   ├── benchmark/         # Benchmark orchestration
│   │   ├── runner.go      # Benchmark execution
│   │   ├── metrics.go     # Statistical calculations
│   │   ├── results.go     # Result data structures
│   │   └── *_test.go      # Unit tests
│   ├── cli/               # CLI/headless mode
│   ├── dns/               # DNS protocol implementations
│   ├── config/            # Configuration management
│   │   └── *_test.go      # Unit tests
│   ├── tui/               # Terminal UI (BubbleTea)
│   └── export/            # Result exporters (CSV/JSON)
│       └── *_test.go      # Unit tests
├── frontend/              # GUI (Wails + Astro)
│   ├── src/
│   │   ├── scripts/
│   │   │   ├── api/       # Wails API client
│   │   │   ├── charts/    # Chart.js visualizations
│   │   │   └── ui/        # UI modules (dialogs, filtering, results, etc.)
│   │   ├── pages/         # Astro pages
│   │   └── layouts/       # Astro layouts
│   └── package.json
├── build/                 # Build outputs
├── servers.json           # Default server list
└── Makefile              # Build automation
```

**Note**: The root-level `.go` files (`main.go`, `app.go`, `gui.go`, `api_dto.go`) are part of the `main` package and serve as the entry point and Wails integration layer. They're kept in the root because Wails requires the main application struct to be accessible from the root package.

### Dependencies

**Go Dependencies:**
- [miekg/dns](https://github.com/miekg/dns) - DNS library
- [quic-go/quic-go](https://github.com/quic-go/quic-go) - QUIC implementation
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - TUI styling
- [wailsapp/wails](https://github.com/wailsapp/wails) - GUI framework

**Frontend Dependencies:**
- [Bun](https://bun.sh/) - JavaScript runtime and package manager
- [Astro](https://astro.build/) - Frontend framework
- [Tailwind CSS](https://tailwindcss.com/) - Styling
- [Chart.js](https://www.chartjs.org/) - Data visualization

### Prerequisites

- Go 1.21 or higher
- Bun (for frontend development): `curl -fsSL https://bun.sh/install | bash`
- Wails CLI (for GUI builds): `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Make Commands

```bash
make help          # Show all available commands
make deps          # Download Go dependencies
make deps-frontend # Install frontend dependencies (Bun)
make build         # Build unified binary (GUI + TUI)
make build-cli     # Build CLI/TUI only
make build-gui     # Build GUI application with Wails
make clean         # Clean build files
make test          # Run all tests
make fmt           # Format code
make tidy          # Tidy go.mod dependencies
make run-tui       # Build and run in TUI mode
make run-gui       # Build and run in GUI mode
make dev           # Run Wails dev server (hot reload)
make install       # Install to $GOPATH/bin
```

### Development Workflow

**GUI Development:**
```bash
# Install dependencies
make deps-frontend

# Run in development mode with hot reload
make dev
```

**TUI Development:**
```bash
# Build and run
make run-tui
```

**Production Build:**
```bash
# Build unified binary
make build-gui
```

## Benchmarking Methodology

1. **Sequential Queries**: Queries to the same server run sequentially (one at a time)
2. **Concurrent Servers**: Multiple servers can be tested concurrently (default: 10)
3. **Timeout**: Default 1 second per query (configurable)
4. **Protocol Detection**: Automatically skips protocols not supported by a server

## Metrics Explained

- **Min**: Fastest query response time
- **Max**: Slowest query response time
- **Mean**: Average response time across all queries
- **Median**: Middle value when all response times are sorted
- **P95**: 95th percentile - 95% of queries were faster than this
- **P99**: 99th percentile - 99% of queries were faster than this
- **Success Rate**: Percentage of queries that completed successfully

## Roadmap

- [x] Core DNS benchmarking engine
- [x] Standard DNS support
- [x] DNS over HTTPS (DoH) support
- [x] DNS over TLS (DoT) support
- [x] DNS over QUIC (DoQ) support
- [x] Terminal UI (TUI)
- [x] CSV/JSON export
- [x] GUI with Wails + Astro
- [x] Result visualization (charts)
- [x] Server selection in GUI
- [x] Result filtering by server, protocol, and success rate
- [x] Headless export mode (command-line export flags)
- [x] Unit tests for core packages (benchmark, config, export)
- [ ] Historical result comparison
- [ ] Custom test domain configuration in UI
- [ ] Windows/Linux platform testing
- [ ] Performance optimizations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Author

Created by [@x86txt](https://github.com/x86txt)

## Support

For issues and questions, please use the [GitHub issue tracker](https://github.com/x86txt/goDnsBench/issues).
