package models

import (
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		expectErr bool
	}{
		{
			name: "valid config",
			config: Config{
				TelegramToken:   "valid_token",
				StartQuestionID: "start",
				DelayMs:         100,
			},
			expectErr: false,
		},
		{
			name: "missing telegram token",
			config: Config{
				StartQuestionID: "start",
				DelayMs:         100,
			},
			expectErr: true,
		},
		{
			name: "missing start question ID",
			config: Config{
				TelegramToken: "valid_token",
				DelayMs:       100,
			},
			expectErr: true,
		},
		{
			name: "negative delay",
			config: Config{
				TelegramToken:   "valid_token",
				StartQuestionID: "start",
				DelayMs:         -100,
			},
			expectErr: true,
		},
		{
			name: "zero delay is valid",
			config: Config{
				TelegramToken:   "valid_token",
				StartQuestionID: "start",
				DelayMs:         0,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("Validate() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestQuestionGetDelayMs(t *testing.T) {
	defaultDelay := 1000

	tests := []struct {
		name     string
		question Question
		expected int
	}{
		{
			name:     "question with nil delay",
			question: Question{DelayMs: nil},
			expected: defaultDelay,
		},
		{
			name:     "question with specific delay",
			question: Question{DelayMs: intPtr(500)},
			expected: 500,
		},
		{
			name:     "question with zero delay",
			question: Question{DelayMs: intPtr(0)},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.question.GetDelayMs(defaultDelay)
			if result != tt.expected {
				t.Errorf("GetDelayMs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuestionHasKeyboard(t *testing.T) {
	tests := []struct {
		name     string
		question Question
		expected bool
	}{
		{
			name: "auto advance question",
			question: Question{
				AutoAdvance: true,
				Options: []Option{
					{Text: "Option 1"},
				},
			},
			expected: false,
		},
		{
			name: "question with options",
			question: Question{
				AutoAdvance: false,
				Options: []Option{
					{Text: "Option 1"},
					{Text: "Option 2"},
				},
			},
			expected: true,
		},
		{
			name: "question with external link",
			question: Question{
				AutoAdvance:  false,
				ExternalLink: "https://example.com",
				ExternalText: "Visit",
			},
			expected: true,
		},
		{
			name: "question with external link but no text",
			question: Question{
				AutoAdvance:  false,
				ExternalLink: "https://example.com",
				ExternalText: "",
			},
			expected: false,
		},
		{
			name: "question without options or external link",
			question: Question{
				AutoAdvance: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.question.HasKeyboard()
			if result != tt.expected {
				t.Errorf("HasKeyboard() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuestionGetDisplayText(t *testing.T) {
	tests := []struct {
		name     string
		question Question
		expected string
	}{
		{
			name: "question with text",
			question: Question{
				Text:     "Main question text",
				Messages: []string{"Message 1", "Message 2"},
			},
			expected: "Main question text",
		},
		{
			name: "question without text but with messages",
			question: Question{
				Text:     "",
				Messages: []string{"Message 1", "Message 2"},
			},
			expected: "Message 2",
		},
		{
			name: "question with single message",
			question: Question{
				Text:     "",
				Messages: []string{"Single message"},
			},
			expected: "Single message",
		},
		{
			name: "question without text or messages",
			question: Question{
				Text:     "",
				Messages: []string{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.question.GetDisplayText()
			if result != tt.expected {
				t.Errorf("GetDisplayText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

const testUserName = "TestUser"

func TestUserStateNewUserState(t *testing.T) {
	name := testUserName
	state := NewUserState(name)

	if state.Name != name {
		t.Errorf("Expected name %s, got %s", name, state.Name)
	}

	if state.Answers == nil {
		t.Error("Expected Answers map to be initialized")
	}

	if len(state.Answers) != 0 {
		t.Errorf("Expected empty Answers map, got %d items", len(state.Answers))
	}
}

func TestUserStateAddAndGetAnswer(t *testing.T) {
	state := NewUserState(testUserName)
	question := "test_question"
	answer := "test_answer"

	// Check that answer doesn't exist yet
	_, exists := state.GetAnswer(question)
	if exists {
		t.Error("Expected answer not to exist initially")
	}

	// Add answer
	state.AddAnswer(question, answer)

	// Check that answer was added
	retrievedAnswer, exists := state.GetAnswer(question)
	if !exists {
		t.Error("Expected answer to exist after adding")
	}
	if retrievedAnswer != answer {
		t.Errorf("Expected answer %s, got %s", answer, retrievedAnswer)
	}
}

func TestUserStateUpdateAnswer(t *testing.T) {
	state := NewUserState(testUserName)
	question := "test_question"
	firstAnswer := "first_answer"
	secondAnswer := "second_answer"

	// Add first answer
	state.AddAnswer(question, firstAnswer)

	// Update answer
	state.AddAnswer(question, secondAnswer)

	// Check that answer was updated
	retrievedAnswer, exists := state.GetAnswer(question)
	if !exists {
		t.Error("Expected answer to exist")
	}
	if retrievedAnswer != secondAnswer {
		t.Errorf("Expected updated answer %s, got %s", secondAnswer, retrievedAnswer)
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
