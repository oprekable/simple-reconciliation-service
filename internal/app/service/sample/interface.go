package sample

import (
	"context"

	"github.com/spf13/afero"
)

type Service interface {
	GenerateReport(ctx context.Context, fs afero.Fs) (returnSummary Summary, err error)
}
