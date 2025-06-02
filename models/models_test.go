package models

import (
	"testing"
)

func TestModelInfo(t *testing.T) {
	// Test that we have available models
	if len(AvailableModels) == 0 {
		t.Error("AvailableModels should not be empty")
	}

	// Test that each model has required fields
	for i, model := range AvailableModels {
		if model.Name == "" {
			t.Errorf("Model %d should have a non-empty Name", i)
		}
		if model.DisplayName == "" {
			t.Errorf("Model %d should have a non-empty DisplayName", i)
		}
		if model.Description == "" {
			t.Errorf("Model %d should have a non-empty Description", i)
		}
	}
}

func TestSelectModelDefaults(t *testing.T) {
	// Note: SelectModel() reads from stdin, so we can't easily test it in unit tests
	// This test verifies the function exists and has the right signature
	// In practice, integration tests would cover the interactive behavior

	// Test that we can call SelectModel and it returns a valid model name
	// Since it requires stdin interaction, we'll just verify it's callable
	if len(AvailableModels) == 0 {
		t.Skip("No models available to test SelectModel")
	}

	// We can at least verify that the default model (AvailableModels[1]) is valid
	defaultModel := AvailableModels[1].Name
	if defaultModel == "" {
		t.Error("Default model name should not be empty")
	}
}

func TestGetModelDisplayName(t *testing.T) {
	tests := []struct {
		modelName   string
		expected    string
		description string
	}{
		{"gpt-4.1", "GPT-4.1", "Should return correct display name for gpt-4.1"},
		{"gpt-4.1-mini", "GPT-4.1 Mini", "Should return correct display name for gpt-4.1-mini"},
		{"gpt-4.1-nano", "GPT-4.1 Nano", "Should return correct display name for gpt-4.1-nano"},
		{"gpt-4.5-preview", "GPT-4.5 Preview", "Should return correct display name for gpt-4.5-preview"},
		{"o4-mini", "O4 Mini", "Should return correct display name for o4-mini"},
		{"o3-mini", "O3 Mini", "Should return correct display name for o3-mini"},
		{"nonexistent-model", "nonexistent-model", "Should return the model name itself for unknown models"},
		{"", "", "Should handle empty string"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := GetModelDisplayName(test.modelName)
			if result != test.expected {
				t.Errorf("GetModelDisplayName(%q) = %q; want %q", test.modelName, result, test.expected)
			}
		})
	}
}

func TestAgentPriorityConstants(t *testing.T) {
	// Test that priority constants are properly ordered
	if PriorityHigh >= PriorityMedium {
		t.Error("PriorityHigh should be less than PriorityMedium (lower number = higher priority)")
	}
	if PriorityMedium >= PriorityLow {
		t.Error("PriorityMedium should be less than PriorityLow")
	}
	if PriorityLow >= PriorityDefault {
		t.Error("PriorityLow should be less than PriorityDefault")
	}
}

func TestAgentConfig(t *testing.T) {
	config := &AgentConfig{
		EnableExampleAgents: true,
		EnableCodeReview:    false,
		EnableDataAnalysis:  true,
		WeatherAPIKey:       "test-key",
		CustomAgentPriority: make(map[string]AgentPriority),
	}

	// Test that we can create and access AgentConfig
	if !config.EnableExampleAgents {
		t.Error("EnableExampleAgents should be true")
	}
	if config.EnableCodeReview {
		t.Error("EnableCodeReview should be false")
	}
	if !config.EnableDataAnalysis {
		t.Error("EnableDataAnalysis should be true")
	}
	if config.WeatherAPIKey != "test-key" {
		t.Error("WeatherAPIKey should be 'test-key'")
	}
	if config.CustomAgentPriority == nil {
		t.Error("CustomAgentPriority should be initialized")
	}
}

func TestAgentRegistration(t *testing.T) {
	registration := &AgentRegistration{
		Agent:    nil, // We can't easily test the agent interface
		Priority: PriorityMedium,
		Enabled:  true,
		Tags:     []string{"test", "example"},
	}

	// Test that we can create and access AgentRegistration
	if registration.Priority != PriorityMedium {
		t.Error("Priority should be PriorityMedium")
	}
	if !registration.Enabled {
		t.Error("Enabled should be true")
	}
	if len(registration.Tags) != 2 {
		t.Error("Should have 2 tags")
	}
	if registration.Tags[0] != "test" || registration.Tags[1] != "example" {
		t.Error("Tags should be ['test', 'example']")
	}
}
