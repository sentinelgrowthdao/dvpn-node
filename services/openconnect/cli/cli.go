package cli

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "openconnect",
		Aliases: nil,
		Short:   "OpenConnect sub-commands",
	}

	cmd.AddCommand(
		configCmd(),
	)

	return cmd
}
