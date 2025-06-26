package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"tlgbot/internal/models"
)

// Constants for environment variables
const (
	EnvTelegramToken   = "TELEGRAM_TOKEN"
	EnvGoogleCreds     = "GOOGLE_CREDS"
	EnvSheetID         = "SHEET_ID"
	EnvDelayMs         = "DELAY_MS"
	EnvStartQuestionID = "START_QUESTION_ID"
)

// Default values
const (
	DefaultGoogleCreds     = "google-credentials.json"
	DefaultSheetID         = "YOUR_GOOGLE_SHEET_ID"
	DefaultDelayMs         = 700
	DefaultStartQuestionID = "start"
)

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*models.Config, error) {
	config := &models.Config{}

	var err error

	// Get required Telegram token
	config.TelegramToken = os.Getenv(EnvTelegramToken)
	if config.TelegramToken == "" {
		return nil, fmt.Errorf("%s environment variable is required", EnvTelegramToken)
	}

	// Get optional parameters with default values
	config.GoogleCreds = getEnvOrDefault(EnvGoogleCreds, DefaultGoogleCreds)
	config.SheetID = getEnvOrDefault(EnvSheetID, DefaultSheetID)
	config.StartQuestionID = getEnvOrDefault(EnvStartQuestionID, DefaultStartQuestionID)

	// Get delay with validation
	config.DelayMs, err = getDelayFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", EnvDelayMs, err)
	}

	return config, nil
}

// LoadFromFile loads configuration from file (fallback)
func LoadFromFile(path string) (*models.Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config file path cannot be empty")
	}

	file, err := os.Open(path) //nolint:gosec // G304: Config file path is controlled by application
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Failed to close config file: %v", closeErr)
		}
	}()

	config := &models.Config{}
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config from %s: %w", path, err)
	}

	return config, nil
}

// LoadQuestions loads questions from JSON file
func LoadQuestions(filename string) (map[string]models.Question, error) {
	if filename == "" {
		return nil, fmt.Errorf("questions file path cannot be empty")
	}

	data, err := os.ReadFile(filename) //nolint:gosec // G304: Questions file path is controlled by application
	if err != nil {
		return nil, fmt.Errorf("failed to read questions file %s: %w", filename, err)
	}

	var questions []models.Question
	if err := json.Unmarshal(data, &questions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal questions from %s: %w", filename, err)
	}

	if len(questions) == 0 {
		return nil, fmt.Errorf("no questions found in %s", filename)
	}

	// Convert to map and validate
	qMap := make(map[string]models.Question, len(questions))
	for i := range questions {
		q := &questions[i]
		if q.ID == "" {
			return nil, fmt.Errorf("question at index %d has empty ID", i)
		}

		if _, exists := qMap[q.ID]; exists {
			return nil, fmt.Errorf("duplicate question ID: %s", q.ID)
		}

		qMap[q.ID] = *q
	}

	return qMap, nil
}

// getEnvOrDefault returns environment variable value or default value
func getEnvOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}

// getDelayFromEnv gets and validates delay from environment variable
func getDelayFromEnv() (int, error) {
	delayStr := os.Getenv(EnvDelayMs)
	if delayStr == "" {
		return DefaultDelayMs, nil
	}

	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		return 0, fmt.Errorf("invalid delay value %s: %w", delayStr, err)
	}

	if delay < 0 {
		return 0, fmt.Errorf("delay cannot be negative: %d", delay)
	}

	return delay, nil
}
