package sample

type TrxData struct {
	TrxID            string  `db:"trxID"            json:"trxID"`
	UniqueIdentifier string  `db:"uniqueIdentifier" json:"uniqueIdentifier"`
	Type             string  `db:"type"             json:"type"`
	Bank             string  `db:"bank"             json:"bank"`
	TransactionTime  string  `db:"transactionTime"  json:"transactionTime"`
	Date             string  `db:"date"             json:"date"`
	IsSystemTrx      bool    `db:"is_system_trx"    json:"is_system_trx"`
	IsBankTrx        bool    `db:"is_bank_trx"      json:"is_bank_trx"`
	Amount           float64 `db:"amount"           json:"amount"`
}
