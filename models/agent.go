package models

import "go-ai-terminal-assistant/agents"

// AgentPriority defines the priority order for agent routing
type AgentPriority int

const (
	PriorityHigh    AgentPriority = 1
	PriorityMedium  AgentPriority = 5
	PriorityLow     AgentPriority = 10
	PriorityDefault AgentPriority = 100 // Default agent should have lowest priority
)

// AgentRegistration holds an agent and its configuration
type AgentRegistration struct {
	Agent    agents.Agent
	Priority AgentPriority
	Enabled  bool
	Tags     []string // Optional tags for categorization
}

// AgentConfig holds configuration for agent behavior
type AgentConfig struct {
	EnableExampleAgents bool
	EnableCodeReview    bool
	EnableDataAnalysis  bool
	WeatherAPIKey       string
	CustomAgentPriority map[string]AgentPriority
}
