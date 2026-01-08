# Makefile for Aether Vault CLI

# Variables
BINARY_NAME=aether-vault-cli
BUILD_DIR=bin
VERSION ?= 1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Default target
.PHONY: all
all: clean build

# Build the CLI binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	
	# Linux ARM64
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	
	# macOS AMD64
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows AMD64
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "Cross-platform builds complete"
	@ls -la $(BUILD_DIR)/

# Install to system PATH
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

# Install to user's local bin
.PHONY: install-local
install-local: build
	@echo "Installing $(BINARY_NAME) to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) ~/.local/bin/
	@echo "Installation complete"
	@echo "Make sure ~/.local/bin is in your PATH"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p coverage
	@go test -v -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report generated: coverage/coverage.html"

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Run formatter
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

# Check for vulnerabilities
.PHONY: security
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not installed. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf coverage
	@go clean
	@echo "Clean complete"

# Clean all (including dependencies)
.PHONY: clean-all
clean-all: clean
	@echo "Cleaning all..."
	@rm -rf vendor/

# Update dependencies
.PHONY: deps
deps:
	@echo "Updating dependencies..."
	@go mod tidy
	@go mod download

# Verify dependencies
.PHONY: deps-verify
deps-verify:
	@echo "Verifying dependencies..."
	@go mod verify

# Download dependencies
.PHONY: deps-download
deps-download:
	@echo "Downloading dependencies..."
	@go mod download

# Generate release artifacts
.PHONY: release
release: clean test lint build-all
	@echo "Creating release artifacts..."
	@mkdir -p release
	@cd $(BUILD_DIR) && \
		for file in $(BINARY_NAME)-*; do \
			if [[ $$file == *.exe ]]; then \
				zip ../release/$$file.zip $$file; \
			else \
				tar -czf ../release/$$file.tar.gz $$file; \
			fi; \
		done
	@echo "Release artifacts created in release/"

# Run the CLI
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run the CLI with development config
.PHONY: run-dev
run-dev:
	@echo "Running $(BINARY_NAME) in development mode..."
	@go run . --help

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t aether-vault-cli:$(VERSION) .
	@docker tag aether-vault-cli:$(VERSION) aether-vault-cli:latest
	@echo "Docker image built"

# Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	@docker run --rm -it aether-vault-cli:latest

# Development targets
.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Development tools installed"

# CI targets
.PHONY: ci
ci: fmt lint security test-coverage
	@echo "CI checks completed"

# Help target
.PHONY: help
help:
	@echo "Aether Vault CLI - Makefile targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build         Build the CLI binary"
	@echo "  build-all     Build for all platforms"
	@echo "  clean         Clean build artifacts"
	@echo "  clean-all     Clean everything including dependencies"
	@echo ""
	@echo "Install targets:"
	@echo "  install       Install to /usr/local/bin"
	@echo "  install-local Install to ~/.local/bin"
	@echo ""
	@echo "Test targets:"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo ""
	@echo "Code quality:"
	@echo "  fmt           Format code"
	@echo "  lint          Run linter"
	@echo "  security      Check for vulnerabilities"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps          Update dependencies"
	@echo "  deps-verify   Verify dependencies"
	@echo "  deps-download Download dependencies"
	@echo ""
	@echo "Release:"
	@echo "  release       Generate release artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Run Docker container"
	@echo ""
	@echo "Development:"
	@echo "  dev-setup     Install development tools"
	@echo "  run           Build and run the CLI"
	@echo "  run-dev       Run the CLI in development mode"
	@echo "  ci            Run all CI checks"