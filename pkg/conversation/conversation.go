package conversation

import (
	"fmt"
	"strings"

	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/ancalabrese/tldr/pkg/kb"
	"github.com/sashabaranov/go-openai"
)

type Convo struct {
	History []openai.ChatCompletionMessage
	Kb      *kb.Kb
}

func New(mode cmdutil.Mode, kb *kb.Kb) (*Convo, error) {
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
	kbContent, err := kb.Parse()
	if err != nil {
		return nil, fmt.Errorf("can't initiate new conversation: %w", err)
	}

	history = append(history, openai.ChatCompletionMessage{
		Role:    "user",
		Content: strings.Join(kbContent, "\n"),
	})

	return &Convo{
		History: history,
		Kb:      kb,
	}, nil
}

func (c *Convo) AddResponse(resp string) {
	c.History = append(c.History, openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: resp,
	})
}

func (c *Convo) AddRequest(req string) {
	c.History = append(c.History, openai.ChatCompletionMessage{
		Role:    "user",
		Content: req,
	})
}
