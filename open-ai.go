package main

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

/* Function that sends the user message to open ai api */
func (b *BotClient) sendToAI(text string) string {
	resp, err := b.Openai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible.",
				},
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
