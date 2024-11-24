package sample

import (
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/err"
	"simple-reconciliation-service/internal/inject"
	"simple-reconciliation-service/internal/pkg/utils/atexit"
	"simple-reconciliation-service/variable"
	"time"

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
	app.GetComponents().Config.Reconciliation.SystemTRXPath = root.FlagSystemTRXPathValue
	app.GetComponents().Config.Reconciliation.BankTRXPath = root.FlagBankTRXPathValue
	app.GetComponents().Config.Reconciliation.ReportTRXPath = root.FlagReportTRXPathValue
	app.GetComponents().Config.Reconciliation.ListBank = root.FlagListBankValue

	toDate, er := time.Parse("2006-01-02", root.FlagToDateValue)
	if er != nil {
		return er
	}

	app.GetComponents().Config.Reconciliation.ToDate = toDate

	fromDate, er := time.Parse("2006-01-02", root.FlagFromDateValue)
	if er != nil {
		return er
	}

	app.GetComponents().Config.Reconciliation.FromDate = fromDate
	app.GetComponents().Config.Reconciliation.TotalData = root.FlagTotalDataSampleToGenerateValue
	app.GetComponents().Config.Reconciliation.PercentageMatch = root.FlagPercentageMatchSampleToGenerateValue

	atexit.Add(cleanup)
	app.Start()

	return nil
}
