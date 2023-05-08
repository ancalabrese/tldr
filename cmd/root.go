package cmd

import (
	"os"

	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	apiToken string
	mode     string
	cmd      *cobra.Command
)

const (
	tokenFlagName      string = "token"
	tokenFlagShortName string = "t"
	modeFlagName       string = "mode"
	modeFlagShortName  string = "m"
)

func init() {
	cobra.OnInitialize(initConfigFunc)
}

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd = &cobra.Command{
		Use:   "tldr",
		Short: "Too Long; Didn't read.",
		Long:  "TL;DR - Summarize any long text and ask any questions for more context.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := cmdutil.WhichMode(mode)
			return err
		},
		Run: func(cmd *cobra.Command, args []string) {
			f.Llm = openai.NewClient(apiToken)
			//Ignoring error. Was checked in PersistentPreRun
			m, _ := cmdutil.WhichMode(mode)
			f.ConversationMode = m
		},
	}
	cmd.PersistentFlags().StringVarP(&apiToken, tokenFlagName, tokenFlagShortName, "", "<API_TOKEN> Set the OpenAI API token")
	cmd.PersistentFlags().StringVarP(&mode, modeFlagName, modeFlagShortName, "TLDR", "[TLDR, INTERACTIVE] select the chat mode.")

	cmd.AddCommand(NewReadCmd(f))
	return cmd
}

func initConfigFunc() {
	if present := cmd.PersistentFlags().Changed(tokenFlagName); !present {
		key := os.Getenv("OPENAI_KEY")
		_ = cmd.PersistentFlags().Set(tokenFlagName, key)
	}
}
