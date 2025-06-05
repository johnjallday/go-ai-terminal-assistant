package reaper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go"
	"go-ai-terminal-assistant/agents"
)

// normalizeName returns a lowercase, alphanumeric-only string (spaces, underscores,
// and hyphens removed) for matching script aliases.
func normalizeName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

// ReaperAgent launches the Reaper application on macOS.
type ReaperAgent struct{}

// New creates a new ReaperAgent.
func New() *ReaperAgent {
	return &ReaperAgent{}
}

// CanHandle returns true if the prompt requests launching Reaper or managing custom scripts.
func (a *ReaperAgent) CanHandle(prompt string) bool {
	lower := strings.ToLower(prompt)
	if strings.Contains(lower, "launch reaper") ||
		strings.Contains(lower, ".lua") ||
		strings.HasPrefix(lower, "list scripts") {
		return true
	}
	// Match script alias without .lua extension (e.g. "auto color tracks").
	normPrompt := normalizeName(lower)
	for _, tool := range a.Tools() {
		base := strings.TrimSuffix(strings.ToLower(tool.Name), ".lua")
		normName := normalizeName(base)
		if normName == normPrompt || strings.HasPrefix(normName, normPrompt) {
			return true
		}
	}
	return false
}

// Handle executes commands to launch Reaper or run/list custom Lua scripts.
func (a *ReaperAgent) Handle(prompt string, client *openai.Client, model string) (string, error) {
	lower := strings.ToLower(strings.TrimSpace(prompt))

	// Determine script directory from env or default
	scriptDir := os.Getenv("REAPER_SCRIPT_DIR")
	if scriptDir == "" {
		// default to custom_scripts folder in the project root (cwd)
		if cwd, err := os.Getwd(); err == nil {
			scriptDir = filepath.Join(cwd, "agents", "reaper", "custom_scripts")
		} else {
			scriptDir = "custom_scripts"
		}
	}

	// List available scripts
	if strings.HasPrefix(lower, "list scripts") {
		entries, err := os.ReadDir(scriptDir)
		if err != nil {
			return "", fmt.Errorf("failed to list scripts in %s: %w", scriptDir, err)
		}
		var scripts []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".lua") {
				scripts = append(scripts, e.Name())
			}
		}
		if len(scripts) == 0 {
			return fmt.Sprintf("No scripts found in %s", scriptDir), nil
		}
		return "Available scripts:\n" + strings.Join(scripts, "\n"), nil
	}

	// Launch script by alias (without .lua extension)
	{
		normPrompt := normalizeName(lower)
		for _, tool := range a.Tools() {
			base := strings.TrimSuffix(strings.ToLower(tool.Name), ".lua")
			normName := normalizeName(base)
			if normName == normPrompt || strings.HasPrefix(normName, normPrompt) {
				scriptName := tool.Name
				fullPath := scriptName
				if !filepath.IsAbs(scriptName) {
					fullPath = filepath.Join(scriptDir, scriptName)
				}
				if _, err := os.Stat(fullPath); err != nil {
					return "", fmt.Errorf("script not found: %s", scriptName)
				}
				cmd := exec.Command("open", "-a", "Reaper", fullPath)
				if err := cmd.Run(); err != nil {
					return "", fmt.Errorf("failed to launch script %s: %w", scriptName, err)
				}
				return fmt.Sprintf("✅ Script %s launched successfully.", scriptName), nil
			}
		}
	}

	// Launch a specific custom script if named
	if strings.Contains(lower, ".lua") {
		parts := strings.Fields(prompt)
		scriptName := parts[len(parts)-1]
		fullPath := scriptName
		if !filepath.IsAbs(scriptName) {
			fullPath = filepath.Join(scriptDir, scriptName)
		}
		if _, err := os.Stat(fullPath); err != nil {
			return "", fmt.Errorf("script not found: %s", scriptName)
		}
		cmd := exec.Command("open", "-a", "Reaper", fullPath)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to launch script %s: %w", scriptName, err)
		}
		return fmt.Sprintf("✅ Script %s launched successfully.", scriptName), nil
	}

	// Fallback: launch Reaper without scripts
	cmd := exec.Command("open", "-a", "Reaper")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to launch Reaper: %w", err)
	}
	return "✅ Reaper launched successfully.", nil
}

// GetName returns the name of the agent.
func (a *ReaperAgent) GetName() string {
	return "Reaper"
}

// GetDescription returns a brief description of the agent.
func (a *ReaperAgent) GetDescription() string {
	return "Agent for launching Reaper and running custom Lua scripts on macOS"
}

// Tools returns the list of custom Lua scripts available to launch as tools.
func (a *ReaperAgent) Tools() []agents.Tool {
	// Determine script directory from REAPER_SCRIPT_DIR or default location
	scriptDir := os.Getenv("REAPER_SCRIPT_DIR")
	if scriptDir == "" {
		// default to custom_scripts folder under agents/reaper in the project root (cwd)
		if cwd, err := os.Getwd(); err == nil {
			scriptDir = filepath.Join(cwd, "agents", "reaper", "custom_scripts")
		} else {
			scriptDir = filepath.Join("agents", "reaper", "custom_scripts")
		}
	}

	entries, err := os.ReadDir(scriptDir)
	if err != nil {
		return nil
	}
	var tools []agents.Tool
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".lua") {
			tools = append(tools, agents.Tool{
				Name:        e.Name(),
				Description: fmt.Sprintf("Launch script %s", e.Name()),
			})
		}
	}
	return tools
}
