package bca

import (
	"context"
	"encoding/csv"
	"io"
	"math"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/utils/log"

	"github.com/jszwec/csvutil"
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
	var dec *csvutil.Decoder
	if d.isHaveHeader {
		dec, err = csvutil.NewDecoder(d.csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, er := csvutil.Header(CSVBankTrxData{}, "csv")
		if er != nil {
			log.AddErr(ctx, er)
			return nil, er
		}

		dec, er = csvutil.NewDecoder(d.csvReader, header...)
		if er != nil {
			log.AddErr(ctx, er)
			return nil, er
		}
	}

	for {
		originalData := &CSVBankTrxData{}
		if err := dec.Decode(originalData); err == io.EOF {
			break
		} else if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}

		bankTrxData := originalData.ToBankTrxData()
		bankTrxData.Bank = d.bank
		bankTrxData.FilePath = filePath
		returnData = append(returnData, bankTrxData)
	}

	return
}
