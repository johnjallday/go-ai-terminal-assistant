# Makefile for go-ai-terminal-assistant

.PHONY: build test test-unit test-integration test-coverage clean help

# Build the application
build:
	@echo "🔨 Building application..."
	@go build .

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
	@rm -f coverage.out coverage.html
	@echo "Cleanup complete"

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

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
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
	@echo "  help           - Show this help message"
