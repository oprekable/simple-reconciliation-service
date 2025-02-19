package bni

import (
	"context"
	"encoding/csv"
	"math"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/_helper"
	"time"
)

type CSVBankTrxData struct {
	BNIUniqueIdentifier string  `csv:"BNIUniqueIdentifier"`
	BNIDate             string  `csv:"BNIDate"`
	BNIBank             string  `csv:"-"`
	BNIAmount           float64 `csv:"BNIAmount"`
}

func (u *CSVBankTrxData) GetUniqueIdentifier() string {
	return u.BNIUniqueIdentifier
}

func (u *CSVBankTrxData) GetDate() string {
	return u.BNIDate
}

func (u *CSVBankTrxData) GetAmount() float64 {
	return u.BNIAmount
}

func (u *CSVBankTrxData) GetAbsAmount() float64 {
	return math.Abs(u.BNIAmount)
}

func (u *CSVBankTrxData) GetType() banks.TrxType {
	if u.BNIAmount <= 0 {
		return banks.DEBIT
	}

	return banks.CREDIT
}

func (u *CSVBankTrxData) GetBank() string {
	return u.BNIBank
}

func (u *CSVBankTrxData) ToBankTrxData() (returnData *banks.BankTrxData, err error) {
	t, e := time.Parse("2006-01-02", u.BNIDate)
	if e != nil {
		return nil, e
	}

	return &banks.BankTrxData{
		UniqueIdentifier: u.BNIUniqueIdentifier,
		Date:             t,
		Type:             u.GetType(),
		Bank:             u.BNIBank,
		FilePath:         "",
		Amount:           u.GetAbsAmount(),
	}, nil
}

type BankParser struct {
	csvReader    *csv.Reader
	banks        banks.BankParserType
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
		banks:        banks.BNIBankParser,
		csvReader:    csvReader,
		isHaveHeader: isHaveHeader,
		bank:         bank,
	}, nil
}

func (d *BankParser) GetBank() string {
	return d.bank
}

func (d *BankParser) GetParser() banks.BankParserType {
	return d.banks
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
