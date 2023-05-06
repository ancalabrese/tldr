package cmd

import (
	"os"

	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	apiToken string
	cmd      *cobra.Command
)

const (
	flagName      string = "token"
	flagShortName string = "t"
)

func init() {
	cobra.OnInitialize(initConfigFunc)
}

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd = &cobra.Command{
		Use:   "tldr",
		Short: "Too Long; Didn't read.",
		Long:  "TL;DR - Summarize any long text and ask any questions for more context.",
		Run: func(cmd *cobra.Command, args []string) {
			f.Llm = openai.NewClient(apiToken)
		},
	}
	cmd.PersistentFlags().StringVarP(&apiToken, flagName, flagShortName, "", "<API_TOKEN> Set the OpenAI API token")

	return cmd
}

func initConfigFunc() {
	if present := cmd.PersistentFlags().Changed(flagName); !present {
		key := os.Getenv("OPENAI_KEY")
		_ = cmd.PersistentFlags().Set(flagName, key)
	}
}
