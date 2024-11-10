package cmd

import (
	"fmt"
	"simple-reconciliation-service/cmd/process"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/variable"

	"github.com/spf13/cobra"
)

var processCmd = &cobra.Command{
	Use:     process.Usage,
	Aliases: process.Aliases,
	Short:   process.Short,
	Long:    process.Long,
	Example: fmt.Sprintf(
		"%s\n",
		fmt.Sprintf("Process data \n\t%s process %s", variable.AppName, root.ProcessUsageFlags),
	),
	RunE: process.Runner,
}

func init() {
	rootCmd.AddCommand(processCmd)
}
