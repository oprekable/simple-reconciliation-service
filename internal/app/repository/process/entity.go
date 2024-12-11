package process

type ReconciliationData struct {
	SystemTrxID           *string  `db:"SystemTrxID"`
	SystemTransactionTime *string  `db:"SystemTransactionTime"`
	Type                  *string  `db:"type"`
	BankUniqueIdentifier  *string  `db:"BankUniqueIdentifier"`
	Bank                  *string  `db:"Bank"`
	SystemFilePath        *string  `db:"-"`
	BankFilePath          *string  `db:"-"`
	Amount                *float64 `db:"amount"`
}

type ReconciliationSummary struct {
	TotalSystemTrx      int64   `db:"total_system_trx"`
	TotalMatchedTrx     int64   `db:"total_matched_trx"`
	TotalNotMatchedTrx  int64   `db:"total_not_matched_trx"`
	SumSystemTrx        float64 `db:"sum_system_trx"`
	SumMatchedTrx       float64 `db:"sum_matched_trx"`
	SumDiscrepanciesTrx float64 `db:"sum_discrepancies_trx"`
}
