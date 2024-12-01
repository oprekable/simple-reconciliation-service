package process

import (
	"context"
	"fmt"
	"os"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/service"
	"strconv"
	"strings"

	"github.com/spf13/afero"

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
	summary, err := h.svc.SvcProcess.GenerateReconciliation(context.Background(), afero.NewOsFs())
	if err != nil {
		return err
	}

	args := [][]string{
		{
			fmt.Sprintf("-%s --%s", root.FlagFromDateShort, root.FlagFromDate),
			root.FlagFromDateValue,
		},
		{
			fmt.Sprintf("-%s --%s", root.FlagToDateShort, root.FlagToDate),
			root.FlagToDateValue,
		},
		{
			fmt.Sprintf("-%s --%s", root.FlagSystemTRXPathShort, root.FlagSystemTRXPath),
			root.FlagSystemTRXPathValue,
		},
		{
			fmt.Sprintf("-%s --%s", root.FlagBankTRXPathShort, root.FlagBankTRXPath),
			root.FlagBankTRXPathValue,
		},
		{
			fmt.Sprintf("-%s --%s", root.FlagReportTRXPathShort, root.FlagReportTRXPath),
			root.FlagReportTRXPathValue,
		},
		{
			fmt.Sprintf("-%s --%s", root.FlagListBankShort, root.FlagListBank),
			strings.Join(root.FlagListBankValue, ","),
		},
	}

	fmt.Println("")
	tableArgs := tablewriter.NewWriter(os.Stdout)
	tableArgs.SetHeader([]string{"Config", "Value"})
	tableArgs.SetBorder(false)
	tableArgs.SetAlignment(tablewriter.ALIGN_LEFT)
	tableArgs.AppendBulk(args)
	tableArgs.Render()
	fmt.Println("")

	dataDesc := [][]string{
		{"Total number of transactions processed", strconv.FormatInt(summary.TotalProcessedSystemTrx, 10)},
		{"Total number of matched transactions", strconv.FormatInt(summary.TotalMatchedSystemTrx, 10)},
		{"Total number of not matched transactions", strconv.FormatInt(summary.TotalNotMatchedSystemTrx, 10)},
		{"Sum amount all transactions", fmt.Sprintf("%f", summary.SumAmountProcessedSystemTrx)},
		{"Sum amount matched transactions", fmt.Sprintf("%f", summary.SumAmountMatchedSystemTrx)},
		{"Total discrepancies", fmt.Sprintf("%f", summary.SumAmountDiscrepanciesSystemTrx)},
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
