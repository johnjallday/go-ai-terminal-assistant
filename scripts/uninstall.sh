#!/bin/bash

# AI Terminal Assistant - Uninstallation Script
# This script removes the AI Terminal Assistant daemon from macOS

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
        print_error "This uninstaller is designed for macOS only."
        exit 1
    fi
}

# Check if installed
check_installed() {
    if [[ ! -f "$INSTALL_DIR/$BINARY_NAME" ]]; then
        print_warning "AI Terminal Assistant is not installed."
        exit 0
    fi
}

# Confirm uninstallation
confirm_uninstall() {
    print_warning "This will completely remove AI Terminal Assistant from your system."
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Uninstallation cancelled."
        exit 0
    fi
}

# Stop and uninstall service
uninstall_service() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        print_info "Stopping service..."
        "$BINARY_NAME" -stop 2>/dev/null || true
        
        print_info "Uninstalling service..."
        if "$BINARY_NAME" -uninstall; then
            print_success "Service uninstalled successfully!"
        else
            print_warning "Service uninstall failed, but continuing..."
        fi
    fi
}

# Remove binary
remove_binary() {
    print_info "Removing binary..."
    sudo rm -f "$INSTALL_DIR/$BINARY_NAME"
    print_success "Binary removed from $INSTALL_DIR/$BINARY_NAME"
}

# Remove configuration files (optional)
remove_config() {
    local env_file="$HOME/.ai-terminal-assistant.env"
    local log_file="/usr/local/var/log/ai-terminal-assistant.log"
    local pid_file="/usr/local/var/run/ai-terminal-assistant.pid"
    
    print_info "Checking for configuration files..."
    
    if [[ -f "$env_file" ]]; then
        read -p "Remove environment configuration file? ($env_file) (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -f "$env_file"
            print_success "Environment file removed"
        fi
    fi
    
    if [[ -f "$log_file" ]]; then
        read -p "Remove log file? ($log_file) (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            sudo rm -f "$log_file"
            print_success "Log file removed"
        fi
    fi
    
    # Always remove PID file
    sudo rm -f "$pid_file" 2>/dev/null || true
}

# Clean up LaunchAgent files manually (in case service uninstall failed)
cleanup_launchagent() {
    local plist_file="$HOME/Library/LaunchAgents/$SERVICE_NAME.plist"
    
    if [[ -f "$plist_file" ]]; then
        print_info "Cleaning up LaunchAgent configuration..."
        
        # Unload the service
        launchctl unload "$plist_file" 2>/dev/null || true
        
        # Remove the plist file
        rm -f "$plist_file"
        print_success "LaunchAgent configuration removed"
    fi
}

# Main uninstallation function
main() {
    echo "🤖 AI Terminal Assistant - macOS Uninstallation"
    echo "================================================"
    echo
    
    check_macos
    check_installed
    confirm_uninstall
    
    echo
    print_info "Starting uninstallation..."
    
    uninstall_service
    cleanup_launchagent
    remove_binary
    remove_config
    
    echo
    print_success "Uninstallation complete! 🎉"
    echo
    print_info "AI Terminal Assistant has been removed from your system."
    print_info "Thank you for using AI Terminal Assistant!"
    echo
    print_info "If you encountered any issues, please report them at:"
    print_info "https://github.com/johnjallday/go-ai-terminal-assistant/issues"
}

# Run main function
main "$@"