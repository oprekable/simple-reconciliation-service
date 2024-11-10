package sample

import (
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/err"
	"simple-reconciliation-service/internal/inject"
	"simple-reconciliation-service/internal/pkg/utils/atexit"
	"simple-reconciliation-service/variable"

	"github.com/spf13/cobra"
)

func Runner(cmd *cobra.Command, _ []string) (er error) {
	defer func() {
		atexit.AtExit()
	}()

	app, cleanup, er := inject.WireApp(
		cmd.Context(),
		root.EmbedFS,
		cconfig.AppName(variable.AppName),
		cconfig.TimeZone(root.FlagTZValue),
		err.RegisteredErrorType,
	)

	if er != nil {
		return er
	}

	app.GetComponents().Config.Reconciliation.Action = cmd.Use
	app.GetComponents().Config.Reconciliation.SystemTRXPath = root.SystemTRXPath
	app.GetComponents().Config.Reconciliation.BankTRXPath = root.BankTRXPath
	app.GetComponents().Config.Reconciliation.ArchivePath = root.ArchivePath
	app.GetComponents().Config.Reconciliation.ListBank = root.ListBank

	atexit.Add(cleanup)
	app.Start()

	return nil
}
