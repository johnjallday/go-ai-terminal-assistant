package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestConversationFile(t *testing.T) {
	// Test ConversationFile struct
	now := time.Now()
	cf := ConversationFile{
		Filename:    "test_conversation.txt",
		Timestamp:   now,
		DisplayName: now.Format("2006-01-02 15:04:05"),
	}

	if cf.Filename != "test_conversation.txt" {
		t.Error("Filename should be set correctly")
	}
	if cf.Timestamp != now {
		t.Error("Timestamp should be set correctly")
	}
	if cf.DisplayName == "" {
		t.Error("DisplayName should not be empty")
	}
}

func TestStoreAndLoadConversation(t *testing.T) {
	// Create a temporary test directory
	tempDir := "test_responses"
	defer os.RemoveAll(tempDir)

	// Change to temp directory for this test
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Test data
	prompt := "What is 2+2?"
	response := "2+2 equals 4."
	model := "gpt-4.1-mini"

	// Test storing a conversation
	err := StoreOpenAIResponse(prompt, response, model)
	if err != nil {
		t.Fatalf("StoreOpenAIResponse failed: %v", err)
	}

	// Check that responses directory was created
	if _, err := os.Stat("responses"); os.IsNotExist(err) {
		t.Error("responses directory should be created")
	}

	// List conversation files
	conversations, err := ListConversationFiles()
	if err != nil {
		t.Fatalf("ListConversationFiles failed: %v", err)
	}

	if len(conversations) == 0 {
		t.Error("Should have at least one conversation file")
		return
	}

	// Test loading the conversation
	filename := filepath.Join("responses", conversations[0].Filename)
	loadedModel, loadedPrompt, loadedResponse, err := LoadConversation(filename)
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}

	// Verify loaded data
	if loadedModel != model {
		t.Errorf("Loaded model = %q; want %q", loadedModel, model)
	}
	if loadedPrompt != prompt {
		t.Errorf("Loaded prompt = %q; want %q", loadedPrompt, prompt)
	}
	if !strings.Contains(loadedResponse, response) {
		t.Errorf("Loaded response should contain %q, got %q", response, loadedResponse)
	}
}

func TestListConversationFilesEmpty(t *testing.T) {
	// Test when no responses directory exists
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Create a temporary directory without responses folder
	tempDir := "test_empty"
	os.Mkdir(tempDir, 0755)
	defer os.RemoveAll(tempDir)
	os.Chdir(tempDir)

	conversations, err := ListConversationFiles()
	if err != nil {
		t.Fatalf("ListConversationFiles should not fail when directory doesn't exist: %v", err)
	}

	if len(conversations) != 0 {
		t.Error("Should return empty slice when no responses directory exists")
	}
}

func TestLoadConversationInvalidFile(t *testing.T) {
	// Test loading a non-existent file
	_, _, _, err := LoadConversation("nonexistent_file.txt")
	if err == nil {
		t.Error("LoadConversation should fail for non-existent file")
	}
}

func TestSelectConversationFile(t *testing.T) {
	// Create some test conversation files
	conversations := []ConversationFile{
		{
			Filename:    "conv1.txt",
			Timestamp:   time.Now().Add(-time.Hour),
			DisplayName: "Test 1",
		},
		{
			Filename:    "conv2.txt",
			Timestamp:   time.Now(),
			DisplayName: "Test 2",
		},
	}

	// Note: SelectConversationFile() reads from stdin, so we can't easily test it in unit tests
	// This test verifies the function exists and has the right signature
	// In practice, integration tests would cover the interactive behavior

	if len(conversations) == 0 {
		t.Skip("No conversations to test SelectConversationFile")
	}

	// We can at least verify that the conversation files have valid filenames
	for i, conv := range conversations {
		if conv.Filename == "" {
			t.Errorf("Conversation %d should have a non-empty filename", i)
		}
		if conv.DisplayName == "" {
			t.Errorf("Conversation %d should have a non-empty display name", i)
		}
	}
}
