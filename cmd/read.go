package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/ancalabrese/tldr/pkg/conversation"
	"github.com/ancalabrese/tldr/pkg/kb"
	"github.com/spf13/cobra"
)

func NewReadCmd(c *conversation.Convo) *cobra.Command {
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
			uri, err := url.Parse(args[0])
			if err != nil {
				return err
			}
			if !strings.Contains(uri.Scheme, "file") || !strings.Contains(uri.Scheme, "http") {
				return fmt.Errorf("URI not supported: %s", uri.String())
			}

			c.Kb = kb.New(*uri)
			return nil
		},
	}

	return cmd
}
