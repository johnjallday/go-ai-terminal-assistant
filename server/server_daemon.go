package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"go-ai-terminal-assistant/models"
	"go-ai-terminal-assistant/utils"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Server represents the HTTP server for the AI assistant
type Server struct {
	client *openai.Client
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
		client: &client,
		model:  models.GetDefaultModel(),
		port:   port,
	}
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
	model := s.model
	if req.Model != "" {
		model = req.Model
	}

	w.Header().Set("Content-Type", "application/json")

	// For daemon mode, use direct OpenAI API without agent routing
	response, err := utils.GetOpenAIResponse(s.client, req.Message, model)

	if err != nil {
		chatResp := ChatResponse{
			Success: false,
			Error:   err.Error(),
			Agent:   "Default",
			Model:   model,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(chatResp)
		return
	}

	chatResp := ChatResponse{
		Success:  true,
		Response: response,
		Agent:    "Default",
		Model:    model,
	}
	json.NewEncoder(w).Encode(chatResp)
}

// handleAgents lists available agents
func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// For daemon mode, return basic agent info
	agentInfos := []AgentInfo{
		{
			Name:        "Default",
			Description: "General AI assistant",
			Enabled:     true,
			Priority:    1,
			Tags:        []string{"general", "default"},
		},
	}

	json.NewEncoder(w).Encode(map[string][]AgentInfo{"agents": agentInfos})
}

// handleAgentStatus returns detailed agent status
func (s *Server) handleAgentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// For daemon mode, return basic status
	agentInfos := []AgentInfo{
		{
			Name:        "Default",
			Description: "General AI assistant",
			Enabled:     true,
			Priority:    1,
			Tags:        []string{"general", "default"},
		},
	}

	json.NewEncoder(w).Encode(AgentStatusResponse{Agents: agentInfos})
}

// handleEnableAgent enables a specific agent
func (s *Server) handleEnableAgent(w http.ResponseWriter, r *http.Request) {
	s.handleAgentToggle(w, r, true)
}

// handleDisableAgent disables a specific agent
func (s *Server) handleDisableAgent(w http.ResponseWriter, r *http.Request) {
	s.handleAgentToggle(w, r, false)
}

// handleAgentToggle handles enable/disable agent requests
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

	// For daemon mode, simulate success for "Default" agent
	success := req.Agent == "Default"
	action := "enabled"
	if !enable {
		action = "disabled"
	}

	response := map[string]interface{}{
		"success": success,
		"agent":   req.Agent,
		"action":  action,
	}

	if !success {
		response["error"] = "Agent not found"
		w.WriteHeader(http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(response)
}

// handleSoloAgent enables solo mode for a specific agent
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

	// For daemon mode, simulate success for "Default" agent
	success := req.Agent == "Default"

	response := map[string]interface{}{
		"success": success,
		"agent":   req.Agent,
		"action":  "solo_enabled",
	}

	if !success {
		response["error"] = "Agent not found"
		w.WriteHeader(http.StatusNotFound)
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
	response := map[string]interface{}{
		"success": true,
		"action":  "solo_disabled",
		"message": "All agents re-enabled",
	}
	json.NewEncoder(w).Encode(response)
}

// handleHealth provides health check endpoint
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
