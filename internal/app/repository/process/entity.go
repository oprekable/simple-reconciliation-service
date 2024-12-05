package process

type ReconciliationData struct {
	SystemTrxID           string  `db:"SystemTrxID"`
	SystemTransactionTime string  `db:"SystemTransactionTime"`
	Type                  string  `db:"type"`
	BankUniqueIdentifier  string  `db:"BankUniqueIdentifier"`
	Bank                  string  `db:"Bank"`
	SystemFilePath        string  `db:"-"`
	BankFilePath          string  `db:"-"`
	Amount                float64 `db:"amount"`
}
