package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// Load .env file
	err := godotenv.Load()

	// Get environment variables
	openAiKey := os.Getenv("OPEN_AI_KEY")
	tBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	// Check if .env file is loaded
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize bot client
	botClient := &BotClient{}

	// Run bot
	runBot(botClient, openAiKey, tBotToken)

}
