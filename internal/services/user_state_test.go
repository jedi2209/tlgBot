package services

import (
	"testing"

	"tlgbot/internal/models"
)

const (
	testUserName   = "TestUser"
	testQuestionID = "start"
)

func TestUserStateManagerGetUserState(t *testing.T) {
	manager := NewUserStateManager()
	userID := int64(123)

	// Test getting non-existent state
	state := manager.GetUserState(userID)
	if state != nil {
		t.Errorf("Expected nil for non-existent user state, got %v", state)
	}
}

func TestUserStateManagerSetAndGetUserState(t *testing.T) {
	manager := NewUserStateManager()
	userID := int64(123)
	userName := testUserName

	// Create user state
	originalState := models.NewUserState(userName)
	originalState.CurrentQuestionID = testQuestionID
	originalState.AddAnswer("test_question", "test_answer")

	// Set and get state
	manager.SetUserState(userID, originalState)
	retrievedState := manager.GetUserState(userID)

	// Check result
	if retrievedState == nil {
		t.Fatal("Expected user state, got nil")
	}

	if retrievedState.Name != userName {
		t.Errorf("Expected name %s, got %s", userName, retrievedState.Name)
	}

	if retrievedState.CurrentQuestionID != testQuestionID {
		t.Errorf("Expected question ID '%s', got %s", testQuestionID, retrievedState.CurrentQuestionID)
	}

	answer, exists := retrievedState.GetAnswer("test_question")
	if !exists {
		t.Error("Expected answer to exist")
	}
	if answer != "test_answer" {
		t.Errorf("Expected answer 'test_answer', got %s", answer)
	}
}

func TestUserStateManagerUpdateCurrentQuestion(t *testing.T) {
	manager := NewUserStateManager()
	userID := int64(123)
	userName := testUserName

	// Create and set state
	state := models.NewUserState(userName)
	state.CurrentQuestionID = testQuestionID
	manager.SetUserState(userID, state)

	// Update current question
	newQuestionID := "question_2"
	manager.UpdateCurrentQuestion(userID, newQuestionID)

	// Check update
	updatedState := manager.GetUserState(userID)
	if updatedState.CurrentQuestionID != newQuestionID {
		t.Errorf("Expected question ID %s, got %s", newQuestionID, updatedState.CurrentQuestionID)
	}
}

func TestUserStateManagerGetOrCreateUserState(t *testing.T) {
	manager := NewUserStateManager()
	userID := int64(123)
	userName := testUserName

	// Get state for non-existent user
	state := manager.GetOrCreateUserState(userID, userName)
	if state == nil {
		t.Fatal("Expected new user state, got nil")
	}

	if state.Name != userName {
		t.Errorf("Expected name %s, got %s", userName, state.Name)
	}

	// Get state for existing user
	state.CurrentQuestionID = "test_question"
	existingState := manager.GetOrCreateUserState(userID, "DifferentName")

	if existingState.Name != userName {
		t.Errorf("Expected original name %s, got %s", userName, existingState.Name)
	}

	if existingState.CurrentQuestionID != "test_question" {
		t.Errorf("Expected question ID 'test_question', got %s", existingState.CurrentQuestionID)
	}
}

func TestUserStateManagerConcurrentAccess(t *testing.T) {
	manager := NewUserStateManager()
	userID := int64(123)
	userName := testUserName

	// Test thread safety with concurrent access
	done := make(chan bool, 2)

	// Goroutine 1: sets state
	go func() {
		defer func() { done <- true }()
		for i := 0; i < 100; i++ {
			state := models.NewUserState(userName)
			manager.SetUserState(userID, state)
		}
	}()

	// Goroutine 2: reads state
	go func() {
		defer func() { done <- true }()
		for i := 0; i < 100; i++ {
			manager.GetUserState(userID)
		}
	}()

	// Wait for goroutines to complete
	<-done
	<-done

	// Test passed if there was no race condition
}
