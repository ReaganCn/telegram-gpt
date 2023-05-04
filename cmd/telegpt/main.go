package telegpt

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	handler "github.com/reagancn/telegram-gpt/internal/handler"
	db "github.com/reagancn/telegram-gpt/pkg/database"
)

func Main() {

	// Load .env file
	err := godotenv.Load()

	// Get environment variables
	openAiKey := os.Getenv("OPEN_AI_KEY")
	tBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	mongoURI := os.Getenv("MONGO_URI")

	// Check if .env file is loaded
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB and get client
	mongoDBClient, ctx := db.ConnectMongoDB(mongoURI)

	// Disconnect from MongoDB when main() ends
	defer mongoDBClient.Disconnect(ctx)

	// Initialize bot client
	botClient := &handler.BotClient{}

	// Run bot
	handler.RunBot(botClient, mongoDBClient, ctx, openAiKey, tBotToken)

}
