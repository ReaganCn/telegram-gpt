<p align="center"><img src="etc/assets/gopher.svg" width="250"></p>
<p align="center">
  <a href="https://goreportcard.com/report/go.mongodb.org/mongo-driver"><img src="https://goreportcard.com/badge/go.mongodb.org/mongo-driver"></a>
  <!-- <a href="https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo"><img src="etc/assets/godev-mongo-blue.svg" alt="docs"></a>
  <a href="https://pkg.go.dev/go.mongodb.org/mongo-driver/bson"><img src="etc/assets/godev-bson-blue.svg" alt="docs"></a>
  <a href="https://www.mongodb.com/docs/drivers/go/current/"><img src="etc/assets/docs-mongodb-green.svg"></a> -->
</p>

# The Telegram GPT Integration

An integration of openai's chatgpt api on telegram bot. Such that anyone can access chatgpt's powerful capabilities at any point and easily.

-------------------------
## Requirements

- Go 1.19 or higher. We aim to support the latest versions of Go.

- MongoDB 3.6 and higher.

-------------------------
## Installation

Make sure you have docker, and docker compose installed.

Start the mongo db server with `docker run -d -p 27017:27017 --name test-mongo mongo:latest`

Clone the project and add the below environments to `.env` in the project root folder:
```
OPEN_AI_KEY="your-open-ai-key"
TELEGRAM_BOT_TOKEN="your-telegram-bot-token"
MONGO_URI="your-mongo-db-URI" or simply "mongodb://127.0.0.1:27017/"
```

Install packages with `go mod download`

Execute `go run .`

## Contribution

For help with the bot, for now, you are free to create an issue. Contribution details are still work in progress.

-------------------------

## License

The Telegram GPT Integration is licensed under the [GNU](LICENSE).