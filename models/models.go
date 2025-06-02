package models

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ModelInfo represents information about an available model
type ModelInfo struct {
	Name        string
	DisplayName string
	Description string
}

// Available models
var AvailableModels = []ModelInfo{
	{
		Name:        "gpt-4.1",
		DisplayName: "GPT-4.1",
		Description: "Latest flagship model - $2.00 input / $8.00 output per 1M tokens",
	},
	{
		Name:        "gpt-4.1-mini",
		DisplayName: "GPT-4.1 Mini",
		Description: "Balanced performance and cost - $0.40 input / $1.60 output per 1M tokens",
	},
	{
		Name:        "gpt-4.1-nano",
		DisplayName: "GPT-4.1 Nano",
		Description: "Ultra-affordable option - $0.10 input / $0.40 output per 1M tokens",
	},
	{
		Name:        "gpt-4.5-preview",
		DisplayName: "GPT-4.5 Preview",
		Description: "Most advanced preview model - $75.00 input / $150.00 output per 1M tokens",
	},
	{
		Name:        "o4-mini",
		DisplayName: "O4 Mini",
		Description: "Efficient reasoning model - $1.10 input / $4.40 output per 1M tokens",
	},
	{
		Name:        "o3-mini",
		DisplayName: "O3 Mini",
		Description: "Advanced reasoning model - $1.10 input / $4.40 output per 1M tokens",
	},
}

// GetDefaultModel returns the default model (GPT-4.1 Nano)
func GetDefaultModel() string {
	return "gpt-4.1-nano"
}

// SelectModel displays available models and lets user choose one
func SelectModel() string {
	fmt.Println("\n🎯 Available Models:")
	fmt.Println("─────────────────────")

	for i, model := range AvailableModels {
		fmt.Printf("%d. %s - %s\n", i+1, model.DisplayName, model.Description)
	}

	fmt.Print("\nSelect a model (1-6) [default: 2]: ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		choice := strings.TrimSpace(scanner.Text())

		// Default to GPT-4.1 Mini if empty
		if choice == "" {
			return AvailableModels[1].Name
		}

		// Parse user choice
		if num, err := strconv.Atoi(choice); err == nil && num >= 1 && num <= len(AvailableModels) {
			return AvailableModels[num-1].Name
		}
	}

	// Default fallback
	return AvailableModels[1].Name
}

// GetModelDisplayName returns the display name for a model
func GetModelDisplayName(modelName string) string {
	for _, model := range AvailableModels {
		if model.Name == modelName {
			return model.DisplayName
		}
	}
	return modelName
}
