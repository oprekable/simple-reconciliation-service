package process

import (
	"context"
	"simple-reconciliation-service/internal/pkg/reconcile/parser"
	"time"
)

type Repository interface {
	Pre(
		ctx context.Context,
		listBank []string,
		startDate time.Time,
		toDate time.Time,
	) (err error)

	ImportSystemTrx(ctx context.Context, data []*parser.SystemTrxData) (err error)
	ImportBankTrx(ctx context.Context, data []*parser.BankTrxData) (err error)
	GenerateReconciliationMap(ctx context.Context) (err error)
	Post(ctx context.Context) (err error)
	Close() (err error)
}
