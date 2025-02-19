package default_bank

import (
	"context"
	"encoding/csv"
	"math"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/_helper"
	"time"
)

type CSVBankTrxData struct {
	DefaultUniqueIdentifier string  `csv:"UniqueIdentifier"`
	DefaultDate             string  `csv:"Date"`
	DefaultBank             string  `csv:"-"`
	DefaultAmount           float64 `csv:"Amount"`
}

func (u *CSVBankTrxData) GetUniqueIdentifier() string {
	return u.DefaultUniqueIdentifier
}

func (u *CSVBankTrxData) GetDate() string {
	return u.DefaultDate
}

func (u *CSVBankTrxData) GetAmount() float64 {
	return u.DefaultAmount
}

func (u *CSVBankTrxData) GetAbsAmount() float64 {
	return math.Abs(u.DefaultAmount)
}

func (u *CSVBankTrxData) GetType() banks.TrxType {
	if u.DefaultAmount <= 0 {
		return banks.DEBIT
	}

	return banks.CREDIT
}

func (u *CSVBankTrxData) GetBank() string {
	return u.DefaultBank
}

func (u *CSVBankTrxData) ToBankTrxData() (returnData *banks.BankTrxData, err error) {
	t, e := time.Parse("2006-01-02", u.DefaultDate)
	if e != nil {
		return nil, e
	}

	return &banks.BankTrxData{
		UniqueIdentifier: u.DefaultUniqueIdentifier,
		Date:             t,
		Type:             u.GetType(),
		Bank:             u.DefaultBank,
		FilePath:         "",
		Amount:           u.GetAbsAmount(),
	}, nil
}

type BankParser struct {
	csvReader    *csv.Reader
	parser       banks.BankParserType
	bank         string
	isHaveHeader bool
}

var _ banks.ReconcileBankData = (*BankParser)(nil)

func NewBankParser(
	bank string,
	csvReader *csv.Reader,
	isHaveHeader bool,
) (*BankParser, error) {
	return &BankParser{
		parser:       banks.DefaultBankParser,
		bank:         bank,
		csvReader:    csvReader,
		isHaveHeader: isHaveHeader,
	}, nil
}

func (d *BankParser) GetParser() banks.BankParserType {
	return d.parser
}

func (d *BankParser) GetBank() string {
	return d.bank
}

func (d *BankParser) ToBankTrxData(ctx context.Context, filePath string) (returnData []*banks.BankTrxData, err error) {
	return _helper.ToBankTrxData(
		ctx,
		filePath,
		d.isHaveHeader,
		d.bank,
		d.csvReader,
		&CSVBankTrxData{},
	)
}
