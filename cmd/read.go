package cmd

import (
	"errors"
	"net/url"

	"github.com/ancalabrese/tldr/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewReadCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read [uri]",
		Short: "Read content from source",
		Long:  "Read from an external source to set the knowledge base for this conversation",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("uri not found")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := url.Parse(args[0])
			if err != nil {
				return err
			}
			//TODO: open URI and set KB path in config then return
			return nil
		},
	}

	return cmd
}
