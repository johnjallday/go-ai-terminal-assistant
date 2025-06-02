# Quick Reference Guide

## 🚀 Starting the Application

```bash
# Set your OpenAI API key
export OPENAI_API_KEY='your-openai-api-key-here'

# Optional: Set weather API key for real-time data
export WEATHER_API_KEY='your-weatherapi-key-here'

# Run the application (starts with GPT-4.1 Nano by default)
./allday-term-agent
```

**Note**: The application now starts with GPT-4.1 Nano by default. Use `/model` to change models during your session.

## 🎯 Available Commands

| Command | Description |
|---------|-------------|
| `/agents` | List all available agents |
| `/status` | Show agent status and priorities |
| `/enable <agent>` | Enable a specific agent |
| `/disable <agent>` | Disable a specific agent |
| `/solo <agent>` | Enable only the specified agent (disable all others) |
| `/unsolo` | Re-enable all agents (exit solo mode) |
| `/tag <tag>` | Find agents by category |
| `/config` | Show current configuration |
| `/model` | Change the AI model |
| `/store` | Save the last conversation |
| `/load` | Load a saved conversation |
| `/list` | List saved conversations |
| `quit` or `exit` | Exit the application |

## 🤖 Available AI Models

1. **GPT-4.1** - Latest flagship model ($2.00/$8.00 per 1M tokens)
2. **GPT-4.1 Mini** - Balanced performance and cost ($0.40/$1.60 per 1M tokens) [Default]
3. **GPT-4.1 Nano** - Ultra-affordable ($0.10/$0.40 per 1M tokens)
4. **GPT-4.5 Preview** - Most advanced preview ($75.00/$150.00 per 1M tokens)
5. **O4 Mini** - Efficient reasoning ($1.10/$4.40 per 1M tokens)
6. **O3 Mini** - Advanced reasoning ($1.10/$4.40 per 1M tokens)

## 🧠 Agent Types

### 🧮 Math Agent
- **Handles:** Calculations, equations, mathematical concepts
- **Triggers:** Keywords like calculate, solve, math, +, -, *, /, =, etc.
- **Examples:**
  - `calculate 25 * 17`
  - `what is 15% of 200?`
  - `solve for x: 2x + 5 = 15`

### 🌤️ Weather Agents
- **Enhanced Weather Agent:** Real-time data with WEATHER_API_KEY
- **Basic Weather Agent:** AI-powered responses (fallback)
- **Triggers:** Keywords like weather, temperature, rain, forecast, etc.
- **Examples:**
  - `what's the weather like today?`
  - `temperature in New York`
  - `weather forecast for London`

### 🔍 Example Agents (Optional)
- **Code Review Agent:** Code analysis and optimization
- **Data Analysis Agent:** Statistics and ML specialist  
- **Enable:** Set environment variables (`ENABLE_EXAMPLE_AGENTS=true`)

### 🤖 Default Agent
- **Handles:** General conversation, creative tasks, knowledge questions
- **Examples:**
  - `hello, how are you?`
  - `tell me a joke`
  - `what's the capital of France?`

## 💡 Tips

- The system automatically routes your questions to the best agent
- Use `/status` to see agent priorities and enabled status
- Use `/agents` to list all available agents and capabilities
- Use `/tag <tag>` to find agents by category (math, weather, code, etc.)
- Conversation history is maintained for context
- Saved conversations preserve model and context information
- Weather agent automatically selects Enhanced/Basic based on API key availability
- Enable example agents with environment variables for specialized tasks

## 🔧 Environment Setup

The application automatically loads environment variables from a `.env` file if present.

Create a `.env` file or set environment variables:

```bash
# Required
OPENAI_API_KEY=your-openai-api-key-here

# Optional (for real-time weather data)  
WEATHER_API_KEY=your-weatherapi-key-here

# Optional (for example agents)
ENABLE_EXAMPLE_AGENTS=true
ENABLE_CODE_REVIEW_AGENT=true
ENABLE_DATA_ANALYSIS_AGENT=true
```

### Weather API Auto-Setup
When you first request weather data without an API key, the weather agent will:
1. 🔑 Prompt you to enter a WeatherAPI.com key
2. 💾 Automatically save it to your `.env` file
3. ✅ Enable real-time weather data for all future sessions

No manual setup required!

## 📁 File Structure

```
allday-term-agent/
├── main.go                    # Main application entry point
├── agent_factory.go          # Agent registration and factory
├── router.go                 # Agent routing logic
├── models.go                 # Model selection and management
├── storage.go                # Conversation persistence
├── openai.go                 # OpenAI API wrapper
├── agents/                   # Modular agent packages
│   ├── interface.go         #   Shared Agent interface
│   ├── default/            #   General conversation agent
│   ├── math/               #   Mathematical calculations
│   ├── weather/            #   Weather data and forecasts
│   └── examples/           #   Optional specialized agents
├── utils/                   # Shared utilities
│   └── openai.go           #   OpenAI API utilities
├── docs/                    # Documentation directory
│   ├── AGENTS.md           #   Agent system documentation
│   ├── QUICK_REFERENCE.md  #   This quick reference guide
│   ├── README_KR.md        #   Korean documentation
│   └── REFACTORING_COMPLETION_REPORT.md  # Technical details
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── README.md                # Main project documentation
├── test_agent_routing.sh    # Agent routing test script
├── test.sh                  # Basic build and functionality test
├── .env.example             # Environment variable template
├── .gitignore               # Git ignore rules
└── responses/               # Saved conversations
```
