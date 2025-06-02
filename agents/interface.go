package agents

import "github.com/openai/openai-go"

// Agent interface defines the contract for all specialized agents
type Agent interface {
	CanHandle(prompt string) bool
	Handle(prompt string, client *openai.Client, model string) (string, error)
	GetName() string
	GetDescription() string
}

// AgentResult holds the result of agent processing
type AgentResult struct {
	AgentName string
	Response  string
	Error     error
}
