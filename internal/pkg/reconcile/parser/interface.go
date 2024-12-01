package parser

import "context"

type ReconcileBankData interface {
	GetBank() string
	GetParser() BankParser
	ToBankTrxData(ctx context.Context, isHaveHeader bool) (returnData []*BankTrxData, err error)
	ToSql(ctx context.Context, isHaveHeader bool, sqlPattern string) (returnData string, err error)
}

type ReconcileSystemData interface {
	GetParser() SystemParser
	ToSystemTrxData(ctx context.Context, isHaveHeader bool) (returnData []*SystemTrxData, err error)
	ToSql(ctx context.Context, isHaveHeader bool, sqlPattern string) (returnData string, err error)
}
