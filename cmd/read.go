package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/ancalabrese/tldr/pkg/kb"
	"github.com/spf13/cobra"
)

func NewReadCmd(f *cmdutil.Factory) *cobra.Command {
	message := ""
	cmd := &cobra.Command{
		Use:   "read [uri]",
		Short: "Read content from source",
		Long:  "Read from an external source to set the knowledge base for this conversation",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				uri *url.URL
				err error
			)
			if len(args) == 0 && message == "" {
				return errors.New("no valid arguments")
			}

			if message != "" {
				uri, err = cmdutil.CreateTmpUri(message, f)
			} else {
				uri, err = url.Parse(args[0])
				if !strings.Contains(uri.Scheme, "file") && !strings.Contains(uri.Scheme, "http") {
					return fmt.Errorf("URI not supported: %s", uri.String())
				}
			}

			if err != nil {
				return err
			}

			f.Kb = kb.New(uri)
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "[Lorem ipsum dolor sit amet] Set the message")

	return cmd
}
