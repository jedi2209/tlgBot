// Package main implements the telegram bot entry point.
package main

import (
	"log"

	"tlgbot/internal/bot"
	"tlgbot/internal/config"
	"tlgbot/internal/handlers"
	"tlgbot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation error: %v", err)
	}

	// Create bot API
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized as %s", botAPI.Self.UserName)

	// Load questions
	questionsMap, err := config.LoadQuestions("configs/questions.json")
	if err != nil {
		log.Fatalf("Failed to load questions: %v", err)
	}

	// Create services
	userStateManager := services.NewUserStateManager()
	questionManager := services.NewQuestionManager(questionsMap)

	// Create bot
	telegramBot := bot.NewTelegramBot(botAPI, cfg, userStateManager, questionManager)

	// Create handler
	handler := handlers.NewTelegramHandler(telegramBot, cfg, userStateManager, questionManager)

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
