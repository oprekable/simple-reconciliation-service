package process

import (
	"fmt"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/service"
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
	fmt.Println(name)
	fmt.Println(h.comp.Config.Reconciliation.BankTRXPath)
	fmt.Println(h.comp.Config.Reconciliation.SystemTRXPath)
	fmt.Println(h.comp.Config.Reconciliation.ArchivePath)
	fmt.Println(h.comp.Config.Reconciliation.ListBank)
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
