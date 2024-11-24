package sample

import "context"

type Service interface {
	GenerateReport(ctx context.Context) (returnSummary Summary, err error)
}
