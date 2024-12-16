package systems

import "context"

type ReconcileSystemData interface {
	GetParser() SystemParserType
	ToSystemTrxData(ctx context.Context, filePath string) (returnData []*SystemTrxData, err error)
	ToSql(ctx context.Context, filePath string, sqlPattern string) (returnData string, err error)
}

type SystemTrxDataInterface interface {
	GetTrxID() string
	GetTransactionTime() string
	GetAmount() float64
	GetType() TrxType
	ToSystemTrxData() *SystemTrxData
}
