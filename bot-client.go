package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
)

func runBot(botClient *BotClient, openAiKey string, tBotToken string) {

	var err error

	botClient.Openai = openai.NewClient(openAiKey)

	botClient.Telegram, err = tgbotapi.NewBotAPI(tBotToken)
	if err != nil {
		log.Panic(err)
	}

	botClient.Telegram.Debug = true

	log.Printf("Authorized on account %s", botClient.Telegram.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botClient.Telegram.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			responseText := botClient.sendToAI(update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
			msg.ReplyToMessageID = update.Message.MessageID

			botClient.Telegram.Send(msg)
		}
	}
}
