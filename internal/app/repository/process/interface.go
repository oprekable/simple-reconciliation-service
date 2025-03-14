package process

import (
	"context"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"time"
)

//go:generate mockery --name "Repository" --output "./_mock" --outpkg "_mock"
type Repository interface {
	Pre(
		ctx context.Context,
		listBank []string,
		startDate time.Time,
		toDate time.Time,
	) (err error)

	ImportSystemTrx(ctx context.Context, data []*systems.SystemTrxData) (err error)
	ImportBankTrx(ctx context.Context, data []*banks.BankTrxData) (err error)
	GenerateReconciliationMap(ctx context.Context, minAmount float64, maxAmount float64) (err error)
	GetReconciliationSummary(ctx context.Context) (returnData ReconciliationSummary, err error)
	Post(ctx context.Context) (err error)
	Close() (err error)
	GetMatchedTrx(ctx context.Context) (returnData []MatchedTrx, err error)
	GetNotMatchedSystemTrx(ctx context.Context) (returnData []NotMatchedSystemTrx, err error)
	GetNotMatchedBankTrx(ctx context.Context) (returnData []NotMatchedBankTrx, err error)
}
