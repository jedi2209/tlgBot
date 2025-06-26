// Package services provides business logic services for user state management.
package services

import (
	"sync"

	"tlgbot/internal/models"
)

// UserStateManager manages user states
type UserStateManager struct {
	mu     sync.RWMutex
	states map[int64]*models.UserState
}

// NewUserStateManager creates a new user state manager
func NewUserStateManager() *UserStateManager {
	return &UserStateManager{
		states: make(map[int64]*models.UserState),
	}
}

// GetUserState returns user state
func (m *UserStateManager) GetUserState(userID int64) *models.UserState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.states[userID]
	if !exists {
		return nil
	}
	return state
}

// SetUserState sets user state
func (m *UserStateManager) SetUserState(userID int64, state *models.UserState) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.states[userID] = state
}

// UpdateCurrentQuestion updates user's current question
func (m *UserStateManager) UpdateCurrentQuestion(userID int64, questionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if state, exists := m.states[userID]; exists {
		state.CurrentQuestionID = questionID
	}
}

// GetOrCreateUserState returns existing state or creates a new one
func (m *UserStateManager) GetOrCreateUserState(userID int64, userName string) *models.UserState {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.states[userID]
	if !exists {
		state = models.NewUserState(userName)
		m.states[userID] = state
	}
	return state
}
