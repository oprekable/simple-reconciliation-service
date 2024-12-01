package process

import (
	"context"

	"github.com/spf13/afero"
)

type Service interface {
	GenerateReconciliation(ctx context.Context, fs afero.Fs) (returnSummary ReconciliationSummary, err error)
}
