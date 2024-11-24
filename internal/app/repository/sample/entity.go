package sample

type TrxData struct {
	TrxID            string `db:"trxID"            json:"trxID"`
	UniqueIdentifier string `db:"uniqueIdentifier" json:"uniqueIdentifier"`
	Type             string `db:"type"             json:"type"`
	Bank             string `db:"bank"             json:"bank"`
	TransactionTime  string `db:"transactionTime"  json:"transactionTime"`
	Date             string `db:"date"             json:"date"`
	IsSystemTrx      int    `db:"is_system_trx"    json:"is_system_trx"`
	Amount           int    `db:"amount"           json:"amount"`
}

type SystemTrxData struct {
	TrxID           string `db:"trxID"           json:"trxID"`
	Type            string `db:"type"            json:"type"`
	TransactionTime string `db:"transactionTime" json:"transactionTime"`
	Amount          int    `db:"amount"          json:"amount"`
}

type BankTrxData struct {
	UniqueIdentifier string `db:"uniqueIdentifier" json:"uniqueIdentifier"`
	Date             string `db:"date"             json:"date"`
	Amount           int    `db:"amount"           json:"amount"`
}
