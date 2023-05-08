package conversation

import (
	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/sashabaranov/go-openai"
)

type Convo struct {
	History []openai.ChatCompletionMessage
	kb      string
}

func New(mode cmdutil.Mode, kbUri string) *Convo {
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

	return &Convo{
		History: history,
		kb:      kbUri,
	}
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
