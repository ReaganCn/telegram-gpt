package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
)

type BotClient struct {
	Id       string
	Openai   *openai.Client
	Telegram *tgbotapi.BotAPI
}
