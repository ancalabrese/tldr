package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command = &cobra.Command{
		Use:   "tldr [command] [flags]",
		Short: "Too Long; Didn't read.",
		Long:  "TL;DR - Summarize any long text and ask any questions for more context.",
	}

	cmdPath *cobra.Command = &cobra.Command{
		Use:   "read <file path>",
		Short: "Read from a file",
		Long:  "Read a file and set its content as the knowledge base for your queries",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}

	token string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(cmdPath)
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "<API_TOKEN> Set the OpenAI API token")
}
