// Package services provides business logic services for question management and user state handling.
package services

import (
	"errors"

	"tlgbot/internal/models"
)

// QuestionManager manages questions
type QuestionManager struct {
	questions map[string]models.Question
}

// NewQuestionManager creates a new question manager
func NewQuestionManager(questions map[string]models.Question) *QuestionManager {
	return &QuestionManager{
		questions: questions,
	}
}

// GetQuestion returns a question by ID
func (m *QuestionManager) GetQuestion(id string) (*models.Question, error) {
	question, exists := m.questions[id]
	if !exists {
		return nil, errors.New("question not found")
	}
	return &question, nil
}

// GetAllQuestions returns all questions
func (m *QuestionManager) GetAllQuestions() map[string]models.Question {
	return m.questions
}

// QuestionExists checks if a question exists
func (m *QuestionManager) QuestionExists(id string) bool {
	_, exists := m.questions[id]
	return exists
}

// GetNextQuestion returns the next question based on option
func (m *QuestionManager) GetNextQuestion(currentQuestionID, optionText string) (*models.Question, error) {
	currentQ, err := m.GetQuestion(currentQuestionID)
	if err != nil {
		return nil, err
	}

	for _, option := range currentQ.Options {
		if option.Text == optionText {
			return m.GetQuestion(option.NextID)
		}
	}

	return nil, errors.New("option not found")
}
