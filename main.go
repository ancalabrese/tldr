package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ancalabrese/tldr/cmd"
	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/ancalabrese/tldr/pkg/conversation"
	"github.com/ancalabrese/tldr/pkg/factory"
	"github.com/sashabaranov/go-openai"
)

func main() {
	factory := factory.Defaults()
	err := cmd.NewRootCmd(factory).Execute()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	inputChan := make(chan string, 1)
	outputChan := make(chan openai.ChatCompletionResponse, 1)
	convo := conversation.New(cmdutil.Tldr)

	go sendNewChatMessage(ctx, factory.Llm, convo, outputChan)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-inputChan:
				{
					convo.AddRequest(msg)
					go sendNewChatMessage(ctx, factory.Llm, convo, outputChan)

				}
			case resp := <-outputChan:
				{
					go printResponse(resp)
					convo.AddResponse(resp.Choices[0].Message.Content)

				}
			}
		}
	}()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan() // use `for scanner.Scan()` to keep reading
		input := scanner.Text()
		inputChan <- input
	}
}

func sendNewChatMessage(ctx context.Context, llm *openai.Client, conversation *conversation.Conversation,
	responseChan chan (openai.ChatCompletionResponse)) {
	req := openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo0301,
		Temperature:      0,
		MaxTokens:        1000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  1,
		Messages:         conversation.History,
	}

	resp, err := llm.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return
	}
	responseChan <- resp
}

func printResponse(resp openai.ChatCompletionResponse) {
	fmt.Println(resp.Choices[0].Message.Content)
}
