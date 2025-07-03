// Package main implements the telegram bot entry point.
package main

import (
	"fmt"
	"log"

	"tlgbot/internal/bot"
	"tlgbot/internal/config"
	"tlgbot/internal/handlers"
	"tlgbot/internal/models"
	"tlgbot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// initializeBot initializes all bot components and returns them or an error
func initializeBot() (*tgbotapi.BotAPI, *handlers.TelegramHandler, *models.Config, error) {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, nil, nil, fmt.Errorf("configuration validation error: %w", err)
	}

	// Create bot API
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create bot: %w", err)
	}

	// Load questions
	questionsMap, err := config.LoadQuestions(cfg.QuestionsFilePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load questions: %w", err)
	}

	// Create services
	userStateManager := services.NewUserStateManager()
	questionManager := services.NewQuestionManager(questionsMap)

	// Create bot
	telegramBot := bot.NewTelegramBot(botAPI, cfg, userStateManager, questionManager)

	// Create handler
	handler := handlers.NewTelegramHandler(telegramBot, cfg, userStateManager, questionManager)

	return botAPI, handler, cfg, nil
}

func main() {
	botAPI, handler, _, err := initializeBot()
	if err != nil {
		log.Fatalf("Bot initialization failed: %v", err)
	}

	log.Printf("Authorized as %s", botAPI.Self.UserName)

	// Start processing updates
	startBot(botAPI, handler)
}

// startBot starts processing updates from Telegram
func startBot(botAPI *tgbotapi.BotAPI, handler *handlers.TelegramHandler) {
	// Settings for receiving updates
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := botAPI.GetUpdatesChan(updateConfig)

	log.Println("Bot started. Waiting for messages...")

	for update := range updates {
		// Process messages
		if update.Message != nil {
			go handler.HandleMessage(update.Message)
		}

		// Process callback queries
		if update.CallbackQuery != nil {
			go handler.HandleCallbackQuery(update.CallbackQuery)
		}
	}
}
