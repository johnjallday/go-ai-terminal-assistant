package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-ai-terminal-assistant/models"
	"go-ai-terminal-assistant/utils"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Router interface defines the methods needed by the server
type Router interface {
	RoutePrompt(prompt string) interface{} // Returns agent
	ListAgents() []interface{}             // Returns list of agents
	ListAllAgents() []interface{}          // Returns all agent registrations
	EnableAgent(agentName string, enabled bool) bool
	SoloAgent(agentName string) bool
	UnsoloAgents()
}

// Agent interface defines the methods needed for agents
type Agent interface {
	GetName() string
	GetDescription() string
}

// Server represents the HTTP server for the AI assistant
type Server struct {
	client *openai.Client
	router Router
	model  string
	port   string
}

// ChatRequest represents the JSON request structure
type ChatRequest struct {
	Message string `json:"message"`
	Model   string `json:"model,omitempty"`
}

// ChatResponse represents the JSON response structure
type ChatResponse struct {
	Response string `json:"response"`
	Agent    string `json:"agent"`
	Model    string `json:"model"`
	Error    string `json:"error,omitempty"`
	Success  bool   `json:"success"`
}

// AgentStatusResponse represents the agent status
type AgentStatusResponse struct {
	Agents []AgentInfo `json:"agents"`
}

// AgentInfo represents individual agent information
type AgentInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Priority    int      `json:"priority"`
	Tags        []string `json:"tags"`
}

// AgentRequest represents requests for agent operations
type AgentRequest struct {
	Agent string `json:"agent"`
}

// NewServer creates a new HTTP server instance
func NewServer(port string) *Server {
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

	return &Server{
		client: client,
		router: nil, // Will be set by SetRouter
		model:  models.GetDefaultModel(),
		port:   port,
	}
}

// SetRouter sets the agent router for the server
func (s *Server) SetRouter(router Router) {
	s.router = router
}

// Start starts the HTTP server
func (s *Server) Start() {
	http.HandleFunc("/", s.handleRoot)
	http.HandleFunc("/chat", s.handleChat)
	http.HandleFunc("/agents", s.handleAgents)
	http.HandleFunc("/agents/status", s.handleAgentStatus)
	http.HandleFunc("/agents/enable", s.handleEnableAgent)
	http.HandleFunc("/agents/disable", s.handleDisableAgent)
	http.HandleFunc("/agents/solo", s.handleSoloAgent)
	http.HandleFunc("/agents/unsolo", s.handleUnsoloAgents)
	http.HandleFunc("/health", s.handleHealth)

	fmt.Printf("🚀 AI Terminal Assistant Server starting on port %s\n", s.port)
	fmt.Printf("📡 API Endpoints:\n")
	fmt.Printf("   POST /chat - Send messages to AI agents\n")
	fmt.Printf("   GET  /agents - List available agents\n")
	fmt.Printf("   GET  /agents/status - Get agent status\n")
	fmt.Printf("   POST /agents/enable - Enable specific agent\n")
	fmt.Printf("   POST /agents/disable - Disable specific agent\n")
	fmt.Printf("   POST /agents/solo - Solo mode for specific agent\n")
	fmt.Printf("   POST /agents/unsolo - Exit solo mode\n")
	fmt.Printf("   GET  /health - Health check\n")
	fmt.Printf("\n💡 Example: curl -X POST http://localhost:%s/chat -H 'Content-Type: application/json' -d '{\"message\":\"What is 2+2?\"}'\n", s.port)

	log.Fatal(http.ListenAndServe(":"+s.port, nil))
}

// handleRoot provides API documentation
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"service": "AI Terminal Assistant API",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"POST /chat":           "Send messages to AI agents",
			"GET /agents":          "List available agents",
			"GET /agents/status":   "Get agent status",
			"POST /agents/enable":  "Enable specific agent",
			"POST /agents/disable": "Disable specific agent",
			"POST /agents/solo":    "Solo mode for specific agent",
			"POST /agents/unsolo":  "Exit solo mode",
			"GET /health":          "Health check",
		},
	}
	json.NewEncoder(w).Encode(response)
}

// handleChat processes chat messages
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Use provided model or default
	model := req.Model
	if model == "" {
		model = s.model
	}

	w.Header().Set("Content-Type", "application/json")

	var response string
	var agentName string
	var err error

	if s.router != nil {
		// Route the prompt to an appropriate agent
		agent := s.router.RoutePrompt(req.Message)
		if agentObj, ok := agent.(Agent); ok {
			agentName = agentObj.GetName()
			response, err = utils.GetOpenAIResponse(s.client, req.Message, model)
		} else {
			agentName = "Default"
			response, err = utils.GetOpenAIResponse(s.client, req.Message, model)
		}
	} else {
		agentName = "Default"
		response, err = utils.GetOpenAIResponse(s.client, req.Message, model)
	}

	chatResponse := ChatResponse{
		Model:   model,
		Agent:   agentName,
		Success: err == nil,
	}

	if err != nil {
		chatResponse.Error = err.Error()
		chatResponse.Response = ""
	} else {
		chatResponse.Response = response
	}

	json.NewEncoder(w).Encode(chatResponse)
}

// handleAgents lists available agents
func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var agents []AgentInfo
	if s.router != nil {
		agentList := s.router.ListAgents()
		for _, agent := range agentList {
			if agentObj, ok := agent.(Agent); ok {
				agents = append(agents, AgentInfo{
					Name:        agentObj.GetName(),
					Description: agentObj.GetDescription(),
					Enabled:     true,
					Priority:    0,
					Tags:        []string{},
				})
			}
		}
	}

	json.NewEncoder(w).Encode(AgentStatusResponse{Agents: agents})
}

// handleAgentStatus provides detailed agent status
func (s *Server) handleAgentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var agents []AgentInfo
	if s.router != nil {
		// This would need to be implemented in the router interface
		// For now, return basic info
		agentList := s.router.ListAgents()
		for _, agent := range agentList {
			if agentObj, ok := agent.(Agent); ok {
				agents = append(agents, AgentInfo{
					Name:        agentObj.GetName(),
					Description: agentObj.GetDescription(),
					Enabled:     true,
					Priority:    0,
					Tags:        []string{},
				})
			}
		}
	}

	json.NewEncoder(w).Encode(AgentStatusResponse{Agents: agents})
}

// handleEnableAgent enables a specific agent
func (s *Server) handleEnableAgent(w http.ResponseWriter, r *http.Request) {
	s.handleAgentToggle(w, r, true)
}

// handleDisableAgent disables a specific agent
func (s *Server) handleDisableAgent(w http.ResponseWriter, r *http.Request) {
	s.handleAgentToggle(w, r, false)
}

// handleAgentToggle handles enabling/disabling agents
func (s *Server) handleAgentToggle(w http.ResponseWriter, r *http.Request, enable bool) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Agent == "" {
		http.Error(w, "Agent name is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var success bool
	if s.router != nil {
		success = s.router.EnableAgent(req.Agent, enable)
	}

	action := "enabled"
	if !enable {
		action = "disabled"
	}

	response := map[string]interface{}{
		"success": success,
		"message": fmt.Sprintf("Agent %s %s", req.Agent, action),
	}

	json.NewEncoder(w).Encode(response)
}

// handleSoloAgent puts a specific agent in solo mode
func (s *Server) handleSoloAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Agent == "" {
		http.Error(w, "Agent name is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var success bool
	if s.router != nil {
		success = s.router.SoloAgent(req.Agent)
	}

	response := map[string]interface{}{
		"success": success,
		"message": fmt.Sprintf("Solo mode for agent %s", req.Agent),
	}

	json.NewEncoder(w).Encode(response)
}

// handleUnsoloAgents exits solo mode
func (s *Server) handleUnsoloAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if s.router != nil {
		s.router.UnsoloAgents()
	}

	response := map[string]interface{}{
		"success": true,
		"message": "All agents re-enabled",
	}

	json.NewEncoder(w).Encode(response)
}

// handleHealth provides health check
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "AI Terminal Assistant",
		"model":   s.model,
	}
	json.NewEncoder(w).Encode(response)
}

// loadEnvFile loads environment variables from .env file if it exists
func loadEnvFile() {
	// Implementation would be similar to main.go
	// For now, we'll assume environment variables are set
}
