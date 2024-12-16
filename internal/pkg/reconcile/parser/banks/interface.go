package banks

import "context"

type ReconcileBankData interface {
	GetBank() string
	GetParser() BankParserType
	ToBankTrxData(ctx context.Context, filePath string) (returnData []*BankTrxData, err error)
	ToSql(ctx context.Context, filePath string, sqlPattern string) (returnData string, err error)
}

type BankTrxDataInterface interface {
	GetUniqueIdentifier() string
	GetDate() string
	GetAmount() float64
	GetAbsAmount() float64
	GetType() TrxType
	GetBank() string
	ToBankTrxData() *BankTrxData
}
