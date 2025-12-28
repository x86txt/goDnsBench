# goDnsBench

A cross-platform DNS benchmarking tool with both TUI and GUI interfaces. Benchmark your DNS servers across multiple protocols including standard DNS, DNS over HTTPS (DoH), DNS over TLS (DoT), and DNS over QUIC (DoQ).

## Features

- 🚀 **Multiple Protocol Support**: DNS, DoH, DoT, and DoQ
- 📊 **Comprehensive Metrics**: Min/Max/Mean/Median/P95/P99 latency statistics
- 🎨 **Dual Interface**: Terminal UI (TUI) and graphical UI (GUI)
- 📈 **Export Results**: Save results in CSV or JSON format
- ⚙️ **Configurable**: Customize timeout, concurrency, and test domains
- 🌐 **Pre-configured Servers**: Ships with popular DNS providers

## Supported Protocols

- **DNS** - Standard DNS over UDP/TCP (port 53)
- **DoH** - DNS over HTTPS (RFC 8484)
- **DoT** - DNS over TLS (RFC 7858)
- **DoQ** - DNS over QUIC (RFC 9250)

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

# Download dependencies
make deps

# Build the binary
make build

# Or build and install to $GOPATH/bin
make install
```

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

### Terminal UI (TUI) Mode

```bash
# Run with default settings
./build/goDnsBench --tui

# Use custom server list
./build/goDnsBench --tui --servers servers.json

# Customize timeout and concurrency
./build/goDnsBench --tui --timeout 2000 --concurrent 5
```

### GUI Mode (Coming Soon)

```bash
# Launch GUI
./build/goDnsBench
```

### Command Line Options

```
Usage of goDnsBench:
  --tui                  Run in terminal UI mode
  --version              Print version and exit
  --servers string       Path to servers JSON/CSV file
  --export-json string   Export results to JSON file
  --export-csv string    Export results to CSV file
  --timeout int          Query timeout in milliseconds (default 1000)
  --concurrent int       Maximum concurrent servers to test (default 10)
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
├── cmd/goDnsBench/        # Main entry point
├── internal/
│   ├── benchmark/         # Benchmark orchestration
│   ├── dns/               # DNS protocol implementations
│   ├── config/            # Configuration management
│   ├── tui/               # Terminal UI
│   └── export/            # Result exporters
├── frontend/              # GUI (Wails + Astro)
├── build/                 # Build outputs
├── servers.json           # Default server list
└── Makefile              # Build automation
```

### Dependencies

- [miekg/dns](https://github.com/miekg/dns) - DNS library
- [quic-go/quic-go](https://github.com/quic-go/quic-go) - QUIC implementation
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - TUI styling

### Make Commands

```bash
make help          # Show all available commands
make deps          # Download dependencies
make build         # Build the binary
make clean         # Clean build files
make test          # Run tests
make fmt           # Format code
make run-tui       # Build and run in TUI mode
make install       # Install to $GOPATH/bin
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
- [x] DoH support
- [x] DoT support
- [x] DoQ support
- [x] Terminal UI (TUI)
- [x] CSV/JSON export
- [ ] GUI with Wails + Astro
- [ ] Result visualization (charts)
- [ ] Advanced filtering and sorting
- [ ] Historical result comparison
- [ ] Custom test domain configuration in UI
- [ ] Windows/Linux platform testing

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Author

Created by [@x86txt](https://github.com/x86txt)

## Support

For issues and questions, please use the [GitHub issue tracker](https://github.com/x86txt/goDnsBench/issues).
