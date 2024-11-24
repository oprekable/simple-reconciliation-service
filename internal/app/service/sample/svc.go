package sample

import (
	"context"
	"fmt"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/repository/sample"

	"github.com/aaronjan/hunch"
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

func (s *Svc) GenerateReport(ctx context.Context) (returnSummary Summary, err error) {
	var trxData []sample.TrxData
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, s.repo.RepoSample.Pre(
				c,
				s.comp.Config.Reconciliation.ListBank,
				s.comp.Config.Reconciliation.FromDate,
				s.comp.Config.Reconciliation.ToDate,
				s.comp.Config.Reconciliation.TotalData,
				s.comp.Config.Reconciliation.PercentageMatch,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			return s.repo.RepoSample.GetTrx(
				c,
			)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			trxData = i.([]sample.TrxData)
			return nil, s.repo.RepoSample.Post(
				c,
			)
		},
	)

	fmt.Printf("Generated %v data\n", len(trxData))

	return
}
