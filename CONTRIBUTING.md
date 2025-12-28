# Contributing to DNS Benchmark Tool

First off, thank you for considering contributing! This project is open to contributions from everyone.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Review Process](#review-process)

## Code of Conduct

This project adheres to a Code of Conduct. By participating, you are expected to uphold this code. Please be respectful and constructive in all interactions.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.22 or later
- **Node.js** 20 or later (for frontend development)
- **Wails CLI** v2.x (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- **golangci-lint** (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)
- **Git**

### macOS-Specific Requirements

```bash
# Install Xcode Command Line Tools
xcode-select --install
```

## Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/REPO_NAME.git
   cd REPO_NAME
   ```

3. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/OWNER/REPO_NAME.git
   ```

4. **Install Go dependencies**
   ```bash
   go mod download
   ```

5. **Install frontend dependencies**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

6. **Verify the setup**
   ```bash
   # Run tests
   go test ./...
   
   # Run linter
   golangci-lint run
   
   # Build the application
   wails build
   ```

## Project Structure

```
.
├── cmd/                    # Application entry points
│   └── dnsbench/
│       └── main.go
├── internal/               # Private application code
│   ├── benchmark/          # Benchmark orchestration
│   ├── dns/                # DNS client implementations
│   ├── config/             # Configuration management
│   ├── tui/                # BubbleTea TUI
│   └── export/             # Export functionality
├── frontend/               # Astro + Tailwind frontend
│   ├── src/
│   └── package.json
├── assets/                 # Static assets
└── build/                  # Build configurations
```

## Making Changes

### Branch Naming

Use descriptive branch names:

- `feature/add-doq-support` - New features
- `fix/timeout-handling` - Bug fixes
- `docs/update-readme` - Documentation
- `refactor/dns-client` - Code refactoring
- `test/benchmark-coverage` - Test additions

### Workflow

1. **Sync with upstream**
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes** (see [Coding Standards](#coding-standards))

4. **Commit with meaningful messages**
   ```bash
   git commit -m "feat: add DNS over QUIC support
   
   - Implement DoQ client using quic-go
   - Add capability detection for DoQ servers
   - Update server config to include DoQ endpoints"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request**

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements

**Examples:**
```
feat(dns): add DNS over HTTPS support
fix(tui): correct table alignment on narrow terminals
docs(readme): add installation instructions for Windows
refactor(benchmark): extract metrics calculation to separate module
```

## Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting (or `goimports`)
- Run `golangci-lint run` before committing
- Write idiomatic Go code

**Naming Conventions:**
```go
// Exported functions: PascalCase
func RunBenchmark() {}

// Unexported functions: camelCase
func calculateMetrics() {}

// Constants: PascalCase or ALL_CAPS for true constants
const DefaultTimeout = 1 * time.Second
const MAX_CONCURRENT_QUERIES = 10

// Interfaces: usually end in -er
type Resolver interface {
    Resolve(domain string) ([]net.IP, error)
}
```

**Error Handling:**
```go
// Always handle errors explicitly
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Use error wrapping for context
if err := validateConfig(cfg); err != nil {
    return fmt.Errorf("invalid configuration: %w", err)
}
```

### Frontend Code (Astro/TypeScript)

- Use TypeScript for all new code
- Follow the existing Tailwind patterns
- Use functional components
- Keep components small and focused

### Documentation

- Add GoDoc comments to all exported functions, types, and packages
- Update README.md for user-facing changes
- Include examples where helpful

```go
// RunBenchmark executes DNS benchmarks against the provided servers.
// It runs queries concurrently across servers but sequentially per server
// to ensure accurate latency measurements.
//
// Example:
//
//	results, err := RunBenchmark(ctx, servers, opts)
//	if err != nil {
//	    log.Fatal(err)
//	}
func RunBenchmark(ctx context.Context, servers []Server, opts Options) ([]Result, error) {
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/dns/...

# Run with verbose output
go test -v ./...

# Run with race detector
go test -race ./...
```

### Writing Tests

- Place tests in `_test.go` files alongside the code
- Use table-driven tests where appropriate
- Mock external dependencies (DNS servers, network calls)
- Aim for >80% coverage on core packages

**Example:**
```go
func TestCalculateP95(t *testing.T) {
    tests := []struct {
        name     string
        latencies []time.Duration
        want     time.Duration
    }{
        {
            name:     "single value",
            latencies: []time.Duration{100 * time.Millisecond},
            want:     100 * time.Millisecond,
        },
        {
            name:     "multiple values",
            latencies: []time.Duration{
                10 * time.Millisecond,
                20 * time.Millisecond,
                30 * time.Millisecond,
                // ... more values
            },
            want: 28 * time.Millisecond, // expected p95
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := calculateP95(tt.latencies)
            if got != tt.want {
                t.Errorf("calculateP95() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Frontend Tests

```bash
cd frontend
npm run test        # Run tests
npm run test:watch  # Watch mode
```

## Submitting Changes

### Before Submitting

1. **Run the full test suite**
   ```bash
   go test -race ./...
   ```

2. **Run the linter**
   ```bash
   golangci-lint run
   ```

3. **Ensure the build succeeds**
   ```bash
   wails build
   ```

4. **Test both interfaces**
   - Run the GUI and verify your changes
   - Run with `--tui` and verify your changes

5. **Update documentation** if needed

### Pull Request Guidelines

- Fill out the PR template completely
- Link related issues using `Fixes #123` or `Relates to #456`
- Keep PRs focused - one feature/fix per PR
- Add screenshots for UI changes
- Ensure CI passes

## Review Process

1. **Automated Checks**: CI will run tests, linting, and CodeQL analysis
2. **Code Review**: A maintainer will review your code
3. **Feedback**: Address any requested changes
4. **Approval**: Once approved, a maintainer will merge

### Review Criteria

- Code quality and readability
- Test coverage
- Documentation
- Performance implications
- Security considerations
- Backward compatibility

## Getting Help

- **Questions?** Open a [Discussion](https://github.com/OWNER/REPO/discussions)
- **Found a bug?** Open an [Issue](https://github.com/OWNER/REPO/issues)
- **Security issue?** See [SECURITY.md](SECURITY.md)

## Recognition

Contributors will be recognized in:
- The project's README
- Release notes for significant contributions
- GitHub's contributors graph

---

Thank you for contributing! 🎉
