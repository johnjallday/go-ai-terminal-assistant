package defaultagent

import (
	"go-ai-terminal-assistant/utils"

	"github.com/openai/openai-go"
)

// DefaultAgent handles general conversations using OpenAI
type DefaultAgent struct{}

func New() *DefaultAgent {
	return &DefaultAgent{}
}

func (a *DefaultAgent) CanHandle(prompt string) bool {
	return true // Default agent can handle anything
}

func (a *DefaultAgent) Handle(prompt string, client *openai.Client, model string) (string, error) {
	return utils.GetOpenAIResponse(client, prompt, model)
}

func (a *DefaultAgent) GetName() string {
	return "Default"
}

func (a *DefaultAgent) GetDescription() string {
	return "General conversation agent using OpenAI models"
}
