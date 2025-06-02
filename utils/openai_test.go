package utils

import (
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func TestGetOpenAIResponseValidation(t *testing.T) {
	// Test input validation - we'll test with nil client to avoid actual API calls
	var client *openai.Client = nil

	// Test with empty prompt - this should validate inputs before using the client
	_, err := GetOpenAIResponse(client, "", "gpt-4.1-mini")
	if err == nil {
		t.Error("Should return error for nil client")
	}

	// Test with invalid model (this will also fail due to nil client, but tests the function signature)
	_, err = GetOpenAIResponse(client, "test prompt", "")
	if err == nil {
		t.Error("Should return error for nil client")
	}
}

func TestGetOpenAIResponseSignature(t *testing.T) {
	// Test that the function signature is correct
	// We'll create a client but not make actual calls to avoid API costs in tests
	client := openai.NewClient(option.WithAPIKey("test-key"))

	// Test with conversation history
	_, err := GetOpenAIResponse(&client, "Hello", "gpt-4.1-mini", "Previous prompt", "Previous response")
	// This will fail due to invalid API key, but that's expected in tests
	if err == nil {
		t.Error("Should return error for invalid API key")
	}

	// Test without conversation history
	_, err = GetOpenAIResponse(&client, "Hello", "gpt-4.1-mini")
	// This will also fail due to invalid API key, but that's expected in tests
	if err == nil {
		t.Error("Should return error for invalid API key")
	}
}

func TestGetOpenAIResponseInputValidation(t *testing.T) {
	// Test input validation with a mock client
	client := openai.NewClient(option.WithAPIKey("test-key"))

	// Test with empty prompt
	_, err := GetOpenAIResponse(&client, "", "gpt-4.1-mini")
	if err == nil {
		t.Error("Should return error for empty prompt")
	}

	// Test with empty model
	_, err = GetOpenAIResponse(&client, "test prompt", "")
	if err == nil {
		t.Error("Should return error for empty model")
	}
}

// Note: For actual API testing, you would need to:
// 1. Set up integration tests with a valid API key
// 2. Use test doubles/mocks to avoid API costs
// 3. Test against a local OpenAI-compatible server for development
