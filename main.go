package main

import (
	"log"
	"os"

	"github.com/korjavin/dstgbot/api"
	"github.com/korjavin/dstgbot/cache"
	"github.com/korjavin/dstgbot/telegram"
)

func main() {
	log.Println("Starting the bot...")

	// Initialize cache
	msgCache := cache.NewMessageCache(300)
	log.Println("Cache initialized.")

	// Get environment variables
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	deepseekKey := os.Getenv("DEEPSEEK_APIKEY")
	groupID := os.Getenv("TG_GROUP_ID")

	// Initialize Telegram bot
	bot, err := telegram.NewBot(telegramToken, groupID, msgCache)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	log.Println("Telegram bot initialized.")

	// Initialize DeepSeek API client
	deepseekClient := api.NewClient(deepseekKey)
	log.Println("DeepSeek API client initialized.")

	// Start bot
	if err := bot.Start(deepseekClient); err != nil {
		log.Fatalf("Bot failed: %v", err)
	}
}
