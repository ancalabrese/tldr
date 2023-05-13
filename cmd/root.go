package cmd

import (
	"errors"
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
	modeFlagShortName  string = "M"
)

func init() {
	cobra.OnInitialize(initConfigFunc)
}

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd = &cobra.Command{
		Use:   "tldr",
		Short: "Too Long; Didn't read.",
		Long:  "TL;DR - Summarize any long text and ask any questions for more context.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("no valid command")
			}

			f.Llm = openai.NewClient(apiToken)
			m, err := cmdutil.WhichMode(mode)
			if err != nil {
				return err
			}
			f.ConversationMode = m
			return nil
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
