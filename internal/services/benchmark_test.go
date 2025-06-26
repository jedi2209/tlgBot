package services

import (
	"fmt"
	"testing"

	"tlgbot/internal/models"
)

func BenchmarkUserStateManagerGetOrCreateUserState(b *testing.B) {
	manager := NewUserStateManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := int64(i % 1000) // Cyclically use 1000 users
		userName := fmt.Sprintf("User%d", userID)
		manager.GetOrCreateUserState(userID, userName)
	}
}

func BenchmarkUserStateManagerSetAndGetUserState(b *testing.B) {
	manager := NewUserStateManager()
	userID := int64(123)
	userName := "BenchmarkUser"
	state := models.NewUserState(userName)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.SetUserState(userID, state)
		manager.GetUserState(userID)
	}
}

func BenchmarkUserStateManagerConcurrentAccess(b *testing.B) {
	manager := NewUserStateManager()
	userID := int64(123)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Use different userIDs to avoid race condition
			currentUserID := userID + int64(i%100)
			userName := fmt.Sprintf("ConcurrentUser%d", currentUserID)
			manager.GetOrCreateUserState(currentUserID, userName)
			manager.UpdateCurrentQuestion(currentUserID, "question_1")
			i++
		}
	})
}

func BenchmarkQuestionManagerGetQuestion(b *testing.B) {
	// Create large set of questions
	questions := make(map[string]models.Question, 1000)
	for i := 0; i < 1000; i++ {
		id := fmt.Sprintf("question_%d", i)
		questions[id] = models.Question{
			ID:   id,
			Text: fmt.Sprintf("Question %d", i),
		}
	}

	manager := NewQuestionManager(questions)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		questionID := fmt.Sprintf("question_%d", i%1000)
		_, _ = manager.GetQuestion(questionID)
	}
}

func BenchmarkQuestionManagerGetNextQuestion(b *testing.B) {
	questions := map[string]models.Question{
		"start": {
			ID:   "start",
			Text: "Start question",
			Options: []models.Option{
				{Text: "Option 1", NextID: "question_1"},
				{Text: "Option 2", NextID: "question_2"},
				{Text: "Option 3", NextID: "question_3"},
			},
		},
		"question_1": {ID: "question_1", Text: "Question 1"},
		"question_2": {ID: "question_2", Text: "Question 2"},
		"question_3": {ID: "question_3", Text: "Question 3"},
	}

	manager := NewQuestionManager(questions)
	options := []string{"Option 1", "Option 2", "Option 3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		option := options[i%len(options)]
		_, _ = manager.GetNextQuestion("start", option)
	}
}

func BenchmarkUserStateAddAnswer(b *testing.B) {
	state := models.NewUserState("BenchmarkUser")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		questionID := fmt.Sprintf("question_%d", i%100)
		answer := fmt.Sprintf("answer_%d", i)
		state.AddAnswer(questionID, answer)
	}
}

func BenchmarkUserStateGetAnswer(b *testing.B) {
	state := models.NewUserState("BenchmarkUser")

	// Pre-populate answers
	for i := 0; i < 100; i++ {
		questionID := fmt.Sprintf("question_%d", i)
		answer := fmt.Sprintf("answer_%d", i)
		state.AddAnswer(questionID, answer)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		questionID := fmt.Sprintf("question_%d", i%100)
		state.GetAnswer(questionID)
	}
}

func BenchmarkQuestionGetDelayMs(b *testing.B) {
	question := models.Question{
		DelayMs: intPtr(1000),
	}
	defaultDelay := 500

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		question.GetDelayMs(defaultDelay)
	}
}

func BenchmarkQuestionHasKeyboard(b *testing.B) {
	question := models.Question{
		AutoAdvance: false,
		Options: []models.Option{
			{Text: "Option 1"},
			{Text: "Option 2"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		question.HasKeyboard()
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
