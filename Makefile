# goDnsBench Makefile

# Variables
BINARY_NAME=goDnsBench
MAIN_PATH=.
BUILD_DIR=./build
VERSION?=0.1.0
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build targets
.PHONY: all build build-gui build-frontend clean test deps deps-frontend fmt help tidy run-tui run-gui dev

all: deps fmt build

## build: Build unified binary (includes both GUI and TUI, defaults to GUI)
build: build-frontend
	@echo "Building unified $(BINARY_NAME) binary (GUI + TUI)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary build complete: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  Run without flags to launch GUI: ./$(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  Run with --tui flag for terminal mode: ./$(BUILD_DIR)/$(BINARY_NAME) --tui"

## build-frontend: Build the frontend to create dist directory
build-frontend: deps-frontend
	@echo "Building frontend..."
	@cd frontend && bun run build
	@echo "Frontend build complete"

## build-gui: Build desktop application bundle with Wails (for distribution)
build-gui: build-frontend
	@echo "Building $(BINARY_NAME) desktop application bundle with Wails..."
	@command -v wails >/dev/null 2>&1 || { echo "Wails CLI not found. Install with: go install github.com/wailsapp/wails/v2/cmd/wails@latest"; exit 1; }
	wails build
	@echo "Desktop application bundle build complete"

## clean: Clean build files
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## deps: Download Go dependencies
deps:
	@echo "Downloading Go dependencies..."
	$(GOGET) github.com/miekg/dns
	$(GOGET) github.com/quic-go/quic-go
	$(GOGET) github.com/charmbracelet/bubbletea
	$(GOGET) github.com/charmbracelet/lipgloss
	$(GOGET) github.com/charmbracelet/bubbles
	$(GOGET) github.com/wailsapp/wails/v2@latest
	$(GOMOD) tidy
	@echo "Go dependencies downloaded"

## deps-frontend: Install frontend dependencies
deps-frontend:
	@echo "Installing frontend dependencies..."
	@cd frontend && bun install
	@echo "Frontend dependencies installed"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Formatting complete"

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	$(GOMOD) tidy
	@echo "Tidy complete"

## run-tui: Build and run in TUI mode
run-tui: build
	@echo "Running in TUI mode..."
	$(BUILD_DIR)/$(BINARY_NAME) --tui

## run-gui: Build and run in GUI mode (default)
run-gui: build
	@echo "Running in GUI mode (default)..."
	$(BUILD_DIR)/$(BINARY_NAME)

## dev: Run Wails in development mode
dev: deps-frontend
	@echo "Running Wails dev server..."
	@command -v wails >/dev/null 2>&1 || { echo "Wails CLI not found. Install with: go install github.com/wailsapp/wails/v2/cmd/wails@latest"; exit 1; }
	wails dev

## install: Install the unified binary to $GOPATH/bin
install: build-frontend
	@echo "Installing unified $(BINARY_NAME) binary..."
	$(GOCMD) install $(LDFLAGS) $(MAIN_PATH)
	@echo "Installation complete"

## help: Show this help message
help:
	@echo "goDnsBench Makefile Commands:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# Platform-specific builds
.PHONY: build-linux build-darwin build-windows build-all

## build-linux: Build unified binary for Linux
build-linux: build-frontend
	@echo "Building unified binary for Linux..."
	@mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/linux/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Linux build complete"

## build-darwin: Build unified binary for macOS
build-darwin: build-frontend
	@echo "Building unified binary for macOS..."
	@mkdir -p $(BUILD_DIR)/darwin
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/darwin/$(BINARY_NAME) $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/darwin/$(BINARY_NAME)-arm64 $(MAIN_PATH)
	@echo "macOS build complete"

## build-windows: Build unified binary for Windows
build-windows: build-frontend
	@echo "Building unified binary for Windows..."
	@mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/windows/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo "Windows build complete"

## build-all: Build for all platforms
build-all: build-linux build-darwin build-windows
	@echo "All platform builds complete"
