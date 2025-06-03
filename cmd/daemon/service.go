package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	serviceName    = "com.allday.ai-terminal-assistant"
	launchAgentDir = "/Users/%s/Library/LaunchAgents"
	globalAgentDir = "/Library/LaunchAgents"
	plistFileName  = serviceName + ".plist"
)

func installService(port, logFile string) error {
	// Get current user
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("unable to determine current user")
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create LaunchAgent directory
	agentDir := fmt.Sprintf(launchAgentDir, user)
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgent directory: %w", err)
	}

	// Create plist content
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>-daemon</string>
        <string>-port</string>
        <string>%s</string>
        <string>-log</string>
        <string>%s</string>
    </array>
    <key>KeepAlive</key>
    <true/>
    <key>RunAtLoad</key>
    <true/>
    <key>StandardOutPath</key>
    <string>%s</string>
    <key>StandardErrorPath</key>
    <string>%s</string>
    <key>WorkingDirectory</key>
    <string>%s</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin</string>
    </dict>
</dict>
</plist>`, serviceName, execPath, port, logFile, logFile, logFile, filepath.Dir(execPath))

	// Write plist file
	plistPath := filepath.Join(agentDir, plistFileName)
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	fmt.Printf("📝 Created LaunchAgent: %s\n", plistPath)
	return nil
}

func uninstallService() error {
	// Stop service first if running
	stopService()

	// Get current user
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("unable to determine current user")
	}

	// Remove plist file
	agentDir := fmt.Sprintf(launchAgentDir, user)
	plistPath := filepath.Join(agentDir, plistFileName)

	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}

	// Unload from launchctl if loaded
	cmd := exec.Command("launchctl", "unload", plistPath)
	cmd.Run() // Ignore errors as service might not be loaded

	fmt.Printf("🗑️  Removed LaunchAgent: %s\n", plistPath)
	return nil
}

func startService() error {
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("unable to determine current user")
	}

	agentDir := fmt.Sprintf(launchAgentDir, user)
	plistPath := filepath.Join(agentDir, plistFileName)

	// Check if plist exists
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return fmt.Errorf("service not installed. Run with -install first")
	}

	// Load service with launchctl
	cmd := exec.Command("launchctl", "load", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start service: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func stopService() error {
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("unable to determine current user")
	}

	agentDir := fmt.Sprintf(launchAgentDir, user)
	plistPath := filepath.Join(agentDir, plistFileName)

	// Unload service with launchctl
	cmd := exec.Command("launchctl", "unload", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Don't return error if service is not loaded
		if !strings.Contains(string(output), "Could not find specified service") {
			return fmt.Errorf("failed to stop service: %w\nOutput: %s", err, string(output))
		}
	}

	return nil
}

func isServiceRunning() (bool, error) {
	cmd := exec.Command("launchctl", "list", serviceName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Service not found
		if strings.Contains(string(output), "Could not find service") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check service status: %w", err)
	}

	// If we get here, service is loaded and potentially running
	// Check if PID is present in output
	return strings.Contains(string(output), "PID"), nil
}
