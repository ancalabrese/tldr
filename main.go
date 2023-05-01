package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ancalabrese/tldr/cmd"
	"github.com/sashabaranov/go-openai"
)

var (
	inputTokens string
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal("Err: ", err)
	}

	key := os.Getenv("OPENAI_KEY")
	ctx := context.Background()
	inputChan := make(chan string, 1)
	outputChan := make(chan openai.ChatCompletionResponse, 1)
	client := openai.NewClient(key)
	messages := make([]openai.ChatCompletionMessage, 0)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "system",
		Content: "Create a summary for the provided text. Then answer any user questions about it.",
	},
		openai.ChatCompletionMessage{
			Role:    "user",
			Content: inputTokens,
		})
	go sendNewChatMessage(ctx, client, messages, outputChan)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-inputChan:
				{
					messages := append(messages, openai.ChatCompletionMessage{
						Role:    "user",
						Content: msg,
					})
					go sendNewChatMessage(ctx, client, messages, outputChan)
				}
			case resp := <-outputChan:
				{
					go printResponse(resp)
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    "assistant",
						Content: resp.Choices[0].Message.Content,
					})
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

func sendNewChatMessage(ctx context.Context, llm *openai.Client, chatMessages []openai.ChatCompletionMessage, responseChan chan (openai.ChatCompletionResponse)) {
	req := openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo0301,
		Temperature:      0,
		MaxTokens:        1000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  1,
		Messages:         chatMessages,
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
