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
)

type FormatterFunc func(string) string
type MaxTokenConverter func(string) int

func formatInput(suffix string) FormatterFunc {
	return func(text string) string {
		return suffix + text
	}
}

func setMaxToken(inputTokens *string) MaxTokenConverter {
	wordCount := len(strings.Split(*inputTokens, " "))

	return func(l string) int {
		switch l {
		case "m":
			return int(0.5 * float32(wordCount))
		case "l":
			return int(0.7 * float32(wordCount))
		case "s":
		default:
			return int(0.3 * float32(wordCount))
		}
		return 0
	}
}

func main() {

	inputTokens = flag.String("m", "", textHelpMessage)
	maxLen = flag.String("l", "s", lenHelpMessage)
	flag.Parse()

	if *inputTokens == "" {
		log.Fatal("Error missing option:  " + textHelpMessage)
	}

	if !strings.ContainsAny(*maxLen, "sml") {
		log.Println("Error wrong format: ", lenHelpMessage, "  defaulting to 's'")
		*maxLen = "s"
	}

	key := os.Getenv("OPENAI_KEY")
	ctx := context.Background()
	client := openai.NewClient(key)
	formatFunc := formatInput("Summarize text:\n")
	setMaxTokenFunc := setMaxToken(inputTokens)

	req := openai.CompletionRequest{
		Model:            openai.GPT3TextCurie001,
		Temperature:      0.7,
		MaxTokens:        setMaxTokenFunc(*maxLen),
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  1,
		Prompt:           formatFunc(*inputTokens),
	}

	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Text)
}
