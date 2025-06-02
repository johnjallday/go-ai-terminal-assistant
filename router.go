package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"go-ai-terminal-assistant/agents"
	defaultagent "go-ai-terminal-assistant/agents/default"
	mathagent "go-ai-terminal-assistant/agents/math"
	weatheragent "go-ai-terminal-assistant/agents/weather"
	"go-ai-terminal-assistant/models"
)

// AgentRouter manages and routes prompts to appropriate agents
type AgentRouter struct {
	registrations []models.AgentRegistration
}

// NewAgentRouter creates a new agent router with dynamic agent loading
func NewAgentRouter() *AgentRouter {
	router := &AgentRouter{
		registrations: make([]models.AgentRegistration, 0),
	}

	// Register default agents
	router.registerDefaultAgents()

	// Sort agents by priority
	router.sortAgentsByPriority()

	return router
}

// registerDefaultAgents registers the standard set of agents
func (r *AgentRouter) registerDefaultAgents() {
	// Register Math Agent
	r.RegisterAgent(mathagent.New(), models.AgentRegistration{
		Priority: models.PriorityHigh,
		Enabled:  true,
		Tags:     []string{"math", "calculation", "computation"},
	})

	// Register Weather Agent (check if API key is available for enhanced features)
	weatherPriority := models.PriorityMedium
	weatherTags := []string{"weather", "forecast", "climate"}

	// If weather API key is available, use enhanced agent, otherwise use basic agent
	if os.Getenv("WEATHER_API_KEY") != "" {
		r.RegisterAgent(weatheragent.NewEnhanced(), models.AgentRegistration{
			Priority: weatherPriority,
			Enabled:  true,
			Tags:     append(weatherTags, "enhanced", "realtime"),
		})
	} else {
		r.RegisterAgent(weatheragent.NewBasic(), models.AgentRegistration{
			Priority: weatherPriority,
			Enabled:  true,
			Tags:     append(weatherTags, "basic"),
		})
	}

	// Register Default Agent (always last)
	r.RegisterAgent(defaultagent.New(), models.AgentRegistration{
		Priority: models.PriorityDefault,
		Enabled:  true,
		Tags:     []string{"general", "default", "fallback"},
	})
}

// RegisterAgent dynamically registers an agent with the router
func (r *AgentRouter) RegisterAgent(agent agents.Agent, config models.AgentRegistration) {
	config.Agent = agent
	r.registrations = append(r.registrations, config)
}

// UnregisterAgent removes an agent by name
func (r *AgentRouter) UnregisterAgent(agentName string) bool {
	for i, reg := range r.registrations {
		if reg.Agent.GetName() == agentName {
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			return true
		}
	}
	return false
}

// EnableAgent enables or disables an agent by name
func (r *AgentRouter) EnableAgent(agentName string, enabled bool) bool {
	for i := range r.registrations {
		if r.registrations[i].Agent.GetName() == agentName {
			r.registrations[i].Enabled = enabled
			return true
		}
	}
	return false
}

// SoloAgent disables all agents except the specified one
func (r *AgentRouter) SoloAgent(agentName string) bool {
	found := false
	for i := range r.registrations {
		if r.registrations[i].Agent.GetName() == agentName {
			r.registrations[i].Enabled = true
			found = true
		} else {
			r.registrations[i].Enabled = false
		}
	}
	return found
}

// UnsoloAgents re-enables all agents
func (r *AgentRouter) UnsoloAgents() {
	for i := range r.registrations {
		r.registrations[i].Enabled = true
	}
}

// sortAgentsByPriority sorts agents by their priority (lower number = higher priority)
func (r *AgentRouter) sortAgentsByPriority() {
	sort.Slice(r.registrations, func(i, j int) bool {
		return r.registrations[i].Priority < r.registrations[j].Priority
	})
}

// GetAgentsByTag returns agents that have a specific tag
func (r *AgentRouter) GetAgentsByTag(tag string) []agents.Agent {
	var agents []agents.Agent
	tag = strings.ToLower(tag)

	for _, reg := range r.registrations {
		if !reg.Enabled {
			continue
		}
		for _, agentTag := range reg.Tags {
			if strings.ToLower(agentTag) == tag {
				agents = append(agents, reg.Agent)
				break
			}
		}
	}
	return agents
}

// RoutePrompt analyzes the prompt and returns the best agent to handle it
func (r *AgentRouter) RoutePrompt(prompt string) agents.Agent {
	// Check each enabled agent in priority order
	for _, reg := range r.registrations {
		if !reg.Enabled {
			continue
		}

		// Skip default agent for now (it should be last)
		if reg.Priority == models.PriorityDefault {
			continue
		}

		if reg.Agent.CanHandle(prompt) {
			fmt.Printf("🎯 Routing to %s Agent\n", reg.Agent.GetName())
			return reg.Agent
		}
	}

	// Return default agent if no specialized agent can handle it
	for _, reg := range r.registrations {
		if reg.Enabled && reg.Priority == models.PriorityDefault {
			return reg.Agent
		}
	}

	// Fallback - should never happen if default agent is properly registered
	return defaultagent.New()
}

// ListAgents returns information about all available agents
func (r *AgentRouter) ListAgents() []agents.Agent {
	var agents []agents.Agent
	for _, reg := range r.registrations {
		if reg.Enabled {
			agents = append(agents, reg.Agent)
		}
	}
	return agents
}

// ListAllAgents returns all agents (enabled and disabled)
func (r *AgentRouter) ListAllAgents() []models.AgentRegistration {
	return r.registrations
}

// GetAgentStatus returns detailed status information about agents
func (r *AgentRouter) GetAgentStatus() string {
	var status strings.Builder
	status.WriteString("🤖 Agent Status:\n")
	status.WriteString("─────────────────────\n")

	for _, reg := range r.registrations {
		statusIcon := "✅"
		if !reg.Enabled {
			statusIcon = "❌"
		}

		priorityStr := fmt.Sprintf("P%d", reg.Priority)
		tags := strings.Join(reg.Tags, ", ")

		status.WriteString(fmt.Sprintf("%s %s Agent (%s) - %s\n",
			statusIcon, reg.Agent.GetName(), priorityStr, reg.Agent.GetDescription()))
		if len(reg.Tags) > 0 {
			status.WriteString(fmt.Sprintf("   Tags: %s\n", tags))
		}
	}

	return status.String()
}
