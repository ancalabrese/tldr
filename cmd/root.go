package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiToken string
	cmd      = &cobra.Command{
		Use:   "tldr",
		Short: "Too Long; Didn't read.",
		Long:  "TL;DR - Summarize any long text and ask any questions for more context.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(args)
			return nil
		},
	}
)

const (
	flagName      string = "token"
	flagShortName string = "t"
)

func init() {
	cobra.OnInitialize(initConfigFunc)

	cmd.PersistentFlags().StringVarP(&apiToken, flagName, flagShortName, "", "<API_TOKEN> Set the OpenAI API token")
}

func Execute() error {
	return cmd.Execute()
}

func initConfigFunc() {
	if present := cmd.PersistentFlags().Changed(flagName); !present {
		key := os.Getenv("OPENAI_KEY")
		_ = cmd.PersistentFlags().Set(flagName, key)
	}
}
