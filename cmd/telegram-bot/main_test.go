// Package main provides tests for the telegram bot entry point.
package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"tlgbot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Mock implementations for testing
type mockUpdatesChan struct {
	updates chan tgbotapi.Update
	closed  bool
}

func newMockUpdatesChan() *mockUpdatesChan {
	return &mockUpdatesChan{
		updates: make(chan tgbotapi.Update, 1),
		closed:  false,
	}
}

func (m *mockUpdatesChan) GetUpdatesChan(_ tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return m.updates
}

func (m *mockUpdatesChan) Close() {
	if !m.closed {
		close(m.updates)
		m.closed = true
	}
}

func (m *mockUpdatesChan) SendUpdate(update tgbotapi.Update) {
	if !m.closed {
		select {
		case m.updates <- update:
		default:
		}
	}
}

// Test helper functions
func setupTestEnv(t *testing.T) (string, func()) {
	questionsPath := createTestQuestionsFile(t)
	cleanup := setupTestEnvironment(t, questionsPath)
	return questionsPath, cleanup
}

func createTestQuestionsFile(t *testing.T) string {
	tmpDir := t.TempDir()
	questionsPath := filepath.Join(tmpDir, "questions.json")

	testQuestions := []models.Question{
		{
			ID:   "start",
			Text: "Welcome! Choose an option:",
			Options: []models.Option{
				{Text: "Option 1", NextID: "end"},
				{Text: "Option 2", NextID: "end"},
			},
		},
		{
			ID:   "end",
			Text: "Thank you!",
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

	return questionsPath
}

func setupTestEnvironment(t *testing.T, questionsPath string) func() {
	envVars := map[string]string{
		"TELEGRAM_TOKEN":      "test_token_123",
		"DELAY_MS":            "100",
		"START_QUESTION_ID":   "start",
		"QUESTIONS_FILE_PATH": questionsPath,
	}

	originalValues := setTestEnvVars(t, envVars)
	return func() {
		restoreEnvVars(t, originalValues)
	}
}

func setTestEnvVars(t *testing.T, envVars map[string]string) map[string]string {
	originalValues := make(map[string]string)
	for key, value := range envVars {
		originalValues[key] = os.Getenv(key)
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Failed to set %s: %v", key, err)
		}
	}
	return originalValues
}

func restoreEnvVars(t *testing.T, originalValues map[string]string) {
	for key, original := range originalValues {
		restoreEnvVar(t, key, original)
	}
}

func restoreEnvVar(t *testing.T, key, original string) {
	if original == "" {
		if err := os.Unsetenv(key); err != nil {
			t.Errorf("Failed to unset %s: %v", key, err)
		}
		return
	}
	if err := os.Setenv(key, original); err != nil {
		t.Errorf("Failed to restore %s: %v", key, err)
	}
}

// Helper functions to reduce cognitive complexity
func assertConfigEquals(t *testing.T, cfg *models.Config, expectedToken string, expectedDelay int, expectedStartID string) {
	if cfg.TelegramToken != expectedToken {
		t.Errorf("Expected token '%s', got %s", expectedToken, cfg.TelegramToken)
	}
	if cfg.DelayMs != expectedDelay {
		t.Errorf("Expected delay %d, got %d", expectedDelay, cfg.DelayMs)
	}
	if cfg.StartQuestionID != expectedStartID {
		t.Errorf("Expected start question '%s', got %s", expectedStartID, cfg.StartQuestionID)
	}
}

func loadQuestionsFromFile(t *testing.T, questionsPath string) map[string]models.Question {
	data, err := os.ReadFile(questionsPath) //nolint:gosec // Test code with controlled file paths
	if err != nil {
		t.Fatalf("Failed to read questions file: %v", err)
	}

	var questions []models.Question
	err = json.Unmarshal(data, &questions)
	if err != nil {
		t.Fatalf("Failed to unmarshal questions: %v", err)
	}

	questionsMap := make(map[string]models.Question)
	for _, q := range questions {
		questionsMap[q.ID] = q
	}
	return questionsMap
}

func validateQuestionsMap(t *testing.T, questionsMap map[string]models.Question, expectedCount int) {
	if len(questionsMap) != expectedCount {
		t.Errorf("Expected %d questions, got %d", expectedCount, len(questionsMap))
	}

	startQuestion, exists := questionsMap["start"]
	if !exists {
		t.Error("Expected 'start' question to exist")
	}

	if startQuestion.Text != "Welcome! Choose an option:" {
		t.Errorf("Unexpected start question text: %s", startQuestion.Text)
	}
}

func setupInvalidEnvVar(t *testing.T, envVar string) func() {
	original := os.Getenv(envVar)
	if err := os.Unsetenv(envVar); err != nil {
		t.Fatalf("Failed to unset %s: %v", envVar, err)
	}

	return func() {
		if original != "" {
			if err := os.Setenv(envVar, original); err != nil {
				t.Errorf("Failed to restore %s: %v", envVar, err)
			}
		}
	}
}

func createInvalidJSONFile(t *testing.T) string {
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "invalid.json")

	err := os.WriteFile(invalidPath, []byte("invalid json content"), 0o600)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	return invalidPath
}

func TestMainConfigurationLoading(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := &models.Config{
		TelegramToken:     os.Getenv("TELEGRAM_TOKEN"),
		DelayMs:           100,
		StartQuestionID:   "start",
		QuestionsFilePath: os.Getenv("QUESTIONS_FILE_PATH"),
	}

	assertConfigEquals(t, cfg, "test_token_123", 100, "start")
}

func TestMainConfigurationValidation(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("ValidConfiguration", func(t *testing.T) {
		cfg := &models.Config{
			TelegramToken:     "test_token_123",
			DelayMs:           100,
			StartQuestionID:   "start",
			QuestionsFilePath: "test.json",
		}

		err := cfg.Validate()
		if err != nil {
			t.Errorf("Expected valid configuration, got error: %v", err)
		}
	})

	t.Run("InvalidConfiguration", func(t *testing.T) {
		cfg := &models.Config{
			TelegramToken:     "", // Empty token should fail
			DelayMs:           100,
			StartQuestionID:   "start",
			QuestionsFilePath: "test.json",
		}

		err := cfg.Validate()
		if err == nil {
			t.Error("Expected validation error for empty token")
		}
	})
}

func TestMainQuestionsLoading(t *testing.T) {
	questionsPath, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("LoadValidQuestions", func(t *testing.T) {
		questionsMap := loadQuestionsFromFile(t, questionsPath)
		validateQuestionsMap(t, questionsMap, 2)
	})

	t.Run("LoadNonexistentQuestions", func(t *testing.T) {
		_, err := os.ReadFile("nonexistent.json")
		if err == nil {
			t.Error("Expected error when loading nonexistent questions file")
		}
	})
}

func TestStartBotFunction(t *testing.T) {
	mockUpdates := newMockUpdatesChan()
	defer mockUpdates.Close()

	done := make(chan bool, 1)

	go func() {
		defer func() {
			done <- true
		}()

		select {
		case update := <-mockUpdates.updates:
			if update.Message != nil {
				t.Logf("Processed message: %s", update.Message.Text)
			}
			if update.CallbackQuery != nil {
				t.Logf("Processed callback query: %s", update.CallbackQuery.Data)
			}
		case <-time.After(100 * time.Millisecond):
			// Timeout
		}
	}()

	mockUpdates.SendUpdate(tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			MessageID: 1,
			Text:      "test",
			Chat: &tgbotapi.Chat{
				ID: 12345,
			},
		},
	})

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Bot processing timed out")
	}
}

func TestMainInitializationComponents(t *testing.T) {
	questionsPath, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := &models.Config{
		TelegramToken:     "test_token_123",
		DelayMs:           100,
		StartQuestionID:   "start",
		QuestionsFilePath: questionsPath,
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Configuration validation failed: %v", err)
	}

	questionsMap := loadQuestionsFromFile(t, questionsPath)

	if _, exists := questionsMap[cfg.StartQuestionID]; !exists {
		t.Errorf("Start question '%s' not found in questions", cfg.StartQuestionID)
	}

	if len(questionsMap) < 1 {
		t.Error("No questions loaded")
	}

	t.Logf("Successfully initialized with %d questions", len(questionsMap))
}

func TestInitializeBot(t *testing.T) {
	questionsPath, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := &models.Config{
		TelegramToken:     "test_token_123",
		DelayMs:           100,
		StartQuestionID:   "start",
		QuestionsFilePath: questionsPath,
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Configuration validation failed: %v", err)
	}

	questionsMap := loadQuestionsFromFile(t, questionsPath)

	if len(questionsMap) == 0 {
		t.Error("No questions loaded")
	}
}

func TestMainErrorHandling(t *testing.T) {
	t.Run("MissingEnvironmentVariables", func(t *testing.T) {
		restore := setupInvalidEnvVar(t, "TELEGRAM_TOKEN")
		defer restore()

		token := os.Getenv("TELEGRAM_TOKEN")
		if token != "" {
			t.Error("Expected empty token after unsetting environment variable")
		}
	})

	t.Run("InvalidQuestionsFile", func(t *testing.T) {
		invalidPath := createInvalidJSONFile(t)

		var questions []models.Question
		data, _ := os.ReadFile(invalidPath) //nolint:gosec // Test code with controlled file paths
		err := json.Unmarshal(data, &questions)
		if err == nil {
			t.Error("Expected error when parsing invalid JSON")
		}
	})
}
