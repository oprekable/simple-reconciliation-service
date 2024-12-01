package process

import (
	"fmt"
	"os"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/service"
	"strings"

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
