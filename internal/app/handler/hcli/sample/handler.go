package sample

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

	"github.com/olekukonko/tablewriter"

	"github.com/spf13/afero"
)

const name = "sample"

type Handler struct {
	comp *component.Components
	svc  *service.Services
	repo *repository.Repositories
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Exec() (err error) {
	if h.comp == nil || h.svc == nil || h.repo == nil {
		return nil
	}
	summary, err := h.svc.SvcSample.GenerateSample(context.Background(), afero.NewOsFs())
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
		{
			fmt.Sprintf("-%s --%s", root.FlagTotalDataSampleToGenerateShort, root.FlagTotalDataSampleToGenerate),
			strconv.FormatInt(root.FlagTotalDataSampleToGenerateValue, 10),
		},
		{
			fmt.Sprintf("-%s --%s", root.FlagPercentageMatchSampleToGenerateShort, root.FlagPercentageMatchSampleToGenerate),
			strconv.Itoa(root.FlagPercentageMatchSampleToGenerateValue),
		},
	}

	fmt.Println("")
	tableArgs := tablewriter.NewWriter(os.Stdout)
	tableArgs.SetHeader([]string{"Config", "Value"})
	tableArgs.SetBorder(false)
	tableArgs.SetAlignment(tablewriter.ALIGN_LEFT)
	tableArgs.AppendBulk(args)
	tableArgs.Render()

	data := [][]string{
		{"System Trx", "-", "Total Trx", strconv.FormatInt(summary.TotalSystemTrx, 10)}, //nolint:gofmt
		{"System Trx", "-", "File Path", summary.FileSystemTrx},
	}

	for k, v := range summary.FileBankTrx {
		data = append(
			data,
			[]string{"Bank Trx", k, "Total Trx", strconv.FormatInt(summary.TotalBankTrx[k], 10)},
			[]string{"Bank Trx", k, "File Path", v},
		)
	}

	fmt.Println("")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type Trx", "Bank", "Title", ""})
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
	fmt.Println("")

	return
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
