package utils

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
)

// GetOpenAIResponse handles communication with OpenAI API
func GetOpenAIResponse(client *openai.Client, prompt string, model string, conversationHistory ...string) (string, error) {
	// Validate inputs
	if client == nil {
		return "", fmt.Errorf("client cannot be nil")
	}
	if prompt == "" {
		return "", fmt.Errorf("prompt cannot be empty")
	}
	if model == "" {
		return "", fmt.Errorf("model cannot be empty")
	}

	ctx := context.Background()

	// Build messages array
	var messages []openai.ChatCompletionMessageParamUnion

	// If we have conversation history, add it as context
	if len(conversationHistory) >= 2 {
		previousPrompt := conversationHistory[0]
		previousResponse := conversationHistory[1]

		messages = append(messages, openai.UserMessage(previousPrompt))
		messages = append(messages, openai.AssistantMessage(previousResponse))
	}

	// Add the current prompt
	messages = append(messages, openai.UserMessage(prompt))

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    model,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get completion: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return completion.Choices[0].Message.Content, nil
}
