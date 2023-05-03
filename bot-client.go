package main

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func runBot(botClient *BotClient, mongoDBClient *mongo.Client, ctx context.Context, openAiKey string, tBotToken string) {

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

			// Convert chatId to string
			chatId := fmt.Sprintf("%d", update.Message.Chat.ID)

			// Get the collection
			col := GetCollection(mongoDBClient, "telegpt", chatId)

			// Check if the collection is empty
			count, err := col.CountDocuments(ctx, bson.D{})

			if err != nil {
				log.Fatal(err)
			}

			// If the collection is empty, insert the first message
			if count == 0 {
				_, err = col.InsertOne(ctx, BotMessage{
					// Get current time in milliseconds
					TimeStamp: time.Now().UnixNano() / int64(time.Millisecond),
					Role:      "system",
					Content:   "You are a helpful assistant.",
				})
			}

			// Insert the message into the database
			_, err = col.InsertOne(ctx, BotMessage{
				// Get current time in milliseconds
				TimeStamp: time.Now().UnixNano() / int64(time.Millisecond),
				Role:      "user",
				Content:   update.Message.Text,
			})

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Get the last 5 messages from the database
			var messagesContext []openai.ChatCompletionMessage

			opts := options.Find().SetSort(bson.D{{"timestamp", 1}})

			cur, err := col.Find(ctx, bson.D{}, opts)

			if err != nil {
				log.Fatal(err)
			}

			defer cur.Close(ctx)

			for cur.Next(ctx) {
				var result BotMessage
				err := cur.Decode(&result)
				if err != nil {
					log.Fatal(err)
				}

				var chatCompletionMessage openai.ChatCompletionMessage

				chatCompletionMessage.Role = result.Role
				chatCompletionMessage.Content = result.Content

				messagesContext = append(messagesContext, chatCompletionMessage)
			}
			if err := cur.Err(); err != nil {
				log.Fatal(err)
			}

			responseText := botClient.sendToAI(update.Message.Text, messagesContext)

			// Insert the bot message into the database
			_, err = col.InsertOne(ctx, BotMessage{
				// Get current time in milliseconds
				TimeStamp: time.Now().UnixNano() / int64(time.Millisecond),
				Role:      "assistant",
				Content:   responseText,
			})

			if err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
			msg.ReplyToMessageID = update.Message.MessageID

			botClient.Telegram.Send(msg)
		}
	}
}

type BotMessage struct {
	TimeStamp int64  `json:"timestamp"`
	Role      string `json:"role"`
	Content   string `json:"content"`
}
