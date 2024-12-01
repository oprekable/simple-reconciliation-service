package process

import (
	"context"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"

	"github.com/spf13/afero"
)

type Svc struct {
	comp *component.Components
	repo *repository.Repositories
}

var _ Service = (*Svc)(nil)

func NewSvc(
	comp *component.Components,
	repo *repository.Repositories,
) *Svc {
	return &Svc{
		comp: comp,
		repo: repo,
	}
}

func (s *Svc) GenerateReconciliation(ctx context.Context, fs afero.Fs) (returnSummary ReconciliationSummary, err error) {
	return
}
