package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"tlgbot/internal/models"
)

const expectedNoErrorMsg = "Expected no error, got %v"

func TestLoadFromEnvSuccess(t *testing.T) {
	// Save original values
	originalToken := os.Getenv(EnvTelegramToken)
	originalDelay := os.Getenv(EnvDelayMs)

	// Set test values
	if err := os.Setenv(EnvTelegramToken, "test_token_123"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv(EnvDelayMs, "1000"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	// Restore after test
	defer func() {
		if originalToken == "" {
			if err := os.Unsetenv(EnvTelegramToken); err != nil {
				t.Logf("Failed to unset environment variable: %v", err)
			}
		} else {
			if err := os.Setenv(EnvTelegramToken, originalToken); err != nil {
				t.Logf("Failed to restore environment variable: %v", err)
			}
		}
		if originalDelay == "" {
			if err := os.Unsetenv(EnvDelayMs); err != nil {
				t.Logf("Failed to unset environment variable: %v", err)
			}
		} else {
			if err := os.Setenv(EnvDelayMs, originalDelay); err != nil {
				t.Logf("Failed to restore environment variable: %v", err)
			}
		}
	}()

	config, err := LoadFromEnv()
	if err != nil {
		t.Errorf(expectedNoErrorMsg, err)
	}

	if config.TelegramToken != "test_token_123" {
		t.Errorf("Expected token 'test_token_123', got %s", config.TelegramToken)
	}

	if config.DelayMs != 1000 {
		t.Errorf("Expected delay 1000, got %d", config.DelayMs)
	}

	// Check default values
	if config.GoogleCreds != DefaultGoogleCreds {
		t.Errorf("Expected default GoogleCreds %s, got %s", DefaultGoogleCreds, config.GoogleCreds)
	}
}

func TestLoadFromEnvMissingToken(t *testing.T) {
	// Save original value
	originalToken := os.Getenv(EnvTelegramToken)

	// Remove token
	if err := os.Unsetenv(EnvTelegramToken); err != nil {
		t.Fatalf("Failed to unset environment variable: %v", err)
	}

	// Restore after test
	defer func() {
		if originalToken != "" {
			if err := os.Setenv(EnvTelegramToken, originalToken); err != nil {
				t.Logf("Failed to restore environment variable: %v", err)
			}
		}
	}()

	_, err := LoadFromEnv()
	if err == nil {
		t.Error("Expected error for missing token")
	}
}

func TestLoadFromEnvInvalidDelay(t *testing.T) {
	// Save original values
	originalToken := os.Getenv(EnvTelegramToken)
	originalDelay := os.Getenv(EnvDelayMs)

	// Set test values
	if err := os.Setenv(EnvTelegramToken, "test_token"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv(EnvDelayMs, "invalid_delay"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	// Restore after test
	defer func() {
		if originalToken == "" {
			if err := os.Unsetenv(EnvTelegramToken); err != nil {
				t.Logf("Failed to unset environment variable: %v", err)
			}
		} else {
			if err := os.Setenv(EnvTelegramToken, originalToken); err != nil {
				t.Logf("Failed to restore environment variable: %v", err)
			}
		}
		if originalDelay == "" {
			if err := os.Unsetenv(EnvDelayMs); err != nil {
				t.Logf("Failed to unset environment variable: %v", err)
			}
		} else {
			if err := os.Setenv(EnvDelayMs, originalDelay); err != nil {
				t.Logf("Failed to restore environment variable: %v", err)
			}
		}
	}()

	_, err := LoadFromEnv()
	if err == nil {
		t.Error("Expected error for invalid delay")
	}
}

func TestLoadFromFileSuccess(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	testConfig := models.Config{
		TelegramToken:     "file_token_123",
		GoogleCreds:       "test-creds.json",
		SheetID:           "test_sheet_id",
		DelayMs:           500,
		StartQuestionID:   "test_start",
		QuestionsFilePath: "configs/test-questions.json",
	}

	configData, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configPath, configData, 0o600)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Load configuration
	config, err := LoadFromFile(configPath)
	if err != nil {
		t.Errorf(expectedNoErrorMsg, err)
	}

	if config.TelegramToken != "file_token_123" {
		t.Errorf("Expected token 'file_token_123', got %s", config.TelegramToken)
	}
	if config.DelayMs != 500 {
		t.Errorf("Expected delay 500, got %d", config.DelayMs)
	}
}

func TestLoadFromFileInvalidPath(t *testing.T) {
	_, err := LoadFromFile("")
	if err == nil {
		t.Error("Expected error for empty path")
	}

	_, err = LoadFromFile("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestLoadQuestionsSuccess(t *testing.T) {
	// Create temporary questions file
	tmpDir := t.TempDir()
	questionsPath := filepath.Join(tmpDir, "questions.json")

	testQuestions := []models.Question{
		{
			ID:   "start",
			Text: "Start question",
			Options: []models.Option{
				{Text: "Continue", NextID: "next"},
			},
		},
		{
			ID:   "next",
			Text: "Next question",
		},
	}

	questionsData, err := json.Marshal(testQuestions)
	if err != nil {
		t.Fatalf("Failed to marshal test questions: %v", err)
	}

	err = os.WriteFile(questionsPath, questionsData, 0o600)
	if err != nil {
		t.Fatalf("Failed to write test questions file: %v", err)
	}

	// Load questions
	questions, err := LoadQuestions(questionsPath)
	if err != nil {
		t.Errorf(expectedNoErrorMsg, err)
	}

	if len(questions) != 2 {
		t.Errorf("Expected 2 questions, got %d", len(questions))
	}

	if _, exists := questions["start"]; !exists {
		t.Error("Expected 'start' question to exist")
	}
	if _, exists := questions["next"]; !exists {
		t.Error("Expected 'next' question to exist")
	}
}

func TestLoadQuestionsInvalidCases(t *testing.T) {
	// Test with empty path
	_, err := LoadQuestions("")
	if err == nil {
		t.Error("Expected error for empty path")
	}

	// Test with non-existent file
	_, err = LoadQuestions("/nonexistent/questions.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}

	// Test with duplicate IDs
	tmpDir := t.TempDir()
	questionsPath := filepath.Join(tmpDir, "duplicate_questions.json")

	duplicateQuestions := []models.Question{
		{ID: "duplicate", Text: "First question"},
		{ID: "duplicate", Text: "Second question"},
	}

	questionsData, _ := json.Marshal(duplicateQuestions)
	if err := os.WriteFile(questionsPath, questionsData, 0o600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err = LoadQuestions(questionsPath)
	if err == nil {
		t.Error("Expected error for duplicate question IDs")
	}
}

func TestGetDelayFromEnv(t *testing.T) {
	originalDelay := os.Getenv(EnvDelayMs)
	defer restoreEnvVar(EnvDelayMs, originalDelay)

	tests := []struct {
		name        string
		envValue    string
		expected    int
		expectError bool
	}{
		{"default value", "", DefaultDelayMs, false},
		{"valid positive", "1000", 1000, false},
		{"valid zero", "0", 0, false},
		{"invalid string", "not_a_number", 0, true},
		{"negative value", "-100", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnvVar(EnvDelayMs, tt.envValue)
			result, err := getDelayFromEnv()
			assertDelayResult(t, result, err, tt.expected, tt.expectError)
		})
	}
}

func setEnvVar(key, value string) {
	var err error
	if value == "" {
		err = os.Unsetenv(key)
	} else {
		err = os.Setenv(key, value)
	}
	if err != nil {
		// In test context, we might want to log this but not fail
		// since this is a helper function used in table tests
		panic(fmt.Sprintf("Failed to set environment variable %s: %v", key, err))
	}
}

func restoreEnvVar(key, originalValue string) {
	var err error
	if originalValue == "" {
		err = os.Unsetenv(key)
	} else {
		err = os.Setenv(key, originalValue)
	}
	if err != nil {
		// In cleanup context, we might want to log this but not fail
		panic(fmt.Sprintf("Failed to restore environment variable %s: %v", key, err))
	}
}

func assertDelayResult(t *testing.T, result int, err error, expected int, expectError bool) {
	if expectError && err == nil {
		t.Error("Expected error but got none")
		return
	}
	if !expectError && err != nil {
		t.Errorf("Expected no error but got: %v", err)
		return
	}
	if !expectError && result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
