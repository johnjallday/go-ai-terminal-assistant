package toolbuilder

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"go-ai-terminal-assistant/agents"
)

// ToolBuilderAgent generates Reaper Lua scripts based on user requirements.
type ToolBuilderAgent struct{}

// New creates a new ToolBuilderAgent.
func New() *ToolBuilderAgent {
	return &ToolBuilderAgent{}
}

// CanHandle returns true if the prompt requests generation of a Reaper Lua script.
func (a *ToolBuilderAgent) CanHandle(prompt string) bool {
	lower := strings.ToLower(prompt)
	return strings.HasPrefix(lower, "build script") ||
		strings.HasPrefix(lower, "generate script") ||
		(strings.Contains(lower, "reaper") && strings.Contains(lower, "script"))
}

// Handle uses the LLM to generate a complete Reaper Lua script.
func (a *ToolBuilderAgent) Handle(prompt string, client *openai.Client, model string) (string, error) {
	sysCtx := `You are a Reaper reascript generator. When the user asks you to build or scaffold a script, generate a complete, runnable Reaper Lua script. Include a header comment with the tool name and description, ensure there is a Main() function, and wrap the entire script in a Lua code block without additional explanation.`
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(sysCtx),
		openai.UserMessage(prompt),
	}
	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate script: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from script generator")
	}
	return resp.Choices[0].Message.Content, nil
}

// GetName returns the name of the agent.
func (a *ToolBuilderAgent) GetName() string {
	return "Script Builder"
}

// GetDescription returns a brief description of the agent.
func (a *ToolBuilderAgent) GetDescription() string {
	return "Generate Reaper Lua scripts based on user requirements"
}

// Tools returns no CLI tools for this agent.
func (a *ToolBuilderAgent) Tools() []agents.Tool {
	return nil
}
