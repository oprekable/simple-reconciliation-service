package process

import (
	"context"
	"fmt"
	"os"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/service"
	"simple-reconciliation-service/internal/pkg/utils/memstats"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/afero"
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

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSpinnerType(17),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	formatText := "-%s --%s"
	args := [][]string{
		{
			fmt.Sprintf(formatText, root.FlagFromDateShort, root.FlagFromDate),
			root.FlagFromDateValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagToDateShort, root.FlagToDate),
			root.FlagToDateValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagSystemTRXPathShort, root.FlagSystemTRXPath),
			root.FlagSystemTRXPathValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagBankTRXPathShort, root.FlagBankTRXPath),
			root.FlagBankTRXPathValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagReportTRXPathShort, root.FlagReportTRXPath),
			root.FlagReportTRXPathValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagListBankShort, root.FlagListBank),
			strings.Join(root.FlagListBankValue, ","),
		},
		{
			fmt.Sprintf(formatText, root.FlagIsDeleteCurrentSampleDirectoryShort, root.FlagIsDeleteCurrentSampleDirectory),
			strconv.FormatBool(root.FlagIsDeleteCurrentSampleDirectoryValue),
		},
		{
			fmt.Sprintf(formatText, root.FlagIsVerboseShort, root.FlagIsVerbose),
			strconv.FormatBool(root.FlagIsVerboseValue),
		},
		{
			fmt.Sprintf(formatText, root.FlagIsDebugShort, root.FlagIsDebug),
			strconv.FormatBool(root.FlagIsDebugValue),
		},
	}

	tableArgs := tablewriter.NewWriter(os.Stdout)
	tableArgs.SetHeader([]string{"Config", "Value"})
	tableArgs.SetBorder(false)
	tableArgs.SetAlignment(tablewriter.ALIGN_LEFT)
	tableArgs.AppendBulk(args)
	tableArgs.Render()
	fmt.Println("")

	summary, err := h.svc.SvcProcess.GenerateReconciliation(context.Background(), afero.NewOsFs(), bar)
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
		{"Missing system transaction data", summary.FileMissingSystemTrx},
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
