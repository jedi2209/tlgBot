package bot

import (
	"testing"
	"tlgbot/internal/models"
	"tlgbot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// mockBotAPI is a mock implementation of the Telegram Bot API
type mockBotAPI struct {
	sentMessages []tgbotapi.Chattable
	sendError    error
}

func (m *mockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if m.sendError != nil {
		return tgbotapi.Message{}, m.sendError
	}
	m.sentMessages = append(m.sentMessages, c)
	return tgbotapi.Message{}, nil
}

func (m *mockBotAPI) SendMediaGroup(config tgbotapi.MediaGroupConfig) ([]tgbotapi.Message, error) {
	if m.sendError != nil {
		return nil, m.sendError
	}
	m.sentMessages = append(m.sentMessages, config)
	return []tgbotapi.Message{}, nil
}

func (m *mockBotAPI) Request(_ tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return &tgbotapi.APIResponse{}, nil
}

func createTestBot(_ *testing.T) (*TelegramBot, *mockBotAPI, *services.UserStateManager, *services.QuestionManager) {
	// Create mock API
	mockAPI := &mockBotAPI{}

	// Create test config
	config := &models.Config{
		TelegramToken:     "test_token",
		DelayMs:           100,
		StartQuestionID:   "start",
		QuestionsFilePath: "test.json",
	}

	// Create services
	userStateManager := services.NewUserStateManager()

	questions := map[string]models.Question{
		"start": {
			ID:   "start",
			Text: "Welcome {name}! Choose an option:",
			Options: []models.Option{
				{Text: "Option 1", NextID: "question1"},
				{Text: "Option 2", NextID: "question2"},
			},
		},
		"question1": {
			ID:   "question1",
			Text: "You chose option 1",
		},
		"question2": {
			ID:      "question2",
			Text:    "You chose option 2",
			DelayMs: intPtr(500),
		},
		"with_images": {
			ID:     "with_images",
			Text:   "Question with images",
			Images: []string{"image1.jpg", "image2.jpg"},
		},
		"auto_advance": {
			ID:          "auto_advance",
			Text:        "Auto advancing question",
			AutoAdvance: true,
			DelayMs:     intPtr(200),
			Options: []models.Option{
				{Text: "Continue", NextID: "end"},
			},
		},
		"with_external": {
			ID:           "with_external",
			Text:         "Question with external link",
			ExternalLink: "https://example.com",
			ExternalText: "Visit Example",
		},
		"input_question": {
			ID:        "input_question",
			Text:      "Please enter your name:",
			InputType: "text",
			Options: []models.Option{
				{NextID: "end"},
			},
		},
		"end": {
			ID:   "end",
			Text: "Thank you for your responses!",
		},
	}

	questionManager := services.NewQuestionManager(questions)

	// Create bot with mock API (we can't use real tgbotapi.BotAPI easily)
	bot := &TelegramBot{
		api:              nil, // We'll mock API calls directly
		config:           config,
		userStateManager: userStateManager,
		questionManager:  questionManager,
	}

	return bot, mockAPI, userStateManager, questionManager
}

func TestNewTelegramBot(t *testing.T) {
	api := &tgbotapi.BotAPI{}
	config := &models.Config{}
	userStateManager := services.NewUserStateManager()
	questionManager := services.NewQuestionManager(map[string]models.Question{})

	bot := NewTelegramBot(api, config, userStateManager, questionManager)

	if bot == nil {
		t.Fatal("Expected bot to be created")
	}
	if bot.api != api {
		t.Error("Expected API to be set correctly")
	}
	if bot.config != config {
		t.Error("Expected config to be set correctly")
	}
}

func TestBuildKeyboard(t *testing.T) {
	bot, _, _, _ := createTestBot(t)

	tests := []struct {
		name     string
		question models.Question
		wantRows int
	}{
		{
			name: "question with options",
			question: models.Question{
				ID: "test",
				Options: []models.Option{
					{Text: "Option 1", NextID: "next1"},
					{Text: "Option 2", NextID: "next2"},
				},
			},
			wantRows: 2,
		},
		{
			name: "question with external link",
			question: models.Question{
				ID:           "test",
				ExternalLink: "https://example.com",
				ExternalText: "Visit",
			},
			wantRows: 1,
		},
		{
			name: "question with options and external link",
			question: models.Question{
				ID: "test",
				Options: []models.Option{
					{Text: "Option 1", NextID: "next1"},
				},
				ExternalLink: "https://example.com",
				ExternalText: "Visit",
			},
			wantRows: 2,
		},
		{
			name: "auto advance question",
			question: models.Question{
				ID:          "test",
				AutoAdvance: true,
				Options: []models.Option{
					{Text: "Option 1", NextID: "next1"},
				},
			},
			wantRows: 0, // Should return nil for auto advance
		},
		{
			name: "question without keyboard",
			question: models.Question{
				ID: "test",
			},
			wantRows: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyboard := bot.BuildKeyboard(&tt.question)

			if tt.wantRows == 0 {
				if keyboard != nil {
					t.Errorf("Expected nil keyboard, got keyboard with %d rows", len(keyboard.InlineKeyboard))
				}
				return
			}

			if keyboard == nil {
				t.Fatal("Expected keyboard, got nil")
			}

			if len(keyboard.InlineKeyboard) != tt.wantRows {
				t.Errorf("Expected %d rows, got %d", tt.wantRows, len(keyboard.InlineKeyboard))
			}
		})
	}
}

func TestReplaceNamePlaceholder(t *testing.T) {
	bot, _, _, _ := createTestBot(t)

	tests := []struct {
		name     string
		text     string
		userName string
		expected string
	}{
		{
			name:     "text with name placeholder",
			text:     "Hello {name}, welcome!",
			userName: "John",
			expected: "Hello John, welcome!",
		},
		{
			name:     "text without placeholder",
			text:     "Hello, welcome!",
			userName: "John",
			expected: "Hello, welcome!",
		},
		{
			name:     "text with multiple placeholders",
			text:     "Hello {name}, your name is {name}",
			userName: "John",
			expected: "Hello John, your name is John",
		},
		{
			name:     "empty username",
			text:     "Hello {name}!",
			userName: "",
			expected: "Hello !",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bot.replaceNamePlaceholder(tt.text, tt.userName)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGenerateAnswersSummary(t *testing.T) {
	bot, _, _, _ := createTestBot(t)

	tests := []struct {
		name     string
		answers  map[string]string
		expected string
	}{
		{
			name:     "empty answers",
			answers:  map[string]string{},
			expected: "",
		},
		{
			name: "single answer",
			answers: map[string]string{
				"question1": "answer1",
			},
			expected: "\n\n📋 Your answers:\n• question1: answer1\n",
		},
		{
			name: "multiple answers",
			answers: map[string]string{
				"question1": "answer1",
				"question2": "answer2",
			},
			// We'll check contains instead of exact match due to map ordering
			expected: "\n\n📋 Your answers:\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bot.generateAnswersSummary(tt.answers)
			if len(tt.answers) <= 1 {
				if result != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result)
				}
			} else {
				// For multiple answers, order might vary, so just check that it contains expected parts
				if !containsAllAnswers(result, tt.answers) {
					t.Errorf("Result %q doesn't contain all expected answers", result)
				}
			}
		})
	}
}

func TestBuildMediaGroup(t *testing.T) {
	bot, _, _, _ := createTestBot(t)

	tests := []struct {
		name     string
		images   []string
		expected int
	}{
		{
			name:     "empty images",
			images:   []string{},
			expected: 0,
		},
		{
			name:     "single image",
			images:   []string{"image1.jpg"},
			expected: 1,
		},
		{
			name:     "multiple images",
			images:   []string{"image1.jpg", "image2.jpg", "image3.jpg"},
			expected: 3,
		},
		{
			name:     "images with empty strings",
			images:   []string{"image1.jpg", "", "image2.jpg"},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bot.buildMediaGroup(tt.images)
			if len(result) != tt.expected {
				t.Errorf("Expected %d media items, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestGetTelegramName(t *testing.T) {
	tests := []struct {
		name     string
		user     *tgbotapi.User
		expected string
	}{
		{
			name: "user with first name only",
			user: &tgbotapi.User{
				FirstName: "John",
			},
			expected: "John",
		},
		{
			name: "user with first and last name", // Last name is ignored in implementation
			user: &tgbotapi.User{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John",
		},
		{
			name: "user with username", // Username is ignored if FirstName exists
			user: &tgbotapi.User{
				FirstName: "John",
				UserName:  "johndoe",
			},
			expected: "John",
		},
		{
			name: "user with all fields", // Only FirstName is used
			user: &tgbotapi.User{
				FirstName: "John",
				LastName:  "Doe",
				UserName:  "johndoe",
			},
			expected: "John",
		},
		{
			name: "user with only username",
			user: &tgbotapi.User{
				UserName: "johndoe",
			},
			expected: "johndoe",
		},
		{
			name:     "user with no name fields",
			user:     &tgbotapi.User{},
			expected: "User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTelegramName(tt.user)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Helper functions

func intPtr(i int) *int {
	return &i
}

func containsAllAnswers(summary string, answers map[string]string) bool {
	for question, answer := range answers {
		expectedLine := question + ": " + answer
		if !containsString(summary, expectedLine) {
			return false
		}
	}
	return containsString(summary, "📋 Your answers:")
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
