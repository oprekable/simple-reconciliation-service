package parser

import "context"

type ReconcileData interface {
	GetBank() string
	GetParser() Parser
	ToBankTrxData(ctx context.Context, isHaveHeader bool) (returnData []*BankTrxData, err error)
	ToSql(ctx context.Context, isHaveHeader bool) (returnData string, err error)
}
