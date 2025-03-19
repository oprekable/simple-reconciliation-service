package helper

import (
	"github.com/spf13/cobra"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/clogger"
	"simple-reconciliation-service/internal/app/component/csqlite"
	"simple-reconciliation-service/internal/app/err"
	"simple-reconciliation-service/internal/inject"
	"simple-reconciliation-service/internal/pkg/utils/atexit"
	"simple-reconciliation-service/variable"
	"strconv"
	"time"
)

func RunnerSubCommand(cmd *cobra.Command, _ []string, dBPath csqlite.DBPath) (er error) {
	defer func() {
		atexit.AtExit()
	}()

	app, cleanup, er := inject.WireApp(
		cmd.Context(),
		root.EmbedFS,
		cconfig.AppName(variable.AppName),
		cconfig.TimeZone(root.FlagTZValue),
		err.RegisteredErrorType,
		clogger.IsShowLog(root.FlagIsVerboseValue),
		dBPath,
	)

	atexit.Add(cleanup)

	if er != nil {
		return er
	}

	app.GetComponents().Config.Data.Reconciliation.Action = cmd.Use
	app.GetComponents().Config.Data.Reconciliation.SystemTRXPath = root.FlagSystemTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.BankTRXPath = root.FlagBankTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.ReportTRXPath = root.FlagReportTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.ListBank = root.FlagListBankValue
	app.GetComponents().Config.Data.Reconciliation.IsDeleteCurrentSampleDirectory = root.FlagIsDeleteCurrentSampleDirectoryValue
	app.GetComponents().Config.Data.IsShowLog = strconv.FormatBool(root.FlagIsVerboseValue)
	app.GetComponents().Config.Data.IsDebug = root.FlagIsDebugValue

	toDate, er := time.Parse(root.DateFormatString, root.FlagToDateValue)
	if er != nil {
		return er
	}

	app.GetComponents().Config.Data.Reconciliation.ToDate = toDate

	fromDate, er := time.Parse(root.DateFormatString, root.FlagFromDateValue)
	if er != nil {
		return er
	}

	app.GetComponents().Config.Data.Reconciliation.FromDate = fromDate
	app.GetComponents().Config.Data.Reconciliation.TotalData = root.FlagTotalDataSampleToGenerateValue
	app.GetComponents().Config.Data.Reconciliation.PercentageMatch = root.FlagPercentageMatchSampleToGenerateValue
	app.Start()
	return nil
}
