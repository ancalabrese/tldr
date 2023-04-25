package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var (
	inputTokens     *string
	textHelpMessage string = "-m to set the paragraph to bu summarized"
	maxLen          *string
	lenHelpMessage  string = "-l [s,m,l] to set the summary length"
	kb              *string
	kbHelpMessage   string = "-kb [filepath] to set a file content as knowledge base"
)

func main() {

	inputTokens = flag.String("m", "", textHelpMessage)
	maxLen = flag.String("l", "s", lenHelpMessage)
	kb = flag.String("kb", "", kbHelpMessage)
	flag.Parse()

	if (*inputTokens == "" && *kb == "") || (*inputTokens != "" && *kb != "") {
		log.Fatal("Error: use -h for help")
	}

	if !strings.ContainsAny(*maxLen, "sml") {
		log.Println("Error wrong format: ", lenHelpMessage, "  defaulting to 's'")
		*maxLen = "s"
	}

	key := os.Getenv("OPENAI_KEY")
	ctx := context.Background()
	client := openai.NewClient(key)
	messages := make([]openai.ChatCompletionMessage, 0)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "system",
		Content: "Give me the tl;dr",
	},
		openai.ChatCompletionMessage{
			Role:    "user",
			Content: *inputTokens,
		})

	req := openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo,
		Temperature:      0,
		MaxTokens:        1000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  1,
		Messages:         messages,
		Stream:           false,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
