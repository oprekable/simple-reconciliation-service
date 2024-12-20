package bca

import (
	"context"
	"encoding/csv"
	"math"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/_helper"
)

type CSVBankTrxData struct {
	BCAUniqueIdentifier string  `csv:"BCAUniqueIdentifier"`
	BCADate             string  `csv:"BCADate"`
	BCABank             string  `csv:"-"`
	BCAAmount           float64 `csv:"BCAAmount"`
}

func (u *CSVBankTrxData) GetUniqueIdentifier() string {
	return u.BCAUniqueIdentifier
}

func (u *CSVBankTrxData) GetDate() string {
	return u.BCADate
}

func (u *CSVBankTrxData) GetAmount() float64 {
	return u.BCAAmount
}

func (u *CSVBankTrxData) GetAbsAmount() float64 {
	return math.Abs(u.BCAAmount)
}

func (u *CSVBankTrxData) GetType() banks.TrxType {
	if u.BCAAmount <= 0 {
		return banks.DEBIT
	}

	return banks.CREDIT
}

func (u *CSVBankTrxData) GetBank() string {
	return u.BCABank
}

func (u *CSVBankTrxData) ToBankTrxData() *banks.BankTrxData {
	return &banks.BankTrxData{
		UniqueIdentifier: u.BCAUniqueIdentifier,
		Date:             u.BCADate,
		Type:             u.GetType(),
		Bank:             u.BCABank,
		FilePath:         "",
		Amount:           u.GetAbsAmount(),
	}
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
		parser:       banks.BCABankParser,
		bank:         bank,
		isHaveHeader: isHaveHeader,
		csvReader:    csvReader,
	}, nil
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

func (d *BankParser) GetParser() banks.BankParserType {
	return d.parser
}

func (d *BankParser) GetBank() string {
	return d.bank
}
