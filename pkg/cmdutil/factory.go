package cmdutil

import (
	"github.com/ancalabrese/tldr/pkg/kb"
	"github.com/sashabaranov/go-openai"
)

type Factory struct {
	Llm              *openai.Client
	Kb               *kb.Kb
	ConversationMode Mode
	ExecutableName   string
}
