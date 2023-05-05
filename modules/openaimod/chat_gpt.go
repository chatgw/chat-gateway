package openaimod

import (
	"bytes"
	"context"
	"strings"

	"github.com/hanyuancheung/gpt-go"
)

type ChatGpt struct {
	client gpt.Client
}

func NewChatGpt(client gpt.Client) *ChatGpt {
	return &ChatGpt{client}
}

// GetResponse get response from gpt3
func (c *ChatGpt) GetResponse(ctx context.Context, question string) (string, error) {
	uestionParam := validateQuestion(question)
	if len(uestionParam) == 0 {
		return "", nil
	}

	buf := bytes.NewBuffer(nil)
	err := c.client.CompletionStreamWithEngine(ctx, &gpt.CompletionRequest{
		Model: gpt.TextDavinci003Engine,
		Prompt: []string{
			uestionParam,
		},
		MaxTokens:   3000,
		Temperature: 0,
	}, func(resp *gpt.CompletionResponse) {
		buf.WriteString(resp.Choices[0].Text)
	})
	if err != nil {
		return "", err
	}
	return strings.Trim(buf.String(), "\n "), nil
}

// NullWriter is a writer on which all Write calls succeed
type NullWriter int

// Write implements io.Writer
func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func validateQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "loop", "break", "continue", "cls", "exit", "block"}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}
