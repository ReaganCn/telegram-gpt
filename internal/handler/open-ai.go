package handler

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

/* Function that sends the user message to open ai api */
func (b *BotClient) sendToAI(text string, messages []openai.ChatCompletionMessage) string {

	messagesToSend := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant. Be as accurate as possible.",
		},
	}

	messagesToSend = append(messagesToSend, messages...)

	resp, err := b.Openai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messagesToSend,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
	}

	return resp.Choices[0].Message.Content
}
