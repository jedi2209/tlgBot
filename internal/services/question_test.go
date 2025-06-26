package services

import (
	"testing"

	"tlgbot/internal/models"
)

func TestQuestionManagerGetQuestion(t *testing.T) {
	// Create test data
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
			Text: "How are you today?",
			Options: []models.Option{
				{Text: "Good", NextID: "end"},
				{Text: "Bad", NextID: "end"},
			},
		},
	}

	manager := NewQuestionManager(questions)

	// Test getting existing question
	question, err := manager.GetQuestion("start")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if question == nil {
		t.Fatal("Expected question, got nil")
	}
	if question.ID != "start" {
		t.Errorf("Expected question ID 'start', got %s", question.ID)
	}
	if question.Text != "Welcome! What's your name?" {
		t.Errorf("Expected specific text, got %s", question.Text)
	}

	// Test getting non-existent question
	question, err = manager.GetQuestion("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent question")
	}
	if question != nil {
		t.Error("Expected nil question for nonexistent ID")
	}
}

func TestQuestionManagerGetAllQuestions(t *testing.T) {
	questions := map[string]models.Question{
		"start": {ID: "start", Text: "Start question"},
		"end":   {ID: "end", Text: "End question"},
	}

	manager := NewQuestionManager(questions)
	allQuestions := manager.GetAllQuestions()

	if len(allQuestions) != 2 {
		t.Errorf("Expected 2 questions, got %d", len(allQuestions))
	}

	if _, exists := allQuestions["start"]; !exists {
		t.Error("Expected 'start' question to exist")
	}
	if _, exists := allQuestions["end"]; !exists {
		t.Error("Expected 'end' question to exist")
	}
}

func TestQuestionManagerQuestionExists(t *testing.T) {
	questions := map[string]models.Question{
		"existing": {ID: "existing", Text: "Test question"},
	}

	manager := NewQuestionManager(questions)

	// Test existing question
	if !manager.QuestionExists("existing") {
		t.Error("Expected question to exist")
	}

	// Test non-existent question
	if manager.QuestionExists("nonexistent") {
		t.Error("Expected question not to exist")
	}
}

func TestQuestionManagerGetNextQuestion(t *testing.T) {
	questions := map[string]models.Question{
		"start": {
			ID:   "start",
			Text: "Choose option:",
			Options: []models.Option{
				{Text: "Option A", NextID: "question_a"},
				{Text: "Option B", NextID: "question_b"},
			},
		},
		"question_a": {ID: "question_a", Text: "You chose A"},
		"question_b": {ID: "question_b", Text: "You chose B"},
	}

	manager := NewQuestionManager(questions)

	// Test getting next question by correct option
	nextQuestion, err := manager.GetNextQuestion("start", "Option A")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if nextQuestion == nil {
		t.Fatal("Expected next question, got nil")
	}
	if nextQuestion.ID != "question_a" {
		t.Errorf("Expected question ID 'question_a', got %s", nextQuestion.ID)
	}

	// Test with non-existent option
	_, err = manager.GetNextQuestion("start", "Nonexistent Option")
	if err == nil {
		t.Error("Expected error for nonexistent option")
	}

	// Test with non-existent question
	_, err = manager.GetNextQuestion("nonexistent", "Option A")
	if err == nil {
		t.Error("Expected error for nonexistent question")
	}
}

func TestQuestionManagerEmptyQuestions(t *testing.T) {
	questions := make(map[string]models.Question)
	manager := NewQuestionManager(questions)

	// Test with empty questions set
	question, err := manager.GetQuestion("any")
	if err == nil {
		t.Error("Expected error for empty questions map")
	}
	if question != nil {
		t.Error("Expected nil question for empty questions map")
	}

	allQuestions := manager.GetAllQuestions()
	if len(allQuestions) != 0 {
		t.Errorf("Expected 0 questions, got %d", len(allQuestions))
	}
}
