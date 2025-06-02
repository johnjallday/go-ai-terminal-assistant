package storage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ConversationFile represents a saved conversation file
type ConversationFile struct {
	Filename    string
	Timestamp   time.Time
	DisplayName string
}

// ListConversationFiles returns a list of saved conversation files
func ListConversationFiles() ([]ConversationFile, error) {
	files, err := os.ReadDir("responses")
	if err != nil {
		if os.IsNotExist(err) {
			return []ConversationFile{}, nil // No responses directory yet
		}
		return nil, fmt.Errorf("failed to read responses directory: %w", err)
	}

	var conversations []ConversationFile
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "conversation_") && strings.HasSuffix(file.Name(), ".txt") {
			// Parse timestamp from filename
			timestampStr := strings.TrimPrefix(file.Name(), "conversation_")
			timestampStr = strings.TrimSuffix(timestampStr, ".txt")

			timestamp, err := time.Parse("2006-01-02_15-04-05", timestampStr)
			if err != nil {
				continue // Skip files with invalid timestamp format
			}

			conversations = append(conversations, ConversationFile{
				Filename:    file.Name(),
				Timestamp:   timestamp,
				DisplayName: timestamp.Format("2006-01-02 15:04:05"),
			})
		}
	}

	// Sort by timestamp (newest first)
	for i := 0; i < len(conversations)-1; i++ {
		for j := i + 1; j < len(conversations); j++ {
			if conversations[i].Timestamp.Before(conversations[j].Timestamp) {
				conversations[i], conversations[j] = conversations[j], conversations[i]
			}
		}
	}

	return conversations, nil
}

// SelectConversationFile displays available conversation files and lets user choose one
func SelectConversationFile() (string, error) {
	conversations, err := ListConversationFiles()
	if err != nil {
		return "", err
	}

	if len(conversations) == 0 {
		return "", fmt.Errorf("no saved conversations found")
	}

	fmt.Println("\n📂 Saved Conversations:")
	fmt.Println("─────────────────────")

	for i, conv := range conversations {
		fmt.Printf("%d. %s\n", i+1, conv.DisplayName)
	}

	fmt.Printf("\nSelect a conversation (1-%d) or press Enter to cancel: ", len(conversations))

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		choice := strings.TrimSpace(scanner.Text())

		// Cancel if empty
		if choice == "" {
			return "", nil
		}

		// Parse user choice
		if num, err := strconv.Atoi(choice); err == nil && num >= 1 && num <= len(conversations) {
			return filepath.Join("responses", conversations[num-1].Filename), nil
		}
	}

	return "", fmt.Errorf("invalid selection")
}

// LoadConversation loads a conversation from file and returns the last prompt and response
func LoadConversation(filename string) (string, string, string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var model, prompt, response string
	var inResponse bool

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Model: ") {
			model = strings.TrimPrefix(line, "Model: ")
		} else if strings.HasPrefix(line, "User: ") {
			prompt = strings.TrimPrefix(line, "User: ")
			inResponse = false
		} else if strings.HasPrefix(line, "Assistant: ") {
			response = strings.TrimPrefix(line, "Assistant: ")
			inResponse = true
		} else if inResponse && line != "" && !strings.HasPrefix(line, "=") {
			// Continue reading multi-line response
			if response != "" {
				response += "\n" + line
			} else {
				response = line
			}
		}
	}

	return model, prompt, response, nil
}

// StoreOpenAIResponse saves a conversation exchange to a file
func StoreOpenAIResponse(prompt, response, model string) error {
	// Create responses directory if it doesn't exist
	if err := os.MkdirAll("responses", 0755); err != nil {
		return fmt.Errorf("failed to create responses directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("responses/conversation_%s.txt", timestamp)

	// Prepare content
	content := fmt.Sprintf("Model: %s\nTimestamp: %s\n\n", model, time.Now().Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("User: %s\n\n", prompt)
	content += fmt.Sprintf("Assistant: %s\n", response)
	content += "\n" + strings.Repeat("=", 50) + "\n\n"

	// Write to file (append mode)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
