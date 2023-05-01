package factory

import (
	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/sashabaranov/go-openai"
)

func Defaults(apiToken string) *cmdutil.Factory {
	return &cmdutil.Factory{
		ExecutableName: "tldr",
		Llm:            openai.NewClient(apiToken),
	}

}
