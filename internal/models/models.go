// Package models provides data structures and interfaces for the Telegram bot application.
package models

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Config structure for storing settings
type Config struct {
	TelegramToken     string `json:"telegram_token"`
	GoogleCreds       string `json:"google_creds"`
	SheetID           string `json:"sheet_id"`
	DelayMs           int    `json:"delay_ms"`
	StartQuestionID   string `json:"start_question_id"`
	QuestionsFilePath string `json:"questions_file_path"`
}

// Validate checks configuration correctness
func (c *Config) Validate() error {
	if c.TelegramToken == "" {
		return errors.New("telegram token is required")
	}
	if c.StartQuestionID == "" {
		return errors.New("start question ID is required")
	}
	if c.DelayMs < 0 {
		return errors.New("delay must be non-negative")
	}
	return nil
}

// Option represents an answer option for a question
type Option struct {
	Text   string `json:"text"`
	NextID string `json:"next_id"`
	Action string `json:"action"`
}

// Question represents a survey question
type Question struct {
	ID                 string   `json:"id"`
	Text               string   `json:"text"`
	Messages           []string `json:"messages"`
	Images             []string `json:"images"`
	Options            []Option `json:"options"`
	ExternalLink       string   `json:"external_link"`
	ExternalText       string   `json:"external_text"`
	Input              string   `json:"input"`
	InputType          string   `json:"input_type"`
	InputPlaceholder   string   `json:"input_placeholder"`
	DelayMs            *int     `json:"delay_ms"`
	AutoAdvance        bool     `json:"auto_advance"`
	AutoAdvanceDelayMs int      `json:"auto_advance_delay_ms"`
}

// GetDelayMs returns delay for question or default value
func (q *Question) GetDelayMs(defaultDelay int) int {
	if q.DelayMs != nil {
		return *q.DelayMs
	}
	return defaultDelay
}

// HasKeyboard checks if keyboard is needed for this question
func (q *Question) HasKeyboard() bool {
	return !q.AutoAdvance && (len(q.Options) > 0 || (q.ExternalLink != "" && q.ExternalText != ""))
}

// GetDisplayText returns text for display
func (q *Question) GetDisplayText() string {
	if q.Text != "" {
		return q.Text
	}
	if len(q.Messages) > 0 {
		return q.Messages[len(q.Messages)-1] // Return last message
	}
	return ""
}

// UserState represents user state
type UserState struct {
	CurrentQuestionID string
	Answers           map[string]string
	Name              string
}

// NewUserState creates new user state
func NewUserState(name string) *UserState {
	return &UserState{
		Name:    name,
		Answers: make(map[string]string),
	}
}

// AddAnswer adds user answer
func (us *UserState) AddAnswer(question, answer string) {
	us.Answers[question] = answer
}

// GetAnswer returns user answer to question
func (us *UserState) GetAnswer(question string) (string, bool) {
	answer, exists := us.Answers[question]
	return answer, exists
}

// BotService interface for working with bot
type BotService interface {
	SendImages(userID int64, images []string, delay int) error
	SendMessage(userID int64, text string, keyboard interface{}) error
	SendMessages(userID int64, messages []string, userName string, keyboard interface{}) error
	ProcessQuestion(userID int64, question *Question) error
	ProcessAnswer(userID int64, answer string) error
	ProcessOptionAnswer(userID int64, optionText string) error
	HandleAutoAdvance(userID int64, question *Question) error
	GetAPI() *tgbotapi.BotAPI // Returns Telegram Bot API instance
}

// QuestionService interface for working with questions
type QuestionService interface {
	GetQuestion(id string) (*Question, error)
	GetAllQuestions() map[string]Question
}

// UserStateService interface for working with user states
type UserStateService interface {
	GetUserState(userID int64) *UserState
	SetUserState(userID int64, state *UserState)
	UpdateCurrentQuestion(userID int64, questionID string)
	GetOrCreateUserState(userID int64, userName string) *UserState
}
