package parser

type BankTrxData struct {
	UniqueIdentifier string
	Date             string
	Type             TrxType
	Amount           float64
}

type TrxType string

const (
	DEBIT  TrxType = "DEBIT"
	CREDIT TrxType = "CREDIT"
)

type Parser string

const (
	DEFAULT Parser = "DEFAULT"
	BCA     Parser = "BCA"
	BNI     Parser = "BNI"
)
