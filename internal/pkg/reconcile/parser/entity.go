package parser

type BankTrxData struct {
	UniqueIdentifier string
	Date             string
	Type             TrxType
	Bank             string
	FilePath         string
	Amount           float64
}

type SystemTrxData struct {
	TrxID           string
	TransactionTime string
	Type            TrxType
	FilePath        string
	Amount          float64
}

type TrxType string

const (
	DEBIT  TrxType = "DEBIT"
	CREDIT TrxType = "CREDIT"
)

type BankParser string

const (
	DefaultBankParser BankParser = "DEFAULT"
	BCABankParser     BankParser = "BCA"
	BNIBankParser     BankParser = "BNI"
)

type SystemParser string

const (
	DefaultSystemParser SystemParser = "DEFAULT"
)
