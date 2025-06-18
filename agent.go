package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, error)
}

func NewAgent(client *anthropic.Client, getUserMessage func() (string, error)) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
	}
}

func (a *Agent) run(ctx context.Context) error {
	conversation := []anthropic.MessageParam{}

	fmt.Println("Chat with claude (Press Ctrl-C to quit)")

	for {
		fmt.Print("> ")

		userInput, err := a.getUserMessage()

		if err != nil {
			return err
		}

		// todo grab previous context
		userMessage := anthropic.NewUserMessage(
			anthropic.NewTextBlock(userInput),
		)

		conversation = append(conversation, userMessage)

		response, err := a.inference(ctx, conversation)

		if err != nil {
			return err
		}

		// todo stream output.
		fmt.Println("┌" + strings.Repeat("─", 60) + "┐")
		fmt.Printf("%+v\n", response.Content[0].Text)
		fmt.Println("└" + strings.Repeat("─", 60) + "┘")

		conversation = append(conversation, response.ToParam())

		fmt.Println()
	}
}

func (a *Agent) inference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	response, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		Messages:  conversation,
		MaxTokens: 1024,
	})

	return response, err
}
