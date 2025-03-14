package process

import (
	"context"
	"fmt"
	"os"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/handler/hcli/_helper"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/service"
	"simple-reconciliation-service/internal/pkg/utils/memstats"

	"github.com/dustin/go-humanize"

	"github.com/olekukonko/tablewriter"
)

const name = "process"

type Handler struct {
	comp *component.Components
	svc  *service.Services
	repo *repository.Repositories
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Exec() error {
	if h.comp == nil || h.svc == nil || h.repo == nil {
		return nil
	}
	bar := _helper.InitProgressBar()
	formatText := "-%s --%s"
	args := _helper.InitCommonArgs(
		[][]string{
			{
				fmt.Sprintf(formatText, root.FlagReportTRXPathShort, root.FlagReportTRXPath),
				root.FlagReportTRXPathValue,
			},
		},
	)

	tableArgs := tablewriter.NewWriter(os.Stdout)
	tableArgs.SetHeader([]string{"Config", "Value"})
	tableArgs.SetBorder(false)
	tableArgs.SetAlignment(tablewriter.ALIGN_LEFT)
	tableArgs.AppendBulk(args)
	tableArgs.Render()
	fmt.Println("")

	summary, err := h.svc.SvcProcess.GenerateReconciliation(context.Background(), h.comp.Fs.LocalStorageFs, bar)
	if err != nil {
		return err
	}

	numberIntegerFormat := "#.###,"
	numberFloatFormat := "#.###,##"

	dataDesc := [][]string{
		{"Total number of transactions processed", humanize.FormatInteger(numberIntegerFormat, int(summary.TotalProcessedSystemTrx))},
		{"Total number of matched transactions", humanize.FormatInteger(numberIntegerFormat, int(summary.TotalMatchedSystemTrx))},
		{"Total number of not matched transactions", humanize.FormatInteger(numberIntegerFormat, int(summary.TotalNotMatchedSystemTrx))},
		{"Sum amount all transactions", humanize.FormatFloat(numberFloatFormat, summary.SumAmountProcessedSystemTrx)},
		{"Sum amount matched transactions", humanize.FormatFloat(numberFloatFormat, summary.SumAmountMatchedSystemTrx)},
		{"Total discrepancies", humanize.FormatFloat(numberFloatFormat, summary.SumAmountDiscrepanciesSystemTrx)},
	}

	fmt.Println("")
	tableDesc := tablewriter.NewWriter(os.Stdout)
	tableDesc.SetHeader([]string{"Description", "Value"})
	tableDesc.SetBorder(false)
	tableDesc.SetAlignment(tablewriter.ALIGN_LEFT)
	tableDesc.SetAutoWrapText(false)
	tableDesc.AppendBulk(dataDesc)
	tableDesc.Render()
	fmt.Println("")

	dataFilePath := [][]string{
		{"Matched system transaction data", summary.FileMatchedSystemTrx},
	}

	if summary.FileMissingSystemTrx != "" {
		dataFilePath = append(
			dataFilePath,
			[]string{"Missing system transaction data", summary.FileMissingSystemTrx},
		)
	}

	for bank, value := range summary.FileMissingBankTrx {
		dataFilePath = append(
			dataFilePath,
			[]string{
				fmt.Sprintf("Missing bank statement data - %s", bank),
				value,
			},
		)
	}

	fmt.Println("")
	tableFilePath := tablewriter.NewWriter(os.Stdout)
	tableFilePath.SetHeader([]string{"Description", "File Path"})
	tableFilePath.SetBorder(false)
	tableFilePath.SetAlignment(tablewriter.ALIGN_LEFT)
	tableFilePath.SetAutoWrapText(false)
	tableFilePath.AppendBulk(dataFilePath)
	tableFilePath.Render()
	fmt.Println("")

	bar.Describe("[cyan]Done")
	memstats.PrintMemoryStats()

	return nil
}

func (h *Handler) Name() string {
	return name
}

func (h *Handler) SetComponents(c *component.Components) {
	h.comp = c
}
func (h *Handler) SetServices(s *service.Services) {
	h.svc = s
}
func (h *Handler) SetRepositories(r *repository.Repositories) {
	h.repo = r
}
