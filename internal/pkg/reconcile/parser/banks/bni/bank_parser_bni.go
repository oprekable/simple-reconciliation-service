package bni

import (
	"context"
	"encoding/csv"
	"errors"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/_helper"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/bni/entity"
)

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
	if csvReader == nil {
		return nil, errors.New("csvReader or dataStruct is nil")
	}

	return &BankParser{
		parser:       banks.BNIBankParser,
		csvReader:    csvReader,
		isHaveHeader: isHaveHeader,
		bank:         bank,
	}, nil
}

func (d *BankParser) GetBank() string {
	return d.bank
}

func (d *BankParser) GetParser() banks.BankParserType {
	return d.parser
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
