package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go-ai-terminal-assistant/agents"
	"go-ai-terminal-assistant/models"
	"go-ai-terminal-assistant/router"
	"go-ai-terminal-assistant/storage"
	"go-ai-terminal-assistant/utils"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// loadEnvFile loads environment variables from .env file if it exists
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		// .env file doesn't exist, which is fine
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Only set if environment variable is not already set
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

func main() {
	// Load environment variables from .env file if it exists
	loadEnvFile()

	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	// Create agent factory and initialize router
	factory := router.NewAgentFactory()
	agentRouter := factory.CreateAgentRouter()

	fmt.Println("🤖 OpenAI Terminal Assistant with Agentic Routing")

	// Use default model (GPT-4.1 Nano) instead of prompting
	selectedModel := models.GetDefaultModel()
	modelDisplayName := models.GetModelDisplayName(selectedModel)

	fmt.Printf("\n✨ Using model: %s\n", modelDisplayName)
	fmt.Println("Commands: 'quit', '/model', '/agents', '/tools', '/status', '/config', '/enable <agent>', '/disable <agent>', '/solo <agent>', '/unsolo', '/tag <tag>', '/store', '/load', '/list', '/clear'")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")

	scanner := bufio.NewScanner(os.Stdin)
	var lastPrompt, lastResponse string

	for {
		fmt.Print("\n💬 You: ")

		// Read user input
		if !scanner.Scan() {
			break
		}

		rawInput := scanner.Text()
		if rawInput == "\f" {
			fmt.Print("\033[H\033[2J")
			continue
		}
		input := strings.TrimSpace(rawInput)

		// Check for exit commands
		if input == "quit" || input == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		// Check for model change command
		if input == "/model" {
			selectedModel = models.SelectModel()
			modelDisplayName = models.GetModelDisplayName(selectedModel)
			fmt.Printf("\n✨ Switched to model: %s\n", modelDisplayName)
			continue
		}

		// Check for store command
		if input == "/store" {
			if lastPrompt != "" && lastResponse != "" {
				if err := storage.StoreOpenAIResponse(lastPrompt, lastResponse, selectedModel); err != nil {
					fmt.Printf("❌ Error storing response: %v\n", err)
				} else {
					fmt.Println("💾 Response saved to file!")
				}
			} else {
				fmt.Println("❌ No previous conversation to store.")
			}
			continue
		}

		// Check for load command
		if input == "/load" {
			filename, err := storage.SelectConversationFile()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}

			if filename != "" {
				// Load and display the conversation
				model, prompt, response, err := storage.LoadConversation(filename)
				if err != nil {
					fmt.Printf("❌ Error loading conversation: %v\n", err)
				} else {
					fmt.Printf("📂 Loaded conversation from %s\n", filename)
					fmt.Printf("Model: %s\nUser: %s\nAssistant: %s\n", model, prompt, response)

					// Set the loaded conversation as the current context
					lastPrompt = prompt
					lastResponse = response

					// Switch to the model used in the loaded conversation
					selectedModel = model
					modelDisplayName = models.GetModelDisplayName(selectedModel)
					fmt.Printf("\n✨ Switched to model: %s\n", modelDisplayName)
					fmt.Println("💡 You can now continue the conversation from where it left off!")
				}
			}
			continue
		}

		// Check for clear screen command
		if input == "/clear" {
			// ANSI escape to clear terminal screen
			fmt.Print("\033[H\033[2J")
			continue
		}

		// Check for agents command
		if input == "/agents" {
			fmt.Println("\n🤖 Available Agents:")
			fmt.Println("─────────────────────")
			agents := agentRouter.ListAgents()
			for i, agent := range agents {
				fmt.Printf("%d. %s Agent - %s\n", i+1, agent.GetName(), agent.GetDescription())
			}
			continue
		}

		// Check for tools command
		if input == "/tools" {
			fmt.Println("\n🛠️ Available tools:")
			fmt.Println("─────────────────")
			for _, ag := range agentRouter.ListAgents() {
				if tp, ok := ag.(agents.ToolProvider); ok {
					for _, tool := range tp.Tools() {
						fmt.Printf(" - %s: %s (%s agent)\n", tool.Name, tool.Description, ag.GetName())
					}
				}
			}
			continue
		}

		// Check for agent status command
		if input == "/status" {
			fmt.Println(agentRouter.GetAgentStatus())
			continue
		}

		// Check for agent enable/disable commands
		if strings.HasPrefix(input, "/enable ") {
			agentName := strings.TrimSpace(strings.TrimPrefix(input, "/enable "))
			if agentRouter.EnableAgent(agentName, true) {
				fmt.Printf("✅ %s Agent enabled\n", agentName)
			} else {
				fmt.Printf("❌ Agent '%s' not found\n", agentName)
			}
			continue
		}

		if strings.HasPrefix(input, "/disable ") {
			agentName := strings.TrimSpace(strings.TrimPrefix(input, "/disable "))
			if agentRouter.EnableAgent(agentName, false) {
				fmt.Printf("❌ %s Agent disabled\n", agentName)
			} else {
				fmt.Printf("❌ Agent '%s' not found\n", agentName)
			}
			continue
		}

		// Check for solo agent command
		if strings.HasPrefix(input, "/solo ") {
			agentName := strings.TrimSpace(strings.TrimPrefix(input, "/solo "))
			if agentRouter.SoloAgent(agentName) {
				fmt.Printf("🎯 Solo mode: Only %s Agent is enabled\n", agentName)
				fmt.Println("💡 Use '/unsolo' to re-enable all agents")
				// Print available tools for the soloed agent, if any
				for _, ag := range agentRouter.ListAgents() {
					if tp, ok := ag.(agents.ToolProvider); ok {
						tools := tp.Tools()
						if len(tools) > 0 {
							fmt.Println()
							fmt.Println("🛠️  Available tools:")
							for _, tool := range tools {
								fmt.Printf(" - %s: %s\n", tool.Name, tool.Description)
							}
						}
					}
				}
			} else {
				fmt.Printf("❌ Agent '%s' not found\n", agentName)
			}
			continue
		}

		// Check for unsolo command
		if input == "/unsolo" {
			agentRouter.UnsoloAgents()
			fmt.Println("✅ All agents re-enabled")
			continue
		}

		// Check for tag-based agent listing
		if strings.HasPrefix(input, "/tag ") {
			tag := strings.TrimSpace(strings.TrimPrefix(input, "/tag "))
			agents := agentRouter.GetAgentsByTag(tag)
			if len(agents) == 0 {
				fmt.Printf("🏷️ No agents found with tag '%s'\n", tag)
			} else {
				fmt.Printf("\n🏷️ Agents with tag '%s':\n", tag)
				fmt.Println("─────────────────────")
				for i, agent := range agents {
					fmt.Printf("%d. %s Agent - %s\n", i+1, agent.GetName(), agent.GetDescription())
				}
			}
			continue
		}

		// Check for configuration command
		if input == "/config" {
			config := factory.GetConfig()
			fmt.Println("\n⚙️ Agent Configuration:")
			fmt.Println("─────────────────────")
			fmt.Printf("Example Agents Enabled: %v\n", config.EnableExampleAgents)
			fmt.Printf("Code Review Agent: %v\n", config.EnableCodeReview)
			fmt.Printf("Data Analysis Agent: %v\n", config.EnableDataAnalysis)
			fmt.Printf("Weather API Key: %s\n", func() string {
				if config.WeatherAPIKey != "" {
					return "✅ Configured"
				}
				return "❌ Not set"
			}())
			if len(config.CustomAgentPriority) > 0 {
				fmt.Println("Custom Priorities:")
				for agent, priority := range config.CustomAgentPriority {
					fmt.Printf("  %s: %d\n", agent, priority)
				}
			}
			continue
		}

		// Check for list command
		if input == "/list" {
			conversations, err := storage.ListConversationFiles()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}

			if len(conversations) == 0 {
				fmt.Println("📂 No saved conversations found.")
			} else {
				fmt.Println("\n📂 Saved Conversations:")
				fmt.Println("─────────────────────")
				for i, conv := range conversations {
					fmt.Printf("%d. %s\n", i+1, conv.DisplayName)
				}
			}
			continue
		}

		// Skip empty input
		if input == "" {
			continue
		}

		// Route the prompt to the appropriate agent
		agent := agentRouter.RoutePrompt(input)

		var response string
		var err error

		// Use the selected agent to handle the prompt
		if agent.GetName() == "Default" {
			// For default agent, include conversation history if available
			if lastPrompt != "" && lastResponse != "" {
				response, err = utils.GetOpenAIResponse(&client, input, selectedModel, lastPrompt, lastResponse)
			} else {
				response, err = utils.GetOpenAIResponse(&client, input, selectedModel)
			}
		} else {
			// For specialized agents, they handle their own context
			response, err = agent.Handle(input, &client, selectedModel)
		}

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("🤖 Assistant: %s\n", response)

		// Store the last conversation for potential saving
		lastPrompt = input
		lastResponse = response
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}
