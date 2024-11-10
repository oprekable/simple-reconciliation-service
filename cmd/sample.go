package cmd

import (
	"fmt"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/cmd/sample"
	"simple-reconciliation-service/variable"

	"github.com/spf13/cobra"
)

var sampleCmd = &cobra.Command{
	Use:     sample.Usage,
	Aliases: sample.Aliases,
	Short:   sample.Short,
	Long:    sample.Long,
	Example: fmt.Sprintf(
		"%s\n",
		fmt.Sprintf("Generate sample \n\t%s sample %s", variable.AppName, root.SampleUsageFlags),
	),
	RunE: sample.Runner,
}

func init() {
	rootCmd.AddCommand(sampleCmd)
}
