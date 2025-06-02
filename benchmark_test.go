package main

import (
	"go-ai-terminal-assistant/models"
	"go-ai-terminal-assistant/storage"
	"testing"
	"time"
)

// BenchmarkGetModelDisplayName benchmarks the model display name lookup
func BenchmarkGetModelDisplayName(b *testing.B) {
	modelName := "gpt-4.1-mini"

	for i := 0; i < b.N; i++ {
		models.GetModelDisplayName(modelName)
	}
}

// BenchmarkModelLookup benchmarks looking up models by different patterns
func BenchmarkModelLookup(b *testing.B) {
	testCases := []string{
		"gpt-4.1",
		"gpt-4.1-mini",
		"gpt-4.1-nano",
		"nonexistent-model",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, modelName := range testCases {
			models.GetModelDisplayName(modelName)
		}
	}
}

// BenchmarkConversationFileCreation benchmarks creating conversation file structs
func BenchmarkConversationFileCreation(b *testing.B) {
	now := time.Now()

	for i := 0; i < b.N; i++ {
		cf := storage.ConversationFile{
			Filename:    "test_conversation.txt",
			Timestamp:   now,
			DisplayName: now.Format("2006-01-02 15:04:05"),
		}
		_ = cf // Use the variable to prevent optimization
	}
}

// BenchmarkAgentConfigCreation benchmarks creating agent configuration
func BenchmarkAgentConfigCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config := &models.AgentConfig{
			EnableExampleAgents: true,
			EnableCodeReview:    false,
			EnableDataAnalysis:  true,
			WeatherAPIKey:       "test-key",
			CustomAgentPriority: make(map[string]models.AgentPriority),
		}
		_ = config // Use the variable to prevent optimization
	}
}

// BenchmarkAgentRegistrationCreation benchmarks creating agent registrations
func BenchmarkAgentRegistrationCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		registration := &models.AgentRegistration{
			Agent:    nil,
			Priority: models.PriorityMedium,
			Enabled:  true,
			Tags:     []string{"test", "example"},
		}
		_ = registration // Use the variable to prevent optimization
	}
}

// BenchmarkAvailableModelsIteration benchmarks iterating through all available models
func BenchmarkAvailableModelsIteration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, model := range models.AvailableModels {
			_ = model.Name
			_ = model.DisplayName
			_ = model.Description
		}
	}
}
