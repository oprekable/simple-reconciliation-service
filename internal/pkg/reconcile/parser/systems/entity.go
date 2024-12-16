package systems

type SystemParserType string

const (
	DefaultSystemParser SystemParserType = "DEFAULT"
)

type TrxType string

type SystemTrxData struct {
	TrxID           string
	TransactionTime string
	Type            TrxType
	FilePath        string
	Amount          float64
}
