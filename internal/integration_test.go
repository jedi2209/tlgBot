// Package internal provides integration tests for the Telegram bot application.
package internal

import (
	"testing"

	"tlgbot/internal/models"
	"tlgbot/internal/services"
)

const notGoodAnswer = "Not good"

func TestIntegrationQuestionFlow(t *testing.T) {
	// Create test questions
	questions := map[string]models.Question{
		"start": {
			ID:   "start",
			Text: "Welcome! What's your name?",
			Options: []models.Option{
				{Text: "Continue", NextID: "question_1"},
			},
		},
		"question_1": {
			ID:   "question_1",
			Text: "How are you feeling today?",
			Options: []models.Option{
				{Text: "Good", NextID: "end"},
				{Text: notGoodAnswer, NextID: "help"},
			},
		},
		"help": {
			ID:   "help",
			Text: "We're here to help!",
			Options: []models.Option{
				{Text: "Thanks", NextID: "end"},
			},
		},
		"end": {
			ID:   "end",
			Text: "Thank you for your time!",
		},
	}

	// Initialize services
	questionManager := services.NewQuestionManager(questions)
	userStateManager := services.NewUserStateManager()

	// Test complete flow
	userID := int64(12345)
	const testUserName = "TestUser"

	userName := testUserName

	// 1. Create user
	userState := userStateManager.GetOrCreateUserState(userID, userName)
	if userState.Name != userName {
		t.Errorf("Expected user name %s, got %s", userName, userState.Name)
	}

	// 2. Start with initial question
	userStateManager.UpdateCurrentQuestion(userID, "start")
	currentQuestion, err := questionManager.GetQuestion("start")
	if err != nil {
		t.Fatalf("Failed to get start question: %v", err)
	}

	if currentQuestion.Text != "Welcome! What's your name?" {
		t.Errorf("Unexpected start question text: %s", currentQuestion.Text)
	}

	// 3. Move to next question
	nextQuestion, err := questionManager.GetNextQuestion("start", "Continue")
	if err != nil {
		t.Fatalf("Failed to get next question: %v", err)
	}

	userStateManager.UpdateCurrentQuestion(userID, nextQuestion.ID)
	userState.AddAnswer("start", "Continue")

	// 4. Check user state
	updatedState := userStateManager.GetUserState(userID)
	if updatedState.CurrentQuestionID != "question_1" {
		t.Errorf("Expected current question 'question_1', got %s", updatedState.CurrentQuestionID)
	}

	answer, exists := updatedState.GetAnswer("start")
	if !exists || answer != "Continue" {
		t.Errorf("Expected answer 'Continue', got %s (exists: %v)", answer, exists)
	}

	// 5. Test branching - choose "Not good"
	helpQuestion, err := questionManager.GetNextQuestion("question_1", notGoodAnswer)
	if err != nil {
		t.Fatalf("Failed to get help question: %v", err)
	}

	if helpQuestion.ID != "help" {
		t.Errorf("Expected help question, got %s", helpQuestion.ID)
	}

	// 6. Complete flow
	userStateManager.UpdateCurrentQuestion(userID, "help")
	userState.AddAnswer("question_1", notGoodAnswer)

	endQuestion, err := questionManager.GetNextQuestion("help", "Thanks")
	if err != nil {
		t.Fatalf("Failed to get end question: %v", err)
	}

	if endQuestion.ID != "end" {
		t.Errorf("Expected end question, got %s", endQuestion.ID)
	}

	// 7. Check final state
	finalState := userStateManager.GetUserState(userID)
	if len(finalState.Answers) != 2 {
		t.Errorf("Expected 2 answers, got %d", len(finalState.Answers))
	}
}

func TestIntegrationConfigAndQuestions(t *testing.T) {
	// Test configuration loading with valid values
	testConfig := &models.Config{
		TelegramToken:     "test_token",
		GoogleCreds:       "test-creds.json",
		SheetID:           "test_sheet",
		DelayMs:           1000,
		StartQuestionID:   "start",
		QuestionsFilePath: "configs/questions.json",
	}

	// Check configuration validation
	err := testConfig.Validate()
	if err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}

	// Test integration with services
	questions := map[string]models.Question{
		testConfig.StartQuestionID: {
			ID:   testConfig.StartQuestionID,
			Text: "Start question",
		},
	}

	questionManager := services.NewQuestionManager(questions)

	// Check that start question exists
	if !questionManager.QuestionExists(testConfig.StartQuestionID) {
		t.Errorf("Start question %s should exist", testConfig.StartQuestionID)
	}

	startQuestion, err := questionManager.GetQuestion(testConfig.StartQuestionID)
	if err != nil {
		t.Errorf("Failed to get start question: %v", err)
	}

	// Check delay
	delay := startQuestion.GetDelayMs(testConfig.DelayMs)
	if delay != testConfig.DelayMs {
		t.Errorf("Expected delay %d, got %d", testConfig.DelayMs, delay)
	}
}

func TestIntegrationConcurrentUsers(t *testing.T) {
	// Test handling multiple users
	userStateManager := services.NewUserStateManager()

	// Create multiple users
	users := []struct {
		id   int64
		name string
	}{
		{1, "User1"},
		{2, "User2"},
		{3, "User3"},
	}

	// Create states for all users
	for i := range users {
		user := users[i]
		state := userStateManager.GetOrCreateUserState(user.id, user.name)
		state.CurrentQuestionID = "start"
		state.AddAnswer("test", user.name)
	}

	// Check that states don't interfere with each other
	for i := range users {
		user := users[i]
		state := userStateManager.GetUserState(user.id)
		if state == nil {
			t.Errorf("State for user %d should exist", user.id)
			continue
		}

		if state.Name != user.name {
			t.Errorf("Expected name %s for user %d, got %s", user.name, user.id, state.Name)
		}

		answer, exists := state.GetAnswer("test")
		if !exists || answer != user.name {
			t.Errorf("Expected answer %s for user %d, got %s", user.name, user.id, answer)
		}
	}
}
