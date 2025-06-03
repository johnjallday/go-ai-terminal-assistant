#!/bin/bash

# AI Terminal Assistant - Installation Script
# This script installs the AI Terminal Assistant daemon on macOS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="ai-terminal-assistant"
INSTALL_DIR="/usr/local/bin"
SERVICE_NAME="com.johnjallday.ai-terminal-assistant"

# Helper functions
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if running on macOS
check_macos() {
    if [[ "$OSTYPE" != "darwin"* ]]; then
        print_error "This installer is designed for macOS only."
        exit 1
    fi
}

# Check if binary exists in build directory or current directory
check_binary() {
    if [[ -f "./build/$BINARY_NAME" ]]; then
        BINARY_PATH="./build/$BINARY_NAME"
    elif [[ -f "./$BINARY_NAME" ]]; then
        BINARY_PATH="./$BINARY_NAME"
    else
        print_error "Binary '$BINARY_NAME' not found in ./build/ or current directory."
        print_info "Please build the daemon first with: make build-daemon"
        exit 1
    fi
}

# Check if already installed
check_existing() {
    if [[ -f "$INSTALL_DIR/$BINARY_NAME" ]]; then
        print_warning "AI Terminal Assistant is already installed."
        read -p "Do you want to update it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled."
            exit 0
        fi
        
        # Stop service if running
        print_info "Stopping existing service..."
        if command -v "$BINARY_NAME" >/dev/null 2>&1; then
            "$BINARY_NAME" -stop 2>/dev/null || true
        fi
    fi
}

# Check for required environment
check_requirements() {
    print_info "Checking requirements..."
    
    # Check for sudo access
    if ! sudo -n true 2>/dev/null; then
        print_info "This installer requires sudo access to install to $INSTALL_DIR"
        print_info "You may be prompted for your password."
    fi
    
    # Check if OpenAI API key is set
    if [[ -z "$OPENAI_API_KEY" ]]; then
        print_warning "OPENAI_API_KEY environment variable is not set."
        print_info "You'll need to configure this before using the service."
        print_info "Add this to your ~/.zshrc or ~/.bash_profile:"
        echo "export OPENAI_API_KEY='your-api-key-here'"
    fi
}

# Install binary
install_binary() {
    print_info "Installing AI Terminal Assistant daemon..."
    
    # Copy binary to install directory
    sudo cp "$BINARY_PATH" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    print_success "Binary installed to $INSTALL_DIR/$BINARY_NAME"
}

# Install service
install_service() {
    print_info "Installing system service..."
    
    # Install as LaunchAgent service
    if "$INSTALL_DIR/$BINARY_NAME" -install; then
        print_success "Service installed successfully!"
    else
        print_error "Failed to install service."
        exit 1
    fi
}

# Create .env file template
create_env_template() {
    local env_file="$HOME/.ai-terminal-assistant.env"
    
    if [[ ! -f "$env_file" ]]; then
        print_info "Creating environment configuration template..."
        
        cat > "$env_file" << EOF
# AI Terminal Assistant Environment Configuration
# Copy this file and set your actual values

# Required: OpenAI API Key
OPENAI_API_KEY=your-openai-api-key-here

# Optional: Custom port (default: 8080)
# PORT=8080

# Optional: Custom log file location
# LOG_FILE=/usr/local/var/log/ai-terminal-assistant.log
EOF
        
        print_success "Environment template created at $env_file"
        print_warning "Please edit this file and set your actual OpenAI API key!"
    fi
}

# Main installation function
main() {
    echo "🤖 AI Terminal Assistant - macOS Installation"
    echo "=============================================="
    echo
    
    check_macos
    check_binary
    check_existing
    check_requirements
    
    echo
    print_info "Starting installation..."
    
    install_binary
    install_service
    create_env_template
    
    echo
    print_success "Installation complete! 🎉"
    echo
    print_info "Next steps:"
    echo "1. Edit ~/.ai-terminal-assistant.env and set your OpenAI API key"
    echo "2. Start the service: $BINARY_NAME -start"
    echo "3. Check status: $BINARY_NAME -status"
    echo "4. Test the API: curl http://localhost:8080/health"
    echo
    print_info "Service management commands:"
    echo "  $BINARY_NAME -start    # Start the service"
    echo "  $BINARY_NAME -stop     # Stop the service"
    echo "  $BINARY_NAME -status   # Check service status"
    echo "  $BINARY_NAME -uninstall # Remove the service"
    echo
    print_info "For more information, visit: https://github.com/johnjallday/go-ai-terminal-assistant"
}

# Run main function
main "$@"
