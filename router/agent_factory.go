package router

import (
	"os"
	"strconv"

	defaultagent "go-ai-terminal-assistant/agents/default"
	mathagent "go-ai-terminal-assistant/agents/math"
	reaper "go-ai-terminal-assistant/agents/reaper"
	toolbuilder "go-ai-terminal-assistant/agents/toolbuilder"
	weatheragent "go-ai-terminal-assistant/agents/weather"
	"go-ai-terminal-assistant/models"
)

// NewAgentFactory creates a new agent factory
func NewAgentFactory() *AgentFactory {
	return &AgentFactory{
		config: LoadAgentConfig(),
	}
}

// LoadAgentConfig loads agent configuration from environment variables
func LoadAgentConfig() *models.AgentConfig {
	config := &models.AgentConfig{
		EnableExampleAgents: getBoolEnv("ENABLE_EXAMPLE_AGENTS", false),
		EnableCodeReview:    getBoolEnv("ENABLE_CODE_REVIEW_AGENT", false),
		EnableDataAnalysis:  getBoolEnv("ENABLE_DATA_ANALYSIS_AGENT", false),
		WeatherAPIKey:       os.Getenv("WEATHER_API_KEY"),
		CustomAgentPriority: make(map[string]models.AgentPriority),
	}
	return config
}

// getBoolEnv gets a boolean environment variable with a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// AgentFactory creates and configures agents based on configuration
type AgentFactory struct {
	config *models.AgentConfig
}

// CreateAgentRouter creates a router with agents based on configuration
func (f *AgentFactory) CreateAgentRouter() *AgentRouter {
	router := &AgentRouter{
		registrations: make([]models.AgentRegistration, 0),
	}
	f.registerCoreAgents(router)
	router.sortAgentsByPriority()
	return router
}

// registerCoreAgents registers the core agents (Math, Weather, Default)
func (f *AgentFactory) registerCoreAgents(router *AgentRouter) {
	router.RegisterAgent(mathagent.New(), models.AgentRegistration{
		Priority: models.PriorityHigh,
		Enabled:  true,
		Tags:     []string{"math", "calculation", "computation"},
	})

	// Always use Enhanced Weather Agent - it will automatically prompt for API key if needed
	router.RegisterAgent(weatheragent.NewEnhanced(), models.AgentRegistration{
		Priority: models.PriorityMedium,
		Enabled:  true,
		Tags:     []string{"weather", "forecast", "climate", "enhanced", "realtime"},
	})

	// Register Reaper Agent for launching Reaper and running custom Lua scripts on macOS
	router.RegisterAgent(reaper.New(), models.AgentRegistration{
		Priority: models.PriorityLow,
		Enabled:  true,
		Tags:     []string{"reaper", "macos", "scripts"},
	})

	// Register Script Builder Agent for generating custom Reaper Lua scripts via LLM
	router.RegisterAgent(toolbuilder.New(), models.AgentRegistration{
		Priority: models.PriorityLow,
		Enabled:  true,
		Tags:     []string{"script", "builder", "reaper", "lua"},
	})

	router.RegisterAgent(defaultagent.New(), models.AgentRegistration{
		Priority: models.PriorityDefault,
		Enabled:  true,
		Tags:     []string{"general", "default", "fallback"},
	})
}

// GetConfig returns the current configuration
func (f *AgentFactory) GetConfig() *models.AgentConfig {
	return f.config
}
