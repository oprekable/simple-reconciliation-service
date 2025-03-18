package bni

import (
	"context"
	"encoding/csv"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/_helper"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/bni/entity"
)

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
		&entity.CSVBankTrxData{},
	)
}
