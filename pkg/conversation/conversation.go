package conversation

import (
	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/sashabaranov/go-openai"
)

type Conversation struct {
	History []openai.ChatCompletionMessage
}

func New(mode cmdutil.Mode) *Conversation {
	command := ""

	switch mode {
	case cmdutil.Interactive:
		command = "Create a summary for the provided text. Then answer any user questions about it."
	case cmdutil.Tldr:
	default:
		command = "Create a summary for the provided text."
	}

	history := make([]openai.ChatCompletionMessage, 1)
	history[0] = openai.ChatCompletionMessage{
		Role:    "system",
		Content: command,
	}

	return &Conversation{
		History: history,
	}
}

func (c *Conversation) AddResponse(resp string) {
	c.History = append(c.History, openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: resp,
	})
}

func (c *Conversation) AddRequest(req string) {
	c.History = append(c.History, openai.ChatCompletionMessage{
		Role:    "user",
		Content: req,
	})
}
