package handlers

import (
	"testing"
	"tlgbot/internal/models"
	"tlgbot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Constants to avoid goconst warnings
const (
	testUserName    = "John"
	startQuestionID = "start"
	userID          = int64(123)
)

// mockTelegramBot is a mock implementation of the BotService interface
type mockTelegramBot struct {
	processQuestionCalled     bool
	processAnswerCalled       bool
	processOptionAnswerCalled bool
	handleAutoAdvanceCalled   bool
	sendMessageCalled         bool
	sendImagesCalled          bool
	sendMessagesCalled        bool
	lastUserID                int64
	lastMessage               string
	lastAnswer                string
	lastOption                string
	lastQuestion              *models.Question

	// Mock API for callback acknowledgment
	api *tgbotapi.BotAPI
}

func (m *mockTelegramBot) SendImages(userID int64, _ []string, _ int) error {
	m.sendImagesCalled = true
	m.lastUserID = userID
	return nil
}

func (m *mockTelegramBot) SendMessage(userID int64, text string, _ interface{}) error {
	m.sendMessageCalled = true
	m.lastUserID = userID
	m.lastMessage = text
	return nil
}

func (m *mockTelegramBot) SendMessages(userID int64, _ []string, _ string, _ interface{}) error {
	m.sendMessagesCalled = true
	m.lastUserID = userID
	return nil
}

func (m *mockTelegramBot) ProcessQuestion(userID int64, question *models.Question) error {
	m.processQuestionCalled = true
	m.lastUserID = userID
	m.lastQuestion = question
	return nil
}

func (m *mockTelegramBot) ProcessAnswer(userID int64, answer string) error {
	m.processAnswerCalled = true
	m.lastUserID = userID
	m.lastAnswer = answer
	return nil
}

func (m *mockTelegramBot) ProcessOptionAnswer(userID int64, optionText string) error {
	m.processOptionAnswerCalled = true
	m.lastUserID = userID
	m.lastOption = optionText
	return nil
}

func (m *mockTelegramBot) HandleAutoAdvance(userID int64, question *models.Question) error {
	m.handleAutoAdvanceCalled = true
	m.lastUserID = userID
	m.lastQuestion = question
	return nil
}

func (m *mockTelegramBot) GetAPI() *tgbotapi.BotAPI {
	if m.api == nil {
		m.api = &tgbotapi.BotAPI{}
	}
	return m.api
}

func createTestHandler(_ *testing.T) (*TelegramHandler, *mockTelegramBot, *services.UserStateManager, *services.QuestionManager) {
	// Create test config
	config := &models.Config{
		TelegramToken:     "test_token",
		DelayMs:           100,
		StartQuestionID:   startQuestionID,
		QuestionsFilePath: "test.json",
	}

	// Create services
	userStateManager := services.NewUserStateManager()

	questions := map[string]models.Question{
		startQuestionID: {
			ID:   startQuestionID,
			Text: "Welcome! What's your name?",
			Options: []models.Option{
				{Text: "Continue", NextID: "question1"},
			},
		},
		"question1": {
			ID:   "question1",
			Text: "How are you today?",
			Options: []models.Option{
				{Text: "Good", NextID: "end"},
				{Text: "Bad", NextID: "end"},
			},
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
			Text: "Thank you!",
		},
	}

	questionManager := services.NewQuestionManager(questions)

	// Create mock bot
	mockBot := &mockTelegramBot{}

	// Create handler
	handler := NewTelegramHandler(mockBot, config, userStateManager, questionManager)

	return handler, mockBot, userStateManager, questionManager
}

func TestNewTelegramHandler(t *testing.T) {
	mockBot := &mockTelegramBot{}
	config := &models.Config{}
	userStateManager := services.NewUserStateManager()
	questionManager := services.NewQuestionManager(map[string]models.Question{})

	handler := NewTelegramHandler(mockBot, config, userStateManager, questionManager)

	if handler == nil {
		t.Fatal("Expected handler to be created")
	}
	if handler.config != config {
		t.Error("Expected config to be set correctly")
	}
}

func TestHandleMessage(t *testing.T) {
	handler, mockBot, userStateManager, _ := createTestHandler(t)

	tests := []struct {
		name        string
		message     *tgbotapi.Message
		expectStart bool
		expectText  bool
	}{
		{
			name: "start command",
			message: &tgbotapi.Message{
				From: &tgbotapi.User{
					ID:        userID,
					FirstName: testUserName,
				},
				Text: "/start",
				Entities: []tgbotapi.MessageEntity{
					{Type: "bot_command", Offset: 0, Length: 6},
				},
			},
			expectStart: true,
			expectText:  false,
		},
		{
			name: "regular text message",
			message: &tgbotapi.Message{
				From: &tgbotapi.User{
					ID:        userID,
					FirstName: testUserName,
				},
				Text: "Hello World",
			},
			expectStart: false,
			expectText:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockBot.processQuestionCalled = false
			mockBot.processAnswerCalled = false

			handler.HandleMessage(tt.message)

			if tt.expectStart && !mockBot.processQuestionCalled {
				t.Error("Expected ProcessQuestion to be called for start command")
			}

			// Check that user state was created
			userState := userStateManager.GetUserState(tt.message.From.ID)
			if userState == nil {
				t.Error("Expected user state to be created")
			}
		})
	}
}

func TestHandleCallbackQuery(t *testing.T) {
	handler, mockBot, userStateManager, _ := createTestHandler(t)

	// Create user state first
	userState := userStateManager.GetOrCreateUserState(userID, testUserName)
	userState.CurrentQuestionID = "question1"

	tests := []struct {
		name         string
		callback     *tgbotapi.CallbackQuery
		expectOption bool
	}{
		{
			name: "valid callback",
			callback: &tgbotapi.CallbackQuery{
				ID: "callback_123",
				From: &tgbotapi.User{
					ID:        userID,
					FirstName: testUserName,
				},
				Data: "Good",
			},
			expectOption: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockBot.processOptionAnswerCalled = false

			handler.HandleCallbackQuery(tt.callback)

			if tt.expectOption && !mockBot.processOptionAnswerCalled {
				t.Error("Expected ProcessOptionAnswer to be called")
			}

			if tt.expectOption && mockBot.lastOption != tt.callback.Data {
				t.Errorf("Expected option %s, got %s", tt.callback.Data, mockBot.lastOption)
			}
		})
	}
}

func TestHandleCommand(t *testing.T) {
	handler, mockBot, userStateManager, _ := createTestHandler(t)

	userState := userStateManager.GetOrCreateUserState(userID, testUserName)

	tests := []struct {
		name        string
		command     string
		expectStart bool
	}{
		{
			name:        "start command",
			command:     "start",
			expectStart: true,
		},
		{
			name:        "unknown command",
			command:     "unknown",
			expectStart: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockBot.processQuestionCalled = false
			mockBot.handleAutoAdvanceCalled = false

			message := &tgbotapi.Message{
				From: &tgbotapi.User{
					ID:        userID,
					FirstName: testUserName,
				},
				Text: "/" + tt.command,
				Entities: []tgbotapi.MessageEntity{
					{Type: "bot_command", Offset: 0, Length: len(tt.command) + 1},
				},
			}

			handler.handleCommand(message, userState)

			if tt.expectStart {
				if !mockBot.processQuestionCalled {
					t.Error("Expected ProcessQuestion to be called")
				}
				if !mockBot.handleAutoAdvanceCalled {
					t.Error("Expected HandleAutoAdvance to be called")
				}
				if userState.CurrentQuestionID != startQuestionID {
					t.Errorf("Expected current question to be '%s', got %s", startQuestionID, userState.CurrentQuestionID)
				}
			}
		})
	}
}

func TestHandleTextInput(t *testing.T) {
	handler, mockBot, userStateManager, _ := createTestHandler(t)

	userState := userStateManager.GetOrCreateUserState(userID, testUserName)

	tests := []struct {
		name              string
		currentQuestionID string
		messageText       string
		expectProcessed   bool
	}{
		{
			name:              "text input for input question",
			currentQuestionID: "input_question",
			messageText:       "My name is John",
			expectProcessed:   true,
		},
		{
			name:              "text input for non-input question",
			currentQuestionID: "question1",
			messageText:       "Some text",
			expectProcessed:   false,
		},
		{
			name:              "no current question",
			currentQuestionID: "",
			messageText:       "Some text",
			expectProcessed:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock and state
			mockBot.processAnswerCalled = false
			userState.CurrentQuestionID = tt.currentQuestionID

			message := &tgbotapi.Message{
				From: &tgbotapi.User{
					ID:        userID,
					FirstName: testUserName,
				},
				Text: tt.messageText,
			}

			handler.handleTextInput(message, userState)

			if tt.expectProcessed {
				if !mockBot.processAnswerCalled {
					t.Error("Expected ProcessAnswer to be called")
				}
				if mockBot.lastAnswer != tt.messageText {
					t.Errorf("Expected answer %s, got %s", tt.messageText, mockBot.lastAnswer)
				}
			} else {
				if mockBot.processAnswerCalled {
					t.Error("Expected ProcessAnswer not to be called")
				}
			}
		})
	}
}

func TestStartConversation(t *testing.T) {
	handler, mockBot, userStateManager, _ := createTestHandler(t)

	userState := userStateManager.GetOrCreateUserState(userID, testUserName)

	// Reset mock
	mockBot.processQuestionCalled = false
	mockBot.handleAutoAdvanceCalled = false

	handler.startConversation(userID, userState)

	// Check that ProcessQuestion was called
	if !mockBot.processQuestionCalled {
		t.Error("Expected ProcessQuestion to be called")
	}

	// Check that HandleAutoAdvance was called
	if !mockBot.handleAutoAdvanceCalled {
		t.Error("Expected HandleAutoAdvance to be called")
	}

	// Check that current question ID was set
	if userState.CurrentQuestionID != startQuestionID {
		t.Errorf("Expected current question ID to be '%s', got %s", startQuestionID, userState.CurrentQuestionID)
	}

	// Check that user state was updated
	retrievedState := userStateManager.GetUserState(userID)
	if retrievedState == nil {
		t.Error("Expected user state to be saved")
		return
	}
	if retrievedState.CurrentQuestionID != startQuestionID {
		t.Errorf("Expected saved current question ID to be '%s', got %s", startQuestionID, retrievedState.CurrentQuestionID)
	}
}

func TestMoveToNextQuestion(t *testing.T) {
	handler, mockBot, userStateManager, _ := createTestHandler(t)

	userState := userStateManager.GetOrCreateUserState(userID, testUserName)
	userState.CurrentQuestionID = startQuestionID

	tests := []struct {
		name           string
		nextQuestionID string
		expectError    bool
	}{
		{
			name:           "valid next question",
			nextQuestionID: "question1",
			expectError:    false,
		},
		{
			name:           "invalid next question",
			nextQuestionID: "nonexistent",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockBot.processQuestionCalled = false
			mockBot.handleAutoAdvanceCalled = false

			err := handler.moveToNextQuestion(userID, tt.nextQuestionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error for invalid question ID")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}

				if !mockBot.processQuestionCalled {
					t.Error("Expected ProcessQuestion to be called")
				}

				if !mockBot.handleAutoAdvanceCalled {
					t.Error("Expected HandleAutoAdvance to be called")
				}

				// Check that current question was updated
				updatedState := userStateManager.GetUserState(userID)
				if updatedState.CurrentQuestionID != tt.nextQuestionID {
					t.Errorf("Expected current question ID to be %s, got %s", tt.nextQuestionID, updatedState.CurrentQuestionID)
				}
			}
		})
	}
}

func TestHandleCallbackQueryWithoutUserState(t *testing.T) {
	handler, mockBot, _, _ := createTestHandler(t)

	// Create callback for user without state
	callback := &tgbotapi.CallbackQuery{
		ID: "callback_123",
		From: &tgbotapi.User{
			ID:        999, // Non-existent user
			FirstName: "Unknown",
		},
		Data: "Some option",
	}

	// Reset mock
	mockBot.processOptionAnswerCalled = false

	handler.HandleCallbackQuery(callback)

	// Should not process option answer when user state doesn't exist
	if mockBot.processOptionAnswerCalled {
		t.Error("Expected ProcessOptionAnswer not to be called for user without state")
	}
}

// Test message command detection
func TestMessageIsCommand(t *testing.T) {
	tests := []struct {
		name      string
		message   *tgbotapi.Message
		isCommand bool
		command   string
	}{
		{
			name: "start command",
			message: &tgbotapi.Message{
				Text: "/start",
				Entities: []tgbotapi.MessageEntity{
					{Type: "bot_command", Offset: 0, Length: 6},
				},
			},
			isCommand: true,
			command:   "start",
		},
		{
			name: "help command",
			message: &tgbotapi.Message{
				Text: "/help",
				Entities: []tgbotapi.MessageEntity{
					{Type: "bot_command", Offset: 0, Length: 5},
				},
			},
			isCommand: true,
			command:   "help",
		},
		{
			name: "regular text",
			message: &tgbotapi.Message{
				Text: "Hello world",
			},
			isCommand: false,
			command:   "",
		},
		{
			name: "text with slash but no command entity",
			message: &tgbotapi.Message{
				Text: "/not_a_command",
			},
			isCommand: false,
			command:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCommand := tt.message.IsCommand()
			if isCommand != tt.isCommand {
				t.Errorf("Expected IsCommand() to return %v, got %v", tt.isCommand, isCommand)
			}

			if tt.isCommand {
				command := tt.message.Command()
				if command != tt.command {
					t.Errorf("Expected command %s, got %s", tt.command, command)
				}
			}
		})
	}
}
