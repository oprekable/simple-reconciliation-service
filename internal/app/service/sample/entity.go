package sample

const (
	DEBIT = "DEBIT"
)

type Summary struct {
	TotalBankTrx   map[string]int64
	FileBankTrx    map[string]string
	FileSystemTrx  string
	TotalSystemTrx int64
}

type SystemTrxDataInterface interface {
	GetTrxID() string
	GetType() string
	GetTransactionTime() string
	GetAmount() float64
}

type SystemTrxData struct {
	TrxID           string
	Type            string
	TransactionTime string
	Amount          float64
}

func NewSystemTrxData(
	TrxID string,
	Type string,
	TransactionTime string,
	Amount float64,
) *SystemTrxData {
	return &SystemTrxData{
		TrxID:           TrxID,
		Type:            Type,
		TransactionTime: TransactionTime,
		Amount:          Amount,
	}
}

func (a *SystemTrxData) GetTrxID() string {
	return a.TrxID
}

func (a *SystemTrxData) GetType() string {
	return a.Type
}

func (a *SystemTrxData) GetTransactionTime() string {
	return a.TransactionTime
}

func (a *SystemTrxData) GetAmount() float64 {
	return a.Amount
}

type BankTrxDataInterface interface {
	GetBank() string
	GetUniqueIdentifier() string
	GetDate() string
	GetAmount() float64
}

type DefaultBankTrxData struct {
	UniqueIdentifier string
	Date             string
	bank             string
	Amount           float64
}

func NewDefaultBankTrxData(
	Bank string,
	UniqueIdentifier string,
	Date string,
	Amount float64,
) *DefaultBankTrxData {
	return &DefaultBankTrxData{
		bank:             Bank,
		UniqueIdentifier: UniqueIdentifier,
		Date:             Date,
		Amount:           Amount,
	}
}

func (b *DefaultBankTrxData) GetBank() string {
	return b.bank
}

func (b *DefaultBankTrxData) GetUniqueIdentifier() string {
	return b.UniqueIdentifier
}

func (b *DefaultBankTrxData) GetDate() string {
	return b.Date
}

func (b *DefaultBankTrxData) GetAmount() float64 {
	return b.Amount
}

type BCABankTrxData struct {
	BCAUniqueIdentifier string
	BCADate             string
	bank                string
	BCAAmount           float64
}

func NewBCABankTrxData(
	Bank string,
	BCAUniqueIdentifier string,
	BCADate string,
	BCAAmount float64,
) *BCABankTrxData {
	return &BCABankTrxData{
		bank:                Bank,
		BCAUniqueIdentifier: BCAUniqueIdentifier,
		BCADate:             BCADate,
		BCAAmount:           BCAAmount,
	}
}

func (c *BCABankTrxData) GetBank() string {
	return c.bank
}

func (c *BCABankTrxData) GetUniqueIdentifier() string {
	return c.BCAUniqueIdentifier
}

func (c *BCABankTrxData) GetDate() string {
	return c.BCADate
}

func (c *BCABankTrxData) GetAmount() float64 {
	return c.BCAAmount
}

type BNIBankTrxData struct {
	BNIUniqueIdentifier string
	BNIDate             string
	bank                string
	BNIAmount           float64
}

func NewBNIBankTrxData(
	Bank string,
	BNIUniqueIdentifier string,
	BNIDate string,
	BNIAmount float64,
) *BNIBankTrxData {
	return &BNIBankTrxData{
		bank:                Bank,
		BNIUniqueIdentifier: BNIUniqueIdentifier,
		BNIDate:             BNIDate,
		BNIAmount:           BNIAmount,
	}
}

func (d *BNIBankTrxData) GetBank() string {
	return d.bank
}

func (d *BNIBankTrxData) GetUniqueIdentifier() string {
	return d.BNIUniqueIdentifier
}

func (d *BNIBankTrxData) GetDate() string {
	return d.BNIDate
}

func (d *BNIBankTrxData) GetAmount() float64 {
	return d.BNIAmount
}
