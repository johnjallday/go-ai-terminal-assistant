# 🤖 Agent System Documentation

## Overview

The allday-term-agent features an intelligent agentic routing system that automatically selects specialized agents based on user prompts. This document provides comprehensive information about the agent system architecture, available agents, and how to use them.

## 🏗️ Agent Architecture

The system uses a modular package structure with priority-based routing:

```
agents/
├── interface.go           # Shared Agent interface & types
├── default/
│   └── default.go        # DefaultAgent implementation
├── math/
│   └── math.go          # MathAgent implementation
├── weather/
│   └── weather.go       # WeatherAgent & EnhancedWeatherAgent
└── examples/
    └── examples.go      # CodeReviewAgent & DataAnalysisAgent
```

## 🎯 Available Agents

### 🧮 Math Agent
**Priority:** 1 (Highest)  
**Tags:** `["math", "calculation", "computation"]`  
**Purpose:** Specialized mathematical calculations, equations, and problem solving

**Triggers on:** Mathematical keywords, symbols, and calculations
- Keywords: calculate, compute, solve, math, equation, formula, add, subtract, multiply, divide, percentage, etc.
- Symbols: +, -, *, /, =, ^, √, %, π, etc.
- Patterns: Number expressions like "25 * 17" or "15% of 200"

**Example prompts:**
```
Calculate 25 * 17
What is 15% of 200?
Solve for x: 2x + 5 = 15
What's the square root of 144?
Convert 68°F to Celsius
What's the area of a circle with radius 5?
Calculate compound interest on $1000 at 5% for 3 years
```

### 🌤️ Weather Agents
**Priority:** 5 (Medium)  
**Auto-Selection:** System automatically chooses Enhanced or Basic based on API key availability

#### Enhanced Weather Agent
**Tags:** `["weather", "forecast", "climate", "enhanced", "realtime"]`  
**Requirements:** `WEATHER_API_KEY` environment variable  
**Purpose:** Real-time weather data with AI fallback using WeatherAPI.com

#### Basic Weather Agent  
**Tags:** `["weather", "forecast", "climate", "basic"]`  
**Requirements:** None (fallback when no API key)  
**Purpose:** AI-powered weather responses without real-time data

**Triggers on:** Weather-related keywords
- Keywords: weather, temperature, rain, snow, sunny, cloudy, forecast, hot, cold, humid, wind, storm, etc.

**Example prompts:**
```
What's the weather in London?
Is it raining in New York?
What's the temperature in Tokyo today?
Will it be sunny tomorrow in Paris?
How humid is it in Miami?
What's the forecast for San Francisco?
Is there a storm coming to Chicago?
```

### 🔍 Example Agents (Optional)
**Priority:** 5 (Medium)  
**Requirements:** Environment variable configuration

#### Code Review Agent
**Tags:** `["code", "review", "programming", "optimization", "debug"]`  
**Enable:** Set `ENABLE_CODE_REVIEW_AGENT=true` or `ENABLE_EXAMPLE_AGENTS=true`  
**Purpose:** Expert in code review, optimization, and debugging

#### Data Analysis Agent
**Tags:** `["data", "analysis", "statistics", "ml", "visualization"]`  
**Enable:** Set `ENABLE_DATA_ANALYSIS_AGENT=true` or `ENABLE_EXAMPLE_AGENTS=true`  
**Purpose:** Specialist in data science, statistics, and machine learning

### 📚 Default Agent
**Priority:** 100 (Lowest - Fallback)  
**Tags:** `["general", "default", "fallback"]`  
**Purpose:** General conversation agent using OpenAI models for all other queries

## 🎮 Agent Commands

### Runtime Agent Management
- `/status` - Shows all agents with priority, enabled status, and tags
- `/agents` - Lists all available agents with descriptions
- `/enable <agent>` - Enable specific agent
- `/disable <agent>` - Disable specific agent
- `/solo <agent>` - Enable only the specified agent (disable all others)
- `/unsolo` - Re-enable all agents (exit solo mode)
- `/tag <tag>` - Find agents by tag
- `/config` - Show current configuration

### Basic Commands
- `/model` - Change AI model during conversation
- `/store` - Save the last conversation to file
- `/load` - Load a previous conversation and continue
- `/list` - View all saved conversations
- `quit` / `exit` - Exit the application

## 🔧 Configuration

### Environment Variables

**Required:**
- `OPENAI_API_KEY` - OpenAI API key

**Optional:**
- `WEATHER_API_KEY` - WeatherAPI.com API key for real-time weather data
- `ENABLE_EXAMPLE_AGENTS=true` - Enable all example agents
- `ENABLE_CODE_REVIEW_AGENT=true` - Enable just code review agent
- `ENABLE_DATA_ANALYSIS_AGENT=true` - Enable just data analysis agent
- `AGENT_PRIORITIES="Math:1,Weather:5"` - Custom priority overrides

### Setup Example
```bash
# Required
export OPENAI_API_KEY="your-openai-api-key"

# Optional - for enhanced weather
export WEATHER_API_KEY="your-weatherapi-key"

# Optional - for example agents
export ENABLE_EXAMPLE_AGENTS=true
```

## 🧪 Testing the System

### Demo Session Flow
1. Start the application: `./allday-term-agent`
2. Try different agent types:
   ```
   💬 You: /agents
   💬 You: What's 25% of 180?  # Math Agent
   💬 You: Weather in Tokyo?   # Weather Agent  
   💬 You: Explain quantum computing  # Default Agent
   ```
3. Test agent management:
   ```
   💬 You: /status
   💬 You: /disable Math
   💬 You: /enable Math
   💬 You: /tag weather
   ```

### Visual Indicators
The system provides clear feedback about which agent is handling each request:
- `🎯 Routing to Math Agent` - Shows when Math Agent is selected
- `🌤️ [Enhanced Weather Agent]` - Weather agent with real-time data
- `🤖 Assistant:` - Default agent for general queries

## 🚀 Extending the System

### Adding a New Agent

1. **Create the agent package:**
```go
// agents/myagent/myagent.go
package myagent

import "allday-term-agent/agents"

type MyAgent struct{}

func New() agents.Agent {
    return &MyAgent{}
}

func (a *MyAgent) CanHandle(prompt string) bool {
    // Add your detection logic
    return strings.Contains(strings.ToLower(prompt), "mykeword")
}

func (a *MyAgent) Handle(prompt, model string) agents.AgentResult {
    // Add your agent logic
    return agents.AgentResult{
        Response: "My custom response",
        Agent:    "My Agent",
    }
}

func (a *MyAgent) GetDescription() string {
    return "My custom agent for specific tasks"
}
```

2. **Register in agent factory:**
```go
// agent_factory.go
import myagent "allday-term-agent/agents/myagent"

// In registerDefaultAgents():
router.RegisterAgent("My Agent", myagent.New(), 5, true, []string{"mytag", "custom"})
```

## 🎯 Agent Priority System

Agents are evaluated in priority order (lowest number = highest priority):

1. **Priority 1** - Math Agent (highest priority)
2. **Priority 5** - Weather Agents, Example Agents  
3. **Priority 10** - Future specialized agents
4. **Priority 100** - Default Agent (fallback)

This ensures that specialized agents get first chance to handle relevant prompts before falling back to the general-purpose Default Agent.

## 🧠 How Routing Works

1. **Input Analysis** - System analyzes the user's prompt
2. **Priority Evaluation** - Agents are checked in priority order
3. **Agent Selection** - First agent that can handle the prompt is selected
4. **Response Generation** - Selected agent processes the prompt
5. **Fallback** - If no specialized agent matches, Default Agent handles it

The routing is completely automatic and transparent to the user, providing specialized expertise when needed while maintaining general conversation capabilities.
