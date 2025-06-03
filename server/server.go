package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go-ai-terminal-assistant/models"
	"go-ai-terminal-assistant/utils"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// DaemonServer represents the HTTP server for the AI assistant daemon
type DaemonServer struct {
	client openai.Client
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
	Model    string `json:"model"`
	Error    string `json:"error,omitempty"`
	Success  bool   `json:"success"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// NewDaemonServer creates a new daemon server instance
func NewDaemonServer(port string) (*DaemonServer, error) {
	// Load environment variables
	if err := loadEnvFile(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	defaultModel := models.GetDefaultModel()

	return &DaemonServer{
		client: client,
		model:  defaultModel,
		port:   port,
	}, nil
}

// Start starts the HTTP server
func (s *DaemonServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/chat", s.handleChat)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/", s.handleRoot)

	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 AI Terminal Assistant daemon starting on port %s", s.port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("🛑 Shutting down daemon server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}

// handleChat handles chat requests
func (s *DaemonServer) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		s.sendErrorResponse(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Use provided model or default
	model := req.Model
	if model == "" {
		model = s.model
	}

	// Use the existing utils.GetOpenAIResponse function for consistency
	response, err := utils.GetOpenAIResponse(&s.client, req.Message, model)
	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		s.sendErrorResponse(w, fmt.Sprintf("AI service error: %v", err), http.StatusInternalServerError)
		return
	}

	// Send successful response
	chatResp := ChatResponse{
		Response: response,
		Model:    model,
		Success:  true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResp)
}

// handleHealth handles health check requests
func (s *DaemonServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleRoot handles requests to the root endpoint
func (s *DaemonServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"service": "AI Terminal Assistant Daemon",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": map[string]string{
			"POST /chat":  "Send chat messages to AI",
			"GET /health": "Health check",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse sends an error response
func (s *DaemonServer) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ChatResponse{
		Error:   message,
		Success: false,
	}

	json.NewEncoder(w).Encode(response)
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}

	return scanner.Err()
}
