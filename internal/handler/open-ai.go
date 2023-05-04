package handler

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

/* Function that sends the user message to open ai api */
func (b *BotClient) sendToAI(text string, messages []openai.ChatCompletionMessage) string {

	resp, err := b.Openai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
	}

	return resp.Choices[0].Message.Content
}
