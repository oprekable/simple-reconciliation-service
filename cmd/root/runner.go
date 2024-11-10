package root

import "github.com/spf13/cobra"

func Runner(cmd *cobra.Command, _ []string) (er error) {
	return cmd.Help()
}
