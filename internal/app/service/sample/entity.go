package sample

type Summary struct {
	TotalBankTrx   map[string]int64
	FileBankTrx    map[string]string
	FileSystemTrx  string
	TotalSystemTrx int64
}

type SystemTrxData struct {
	TrxID           string
	Type            string
	TransactionTime string
	Amount          float64
}

type DefaultBankTrxData struct {
	UniqueIdentifier string
	Date             string
	Amount           float64
}

type BCABankTrxData struct {
	BCAUniqueIdentifier string
	BCADate             string
	BCAAmount           float64
}

type BNIBankTrxData struct {
	BNIUniqueIdentifier string
	BNIDate             string
	BNIAmount           float64
}
