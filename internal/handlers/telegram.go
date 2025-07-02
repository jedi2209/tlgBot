// Package handlers provides Telegram event handlers for processing user interactions.
package handlers

import (
	"log"

	"tlgbot/internal/bot"
	"tlgbot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" //nolint:depguard // Required for Telegram bot functionality
)

// TelegramHandler handles Telegram events
type TelegramHandler struct {
	bot              models.BotService
	config           *models.Config
	userStateManager models.UserStateService
	questionManager  models.QuestionService
}

// NewTelegramHandler creates a new Telegram handler
func NewTelegramHandler(
	telegramBot models.BotService,
	config *models.Config,
	userStateManager models.UserStateService,
	questionManager models.QuestionService,
) *TelegramHandler {
	return &TelegramHandler{
		bot:              telegramBot,
		config:           config,
		userStateManager: userStateManager,
		questionManager:  questionManager,
	}
}

// HandleMessage handles incoming messages
func (h *TelegramHandler) HandleMessage(message *tgbotapi.Message) {
	userID := message.From.ID
	userName := bot.GetTelegramName(message.From)

	// Get or create user state
	userState := h.userStateManager.GetOrCreateUserState(userID, userName)

	if message.IsCommand() {
		h.handleCommand(message, userState)
	} else {
		h.handleTextInput(message, userState)
	}
}

// HandleCallbackQuery handles callback queries
func (h *TelegramHandler) HandleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	data := callback.Data

	// Acknowledge callback - это КРИТИЧЕСКИ ВАЖНО для Telegram!
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := h.bot.GetAPI().Request(callbackConfig); err != nil {
		log.Printf("Error acknowledging callback: %v", err)
	}

	userState := h.userStateManager.GetUserState(userID)
	if userState == nil {
		log.Printf("User state not found for user %d", userID)
		return
	}

	// Process option selection
	if err := h.bot.ProcessOptionAnswer(userID, data); err != nil {
		log.Printf("Error processing option answer: %v", err)
	}
}

// handleCommand handles commands
func (h *TelegramHandler) handleCommand(message *tgbotapi.Message, userState *models.UserState) {
	switch message.Command() {
	case "start":
		h.startConversation(message.From.ID, userState)
	default:
		log.Printf("Unknown command: %s", message.Command())
	}
}

// handleTextInput handles text input
func (h *TelegramHandler) handleTextInput(message *tgbotapi.Message, userState *models.UserState) {
	if userState.CurrentQuestionID == "" {
		return
	}

	currentQuestion, err := h.questionManager.GetQuestion(userState.CurrentQuestionID)
	if err != nil {
		log.Printf("Failed to get current question: %v", err)
		return
	}

	// Check if text input is expected
	if currentQuestion.InputType != "" {
		// Save answer and move to next question
		if err := h.bot.ProcessAnswer(message.From.ID, message.Text); err != nil {
			log.Printf("Failed to process answer: %v", err)
			return
		}

		// Move to next question if there are options
		if len(currentQuestion.Options) > 0 {
			nextQuestionID := currentQuestion.Options[0].NextID
			if err := h.moveToNextQuestion(message.From.ID, nextQuestionID); err != nil {
				log.Printf("Failed to move to next question: %v", err)
			}
		}
	}
}

// startConversation starts conversation with user
func (h *TelegramHandler) startConversation(userID int64, userState *models.UserState) {
	startQuestion, err := h.questionManager.GetQuestion(h.config.StartQuestionID)
	if err != nil {
		log.Printf("Failed to get start question: %v", err)
		return
	}

	userState.CurrentQuestionID = h.config.StartQuestionID
	h.userStateManager.SetUserState(userID, userState)

	if err := h.bot.ProcessQuestion(userID, startQuestion); err != nil {
		log.Printf("Failed to process start question: %v", err)
		return
	}

	// Handle automatic transition
	if err := h.bot.HandleAutoAdvance(userID, startQuestion); err != nil {
		log.Printf("Failed to handle auto advance: %v", err)
	}
}

// moveToNextQuestion moves to next question
func (h *TelegramHandler) moveToNextQuestion(userID int64, nextQuestionID string) error {
	nextQuestion, err := h.questionManager.GetQuestion(nextQuestionID)
	if err != nil {
		return err
	}

	h.userStateManager.UpdateCurrentQuestion(userID, nextQuestionID)

	if err := h.bot.ProcessQuestion(userID, nextQuestion); err != nil {
		return err
	}

	return h.bot.HandleAutoAdvance(userID, nextQuestion)
}
