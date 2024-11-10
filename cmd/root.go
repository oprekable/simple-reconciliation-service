package cmd

import (
	"embed"
	"fmt"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/pkg/utils/filepathhelper"
	"simple-reconciliation-service/variable"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   variable.AppName,
	Short: variable.AppDescShort,
	Long:  variable.AppDescLong,
	Example: fmt.Sprintf(
		"%s\n%s\n",
		fmt.Sprintf("Generate sample \n\t%s sample %s", variable.AppName, root.SampleUsageFlags),
		fmt.Sprintf("Process data \n\t%s process %s", variable.AppName, root.ProcessUsageFlags),
	),
	RunE: root.Runner,
}

func init() {
	defaultTZ := variable.TimeZone
	if defaultTZ == "" {
		defaultTZ = "Asia/Jakarta"
	}

	rootCmd.PersistentFlags().StringVarP(
		&root.FlagTZValue,
		root.FlagTimeZone,
		root.FlagTimeZoneShort,
		defaultTZ,
		root.FlagTimeZoneUsage,
	)

	workDir := filepathhelper.GetWorkDir()
	rootCmd.PersistentFlags().StringVarP(
		&root.SystemTRXPath,
		root.FlagSystemTRXPath,
		root.FlagSystemTRXPathShort,
		fmt.Sprintf("%s/system", workDir),
		root.FlagSystemTRXPathUsage,
	)

	rootCmd.PersistentFlags().StringVarP(
		&root.BankTRXPath,
		root.FlagBankTRXPath,
		root.FlagBankTRXPathShort,
		fmt.Sprintf("%s/bank", workDir),
		root.FlagBankTRXPathUsage,
	)

	rootCmd.PersistentFlags().StringVarP(
		&root.ArchivePath,
		root.FlagArchiveTRXPath,
		root.FlagArchiveTRXPathShort,
		fmt.Sprintf("%s/archive", workDir),
		root.FlagArchiveTRXPathUsage,
	)

	rootCmd.PersistentFlags().StringSliceVarP(
		&root.ListBank,
		root.FlagListBank,
		root.FlagListBankShort,
		root.DefaultListBank,
		root.FlagListBankUsage,
	)
}

func Execute(embedFS *embed.FS) {
	root.RootEmbedFS = embedFS
	if er := rootCmd.Execute(); er != nil {
		fmt.Println(er)
	}
}
