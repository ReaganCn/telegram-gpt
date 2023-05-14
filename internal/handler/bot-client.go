package handler

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

	db "github.com/reagancn/telegram-gpt/pkg/database"
	"github.com/reagancn/telegram-gpt/pkg/utils"
)

func RunBot(botClient *BotClient, mongoDBClient *mongo.Client, ctx context.Context, openAiKey string, tBotToken string) {

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
			col := db.GetCollection(mongoDBClient, "telegpt", chatId)

			var assistantMessage openai.ChatCompletionMessage
			var messagesContext []openai.ChatCompletionMessage

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// messageObject := update.Message
			// messageO, _ := json.Marshal(messageObject)

			if update.Message.ReplyToMessage != nil {

				fmt.Println("\nReplying to message: ", update.Message.ReplyToMessage.Text, "\n")

				// Set the assistant message
				assistantMessage = openai.ChatCompletionMessage{
					Role:    "assistant",
					Content: update.Message.ReplyToMessage.Text,
				}
			} else {
				// Get the last message from the database
				fmt.Println("\nNo reply message. Getting last message from database...\n")
				// Sort by timestamp in descending order
				opts := options.Find().SetSort(bson.D{{"timestamp", -1}})
				opts.SetLimit(1)

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

					assistantMessage.Role = "assistant"
					assistantMessage.Content = result.Content
				}
				if err := cur.Err(); err != nil {
					log.Fatal(err)
				}

				fmt.Println("\nLast message: ", assistantMessage.Content, "\n")
			}

			// Create openai.ChatCompletionMessage for the user message
			userMessage := openai.ChatCompletionMessage{
				Role:    "user",
				Content: update.Message.Text,
			}

			messagesContext = append(messagesContext, assistantMessage, userMessage)

			responseText := botClient.sendToAI(update.Message.Text, messagesContext)

			// Insert the received user message into the database
			_, err = col.InsertOne(ctx, BotMessage{
				// Get current time in milliseconds
				TimeStamp: utils.GetTimeInMilliseconds(),
				Role:      "user",
				Content:   update.Message.Text,
			})

			if err != nil {
				log.Fatal(err)
			}

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
