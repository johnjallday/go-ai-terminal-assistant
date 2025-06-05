package chat

import (
	"bufio"
	"log"
	"os"
	"strings"

	"go-ai-terminal-assistant/models"
	"go-ai-terminal-assistant/router"
	"go-ai-terminal-assistant/utils"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// LoadEnvFile loads environment variables from .env file if it exists
func LoadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

// ChatSession holds state for conversational interactions
type ChatSession struct {
	Client           openai.Client
	Factory          *router.AgentFactory
	Router           *router.AgentRouter
	SelectedModel    string
	ModelDisplayName string
	LastPrompt       string
	LastResponse     string
}

// NewSession initializes and returns a new ChatSession
func NewSession() *ChatSession {
	LoadEnvFile()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}
	client := openai.NewClient(option.WithAPIKey(apiKey))
	factory := router.NewAgentFactory()
	agentRouter := factory.CreateAgentRouter()
	defaultModel := models.GetDefaultModel()
	modelDisplayName := models.GetModelDisplayName(defaultModel)
	return &ChatSession{
		Client:           client,
		Factory:          factory,
		Router:           agentRouter,
		SelectedModel:    defaultModel,
		ModelDisplayName: modelDisplayName,
	}
}

// ProcessMessage routes the message to the appropriate agent and returns the response
func (s *ChatSession) ProcessMessage(message string) (string, error) {
	agent := s.Router.RoutePrompt(message)
	var response string
	var err error
	if agent.GetName() == "Default" {
		if s.LastPrompt != "" && s.LastResponse != "" {
			response, err = utils.GetOpenAIResponse(&s.Client, message, s.SelectedModel, s.LastPrompt, s.LastResponse)
		} else {
			response, err = utils.GetOpenAIResponse(&s.Client, message, s.SelectedModel)
		}
	} else {
		response, err = agent.Handle(message, &s.Client, s.SelectedModel)
	}
	if err != nil {
		return "", err
	}
	s.LastPrompt = message
	s.LastResponse = response
	return response, nil
}
