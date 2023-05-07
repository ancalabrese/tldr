package cmdutil

import (
	"github.com/sashabaranov/go-openai"
)

type Factory struct {
	Llm              *openai.Client
	ConversationMode Mode

	ExecutableName string
}
