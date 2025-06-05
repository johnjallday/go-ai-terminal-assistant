# Makefile for go-ai-terminal-assistant

# Binary name and version
BINARY_NAME=ai-terminal-assistant
VERSION?=1.0.0
BUILD_DIR=build

# Go build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"

.PHONY: build build-daemon build-terminal build-all install uninstall package test test-unit test-integration test-coverage clean help

# Build the terminal application
build:
	@echo "🔨 Building terminal application..."
	@go build .

# Build daemon binary
build-daemon: $(BUILD_DIR)
	@echo "🔨 Building AI Terminal Assistant daemon..."
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/daemon
	@echo "✅ Daemon build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build terminal version with custom name
build-terminal: $(BUILD_DIR)
	@echo "🔨 Building AI Terminal Assistant (terminal version)..."
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-terminal .
	@echo "✅ Terminal build complete: $(BUILD_DIR)/$(BINARY_NAME)-terminal"

# Build both versions
build-all: build-daemon build-terminal

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Install daemon binary to /usr/local/bin
install: build-daemon
	@echo "📦 Installing AI Terminal Assistant daemon..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Installation complete!"
	@echo "💡 Run '$(BINARY_NAME) -install' to set up as a service"

# Uninstall daemon binary
uninstall:
	@echo "🗑️  Uninstalling AI Terminal Assistant daemon..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstallation complete!"
	@echo "💡 Run '$(BINARY_NAME) -uninstall' first if you have the service installed"

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	@echo "🧪 Running unit tests..."
	@go test ./... -v

# Run integration tests (requires OPENAI_API_KEY)
test-integration:
	@echo "🔗 Running integration tests..."
	@./tests/test_agent_routing.sh

# Run comprehensive test suite
test-all:
	@echo "🚀 Running comprehensive test suite..."
	@./tests/run_all_tests.sh

# Generate test coverage
test-coverage:
	@echo "📊 Generating test coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run static analysis
lint:
	@echo "🔍 Running static analysis..."
	@go vet ./...
	@gofmt -l .

# Clean build artifacts and test files
clean:
	@echo "🧹 Cleaning up..."
	@rm -f allday-term-agent
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Cleanup complete"

# Create release package
package: build-daemon
	@echo "📦 Creating release package..."
	@mkdir -p $(BUILD_DIR)/release
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/release/
	@cp scripts/install.sh $(BUILD_DIR)/release/ 2>/dev/null || echo "⚠️  install.sh not found, skipping"
	@cp scripts/uninstall.sh $(BUILD_DIR)/release/ 2>/dev/null || echo "⚠️  uninstall.sh not found, skipping"
	@cp README.md $(BUILD_DIR)/release/ 2>/dev/null || echo "⚠️  README.md not found, skipping"
	@cp docs/DAEMON.md $(BUILD_DIR)/release/ 2>/dev/null || echo "⚠️  DAEMON.md not found, skipping"
	@cd $(BUILD_DIR) && tar -czf ai-terminal-assistant-$(VERSION)-macos.tar.gz release/
	@echo "✅ Release package created: $(BUILD_DIR)/ai-terminal-assistant-$(VERSION)-macos.tar.gz"

# Run daemon locally (for development)
run-daemon: build-daemon
	@echo "🚀 Running AI Terminal Assistant daemon locally..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -port 8080

# Run terminal version locally (for development)
run-terminal: build-terminal
	@echo "🚀 Running AI Terminal Assistant (terminal version)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)-terminal

# Format code
fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...

# Tidy dependencies
tidy:
	@echo "📦 Tidying dependencies..."
	@go mod tidy

# Development setup
dev-setup: tidy build test-unit
	@echo "🔧 Development environment ready!"

# Quick development check
dev-check: fmt lint test-unit
	@echo "✅ Development check complete!"

# Create a new Reaper Lua script template in agents/reaper/custom_scripts
# Usage: make reaper-script NAME=script_name
reaper-script:
	@if [ -z "${NAME}" ]; then \
		echo "Usage: make reaper-script NAME=script_name"; exit 1; \
	fi
	@dst=agents/reaper/custom_scripts/${NAME}.lua; \
	if [ -e "$$dst" ]; then \
		echo "Error: $$dst already exists"; exit 1; \
	fi; \
	mkdir -p $$(dirname "$$dst"); \
	printf "--[[\nTool: %s\nDescription: TODO: add description\n--]]\n\n-- Main function\nfunction Main()\n    -- TODO: implement script logic\nend\n\nMain()\n" "${NAME}" > "$$dst"; \
	echo "Created new Reaper script template: $$dst"

# Help
help:
	@echo "AI Terminal Assistant Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build          - Build the terminal application"
	@echo "  build-daemon   - Build daemon binary"
	@echo "  build-terminal - Build terminal version with custom name"
	@echo "  build-all      - Build both daemon and terminal versions"
	@echo "  install        - Install daemon to /usr/local/bin"
	@echo "  uninstall      - Remove daemon from /usr/local/bin"
	@echo "  run-daemon     - Run daemon locally"
	@echo "  run-terminal   - Run terminal version locally"
	@echo "  package        - Create release package"
	@echo "  test           - Run all tests (unit + integration)"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests"
	@echo "  test-all       - Run comprehensive test suite"
	@echo "  test-coverage  - Generate test coverage report"
	@echo "  lint           - Run static analysis"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy dependencies"
	@echo "  dev-setup      - Set up development environment"
	@echo "  dev-check      - Quick development check"
	@echo "  reaper-script  - Scaffold a new Reaper Lua script template (NAME=<script_name>)"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  VERSION        - Set version (default: 1.0.0)"
