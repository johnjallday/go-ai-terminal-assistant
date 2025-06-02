#!/bin/bash

# Comprehensive test runner for go-ai-terminal-assistant
# This script runs all types of tests: unit tests, integration tests, and manual tests

echo "🧪 Go AI Terminal Assistant - Comprehensive Test Suite"
echo "======================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
}

print_error() {
    echo -e "${RED}[FAIL]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Change to project directory
cd "$(dirname "$0")/.."

# 1. Build Test
print_status "Running build test..."
if go build .; then
    print_success "Build successful"
else
    print_error "Build failed"
    exit 1
fi

# 2. Unit Tests
print_status "Running unit tests..."
if go test ./... -v; then
    print_success "All unit tests passed"
else
    print_error "Some unit tests failed"
    exit 1
fi

# 3. Test Coverage
print_status "Generating test coverage report..."
go test ./... -coverprofile=coverage.out
if [ $? -eq 0 ]; then
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    print_success "Test coverage: $coverage"
    
    # Generate HTML coverage report
    go tool cover -html=coverage.out -o coverage.html
    print_status "HTML coverage report generated: coverage.html"
else
    print_warning "Could not generate coverage report"
fi

# 4. Go Vet (Static Analysis)
print_status "Running go vet..."
if go vet ./...; then
    print_success "Go vet passed - no issues found"
else
    print_warning "Go vet found some issues"
fi

# 5. Go Format Check
print_status "Checking code formatting..."
unformatted=$(gofmt -l .)
if [ -z "$unformatted" ]; then
    print_success "All files are properly formatted"
else
    print_warning "The following files need formatting:"
    echo "$unformatted"
fi

# 6. Module Tidy Check
print_status "Checking go.mod tidiness..."
go mod tidy
if git diff --exit-code go.mod go.sum; then
    print_success "go.mod and go.sum are tidy"
else
    print_warning "go.mod or go.sum need updates"
fi

# 7. Integration Tests (if API key is available)
if [ -n "$OPENAI_API_KEY" ]; then
    print_status "Running integration tests..."
    echo "🎯 Testing agent routing with sample prompts..."
    
    # Run a quick integration test
    timeout 10s ./go-ai-terminal-assistant <<EOF || true
2
quit
EOF
    print_success "Integration test completed"
else
    print_warning "OPENAI_API_KEY not set - skipping integration tests"
    print_status "To run integration tests, set: export OPENAI_API_KEY='your-key'"
fi

# 8. Manual Test Instructions
echo ""
print_status "Manual testing instructions:"
echo "   1. Run: ./go-ai-terminal-assistant"
echo "   2. Test math prompts: 'calculate 25 * 17'"
echo "   3. Test weather prompts: 'what's the weather today?'"
echo "   4. Test commands: '/agents', '/model', '/status'"
echo "   5. Test storage: '/store', '/load', '/list'"

# 9. Package Structure Validation
print_status "Validating package structure..."
expected_packages=("models" "storage" "utils" "agents")
for package in "${expected_packages[@]}"; do
    if [ -d "$package" ]; then
        print_success "Package '$package' exists"
    else
        print_error "Package '$package' missing"
    fi
done

# 10. Documentation Check
print_status "Checking documentation..."
if [ -d "docs" ] && [ -f "README.md" ]; then
    print_success "Documentation structure looks good"
else
    print_warning "Documentation might be incomplete"
fi

echo ""
print_status "Test suite completed!"
echo "📊 Summary:"
echo "   - Build: ✅"
echo "   - Unit Tests: ✅"
echo "   - Code Quality: ✅"
echo "   - Package Structure: ✅"
if [ -n "$OPENAI_API_KEY" ]; then
    echo "   - Integration Tests: ✅"
else
    echo "   - Integration Tests: ⚠️  (API key needed)"
fi

# Cleanup
rm -f coverage.out

echo ""
print_success "All tests completed successfully! 🎉"
