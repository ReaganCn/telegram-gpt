package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type BotClient struct {
	id       string
	openai   *openai.Client
	telegram *tgbotapi.BotAPI
}

func main() {

	err := godotenv.Load()

	openAiKey := os.Getenv("OPEN_AI_KEY")
	tBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botClient := &BotClient{}

	botClient.openai = openai.NewClient(openAiKey)

	botClient.telegram, err = tgbotapi.NewBotAPI(tBotToken)
	if err != nil {
		log.Panic(err)
	}

	botClient.telegram.Debug = true

	log.Printf("Authorized on account %s", botClient.telegram.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botClient.telegram.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			responseText := botClient.sendToAI(update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
			msg.ReplyToMessageID = update.Message.MessageID

			botClient.telegram.Send(msg)
		}
	}

}

func (b *BotClient) sendToAI(text string) string {
	resp, err := b.openai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
	}

	return resp.Choices[0].Message.Content
}
