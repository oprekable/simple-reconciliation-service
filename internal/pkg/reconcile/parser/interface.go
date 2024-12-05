package parser

import "context"

type ReconcileBankData interface {
	GetBank() string
	GetParser() BankParser
	ToBankTrxData(ctx context.Context, filePath string) (returnData []*BankTrxData, err error)
	ToSql(ctx context.Context, filePath string, sqlPattern string) (returnData string, err error)
}

type ReconcileSystemData interface {
	GetParser() SystemParser
	ToSystemTrxData(ctx context.Context, filePath string) (returnData []*SystemTrxData, err error)
	ToSql(ctx context.Context, filePath string, sqlPattern string) (returnData string, err error)
}
