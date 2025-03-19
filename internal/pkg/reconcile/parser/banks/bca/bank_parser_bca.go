package bca

import (
	"context"
	"encoding/csv"
	"errors"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/bca/entity"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/helper"
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
		parser:       banks.BCABankParser,
		bank:         bank,
		isHaveHeader: isHaveHeader,
		csvReader:    csvReader,
	}, nil
}

func (d *BankParser) ToBankTrxData(ctx context.Context, filePath string) (returnData []*banks.BankTrxData, err error) {
	return helper.ToBankTrxData(
		ctx,
		filePath,
		d.isHaveHeader,
		d.bank,
		d.csvReader,
		&entity.CSVBankTrxData{},
	)
}

func (d *BankParser) GetParser() banks.BankParserType {
	return d.parser
}

func (d *BankParser) GetBank() string {
	return d.bank
}
