package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	openAiKey := os.Getenv("OPEN_AI_KEY")
	tBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botClient := &BotClient{}

	runBot(botClient, openAiKey, tBotToken)

}
