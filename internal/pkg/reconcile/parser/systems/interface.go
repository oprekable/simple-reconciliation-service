package systems

import "context"

//go:generate mockery --name "ReconcileSystemData" --output "./_mock" --outpkg "_mock"
type ReconcileSystemData interface {
	GetParser() SystemParserType
	ToSystemTrxData(ctx context.Context, filePath string) (returnData []*SystemTrxData, err error)
	ToSql(ctx context.Context, filePath string, sqlPattern string) (returnData string, err error)
}

//go:generate mockery --name "SystemTrxDataInterface" --output "./_mock" --outpkg "_mock"
type SystemTrxDataInterface interface {
	GetTrxID() string
	GetTransactionTime() string
	GetAmount() float64
	GetType() TrxType
	ToSystemTrxData() *SystemTrxData
}
