package parser

type BankTrxData struct {
	UniqueIdentifier string
	Date             string
	Type             TrxType
	Bank             string
	Amount           float64
}

type SystemTrxData struct {
	TrxID           string
	TransactionTime string
	Type            TrxType
	Amount          float64
}

type TrxType string

const (
	DEBIT  TrxType = "DEBIT"
	CREDIT TrxType = "CREDIT"
)

type BankParser string

const (
	DEFAULT_BANK BankParser = "DEFAULT"
	BCA          BankParser = "BCA"
	BNI          BankParser = "BNI"
)

type SystemParser string

const (
	DEFAULT_SYSTEM SystemParser = "DEFAULT"
)
