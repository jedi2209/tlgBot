// Package bot provides Telegram bot functionality for handling messages and user interactions.
package bot

import (
	"fmt"
	"strings"
	"time"

	"tlgbot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" //nolint:depguard // Required for Telegram bot functionality
)

// TelegramBot represents a Telegram bot
type TelegramBot struct {
	api              *tgbotapi.BotAPI
	config           *models.Config
	userStateManager models.UserStateService
	questionManager  models.QuestionService
}

// NewTelegramBot creates a new bot instance
func NewTelegramBot(
	api *tgbotapi.BotAPI,
	config *models.Config,
	userStateManager models.UserStateService,
	questionManager models.QuestionService,
) *TelegramBot {
	return &TelegramBot{
		api:              api,
		config:           config,
		userStateManager: userStateManager,
		questionManager:  questionManager,
	}
}

// SendImages sends images to user
func (bot *TelegramBot) SendImages(userID int64, images []string, delay int) error {
	if len(images) == 0 {
		return nil
	}

	if len(images) > 1 {
		return bot.sendMultipleImages(userID, images, delay)
	}

	return bot.sendSingleImage(userID, images[0], delay)
}

// sendMultipleImages sends multiple images as media group
func (bot *TelegramBot) sendMultipleImages(userID int64, images []string, delay int) error {
	mediaPhotos := bot.buildMediaGroup(images)
	if len(mediaPhotos) == 0 {
		return nil
	}

	// Convert to []interface{} for compatibility with NewMediaGroup
	media := make([]interface{}, len(mediaPhotos))
	for i, photo := range mediaPhotos {
		media[i] = photo
	}

	mediaGroup := tgbotapi.NewMediaGroup(userID, media)
	// Send media group returns []tgbotapi.Message, not tgbotapi.Message
	_, err := bot.api.SendMediaGroup(mediaGroup)
	if err != nil {
		return fmt.Errorf("failed to send media group: %w", err)
	}

	time.Sleep(time.Duration(delay) * time.Millisecond)
	return nil
}

// sendSingleImage sends a single image
func (bot *TelegramBot) sendSingleImage(userID int64, imagePath string, delay int) error {
	if imagePath == "" {
		return nil
	}

	photoMsg := tgbotapi.NewPhoto(userID, tgbotapi.FilePath(imagePath))
	_, err := bot.api.Send(photoMsg)
	if err != nil {
		return fmt.Errorf("failed to send photo: %w", err)
	}

	time.Sleep(time.Duration(delay) * time.Millisecond)
	return nil
}

// buildMediaGroup creates media group from image paths
func (bot *TelegramBot) buildMediaGroup(images []string) []tgbotapi.InputMediaPhoto {
	media := make([]tgbotapi.InputMediaPhoto, 0, len(images))
	for _, imgPath := range images {
		if imgPath != "" {
			photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(imgPath))
			media = append(media, photo)
		}
	}
	return media
}

// SendMessage sends a single message with keyboard
func (bot *TelegramBot) SendMessage(userID int64, text string, keyboard interface{}) error {
	msg := tgbotapi.NewMessage(userID, text)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	_, err := bot.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// SendMessages sends multiple messages with keyboard on the last one
func (bot *TelegramBot) SendMessages(userID int64, messages []string, userName string, keyboard interface{}) error {
	for i, msgTmpl := range messages {
		msgText := bot.replaceNamePlaceholder(msgTmpl, userName)
		msg := tgbotapi.NewMessage(userID, msgText)

		// Add keyboard to the last message
		if i == len(messages)-1 && keyboard != nil {
			msg.ReplyMarkup = keyboard
		}

		_, err := bot.api.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message %d: %w", i, err)
		}

		// Add delay between messages, except for the last one
		if i != len(messages)-1 {
			time.Sleep(time.Duration(bot.config.DelayMs) * time.Millisecond)
		}
	}
	return nil
}

// BuildKeyboard builds inline keyboard for a question
func (bot *TelegramBot) BuildKeyboard(q *models.Question) *tgbotapi.InlineKeyboardMarkup {
	if !q.HasKeyboard() {
		return nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	// Add option buttons
	for _, opt := range q.Options {
		btn := tgbotapi.NewInlineKeyboardButtonData(opt.Text, opt.Text)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	// Add external link if present
	if q.ExternalLink != "" && q.ExternalText != "" {
		btn := tgbotapi.NewInlineKeyboardButtonURL(q.ExternalText, q.ExternalLink)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard
}

// ProcessQuestion processes sending a question to user
func (bot *TelegramBot) ProcessQuestion(userID int64, question *models.Question) error {
	userState := bot.userStateManager.GetUserState(userID)
	if userState == nil {
		return fmt.Errorf("user state not found for user %d", userID)
	}

	// Send images
	if len(question.Images) > 0 {
		delay := question.GetDelayMs(bot.config.DelayMs)
		if err := bot.SendImages(userID, question.Images, delay); err != nil {
			return fmt.Errorf("failed to send images: %w", err)
		}
	}

	// Build keyboard
	keyboard := bot.BuildKeyboard(question)

	// Send messages
	if len(question.Messages) > 0 {
		if err := bot.SendMessages(userID, question.Messages, userState.Name, keyboard); err != nil {
			return fmt.Errorf("failed to send messages: %w", err)
		}
	} else if question.Text != "" {
		text := bot.replaceNamePlaceholder(question.Text, userState.Name)

		// For final question add answers summary
		if question.ID == "end" {
			text += bot.generateAnswersSummary(userState.Answers)
		}

		if err := bot.SendMessage(userID, text, keyboard); err != nil {
			return fmt.Errorf("failed to send text message: %w", err)
		}
	}

	return nil
}

// HandleAutoAdvance handles automatic transition to next question
func (bot *TelegramBot) HandleAutoAdvance(userID int64, question *models.Question) error {
	if !question.AutoAdvance {
		return nil
	}

	delayMs := question.AutoAdvanceDelayMs
	if delayMs == 0 {
		delayMs = bot.config.DelayMs
	}

	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	if len(question.Options) > 0 {
		nextQuestionID := question.Options[0].NextID
		nextQuestion, err := bot.questionManager.GetQuestion(nextQuestionID)
		if err != nil {
			return fmt.Errorf("failed to get next question: %w", err)
		}

		bot.userStateManager.UpdateCurrentQuestion(userID, nextQuestionID)

		if err := bot.ProcessQuestion(userID, nextQuestion); err != nil {
			return fmt.Errorf("failed to process next question: %w", err)
		}

		return bot.HandleAutoAdvance(userID, nextQuestion)
	}

	return nil
}

// ProcessAnswer processes user's answer
func (bot *TelegramBot) ProcessAnswer(userID int64, answer string) error {
	userState := bot.userStateManager.GetUserState(userID)
	if userState == nil {
		return fmt.Errorf("user state not found")
	}

	currentQuestion, err := bot.questionManager.GetQuestion(userState.CurrentQuestionID)
	if err != nil {
		return fmt.Errorf("failed to get current question: %w", err)
	}

	// Save user's answer
	questionText := currentQuestion.GetDisplayText()
	userState.AddAnswer(questionText, answer)

	return nil
}

// ProcessOptionAnswer handles user option selection
func (bot *TelegramBot) ProcessOptionAnswer(userID int64, optionText string) error {
	userState := bot.userStateManager.GetUserState(userID)
	if userState == nil {
		return fmt.Errorf("user state not found")
	}

	currentQuestion, err := bot.questionManager.GetQuestion(userState.CurrentQuestionID)
	if err != nil {
		return fmt.Errorf("failed to get current question: %w", err)
	}

	// Find selected option
	var selectedOption *models.Option
	for _, option := range currentQuestion.Options {
		if option.Text == optionText {
			selectedOption = &option
			break
		}
	}

	if selectedOption == nil {
		return fmt.Errorf("option not found: %s", optionText)
	}

	// Save answer
	if err := bot.ProcessAnswer(userID, optionText); err != nil {
		return fmt.Errorf("failed to process answer: %w", err)
	}

	// Handle special actions
	if selectedOption.Action == "get_location" {
		return bot.requestLocation(userID)
	}

	// Move to next question
	return bot.moveToNextQuestion(userID, selectedOption.NextID)
}

// moveToNextQuestion moves to next question
func (bot *TelegramBot) moveToNextQuestion(userID int64, nextQuestionID string) error {
	nextQuestion, err := bot.questionManager.GetQuestion(nextQuestionID)
	if err != nil {
		return fmt.Errorf("failed to get next question: %w", err)
	}

	bot.userStateManager.UpdateCurrentQuestion(userID, nextQuestionID)

	if err := bot.ProcessQuestion(userID, nextQuestion); err != nil {
		return fmt.Errorf("failed to process next question: %w", err)
	}

	return bot.HandleAutoAdvance(userID, nextQuestion)
}

// requestLocation requests user's location
func (bot *TelegramBot) requestLocation(userID int64) error {
	msg := tgbotapi.NewMessage(userID, "Please share your location by clicking the button below:")
	locationBtn := tgbotapi.NewKeyboardButtonLocation("📍 Share Location")
	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{locationBtn})
	keyboard.OneTimeKeyboard = true
	msg.ReplyMarkup = keyboard

	_, err := bot.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send location request: %w", err)
	}
	return nil
}

// replaceNamePlaceholder replaces {name} placeholder with user's name
func (bot *TelegramBot) replaceNamePlaceholder(text, name string) string {
	return strings.ReplaceAll(text, "{name}", name)
}

// generateAnswersSummary generates user's answers summary
func (bot *TelegramBot) generateAnswersSummary(answers map[string]string) string {
	if len(answers) == 0 {
		return ""
	}

	summary := "\n\n📋 Your answers:\n"
	for question, answer := range answers {
		summary += fmt.Sprintf("• %s: %s\n", question, answer)
	}
	return summary
}

// GetAPI returns the Telegram bot API instance
func (bot *TelegramBot) GetAPI() *tgbotapi.BotAPI {
	return bot.api
}

// GetTelegramName gets user name from User object
func GetTelegramName(u *tgbotapi.User) string {
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.UserName != "" {
		return u.UserName
	}
	return "User"
}
