package cmd

import (
	"simple-reconciliation-service/cmd/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     version.Usage,
	Aliases: version.Aliases,
	Short:   version.Short,
	Long:    version.Long,
	RunE:    version.Runner,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
