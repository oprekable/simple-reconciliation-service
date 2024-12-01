package parser

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"math"
	"simple-reconciliation-service/internal/pkg/utils/log"

	loparallel "github.com/samber/lo/parallel"

	"github.com/ulule/deepcopier"

	"github.com/jszwec/csvutil"
)

type BNIBankTrxData struct {
	BNIUniqueIdentifier string  `csv:"BNIUniqueIdentifier"`
	BNIDate             string  `csv:"BNIDate"`
	BNIAmount           float64 `csv:"BNIAmount"`
}

func (u *BNIBankTrxData) UniqueIdentifier() string {
	return u.BNIUniqueIdentifier
}

func (u *BNIBankTrxData) Date() string {
	return u.BNIDate
}

func (u *BNIBankTrxData) Amount() float64 {
	return math.Abs(u.BNIAmount)
}

func (u *BNIBankTrxData) Type() TrxType {
	if u.BNIAmount <= 0 {
		return DEBIT
	}

	return CREDIT
}

type BNIBank struct {
	csvReader *csv.Reader
	parser    BankParser
	bank      string
}

var _ ReconcileBankData = (*BNIBank)(nil)

func NewBNIBank(
	bank string,
	csvReader *csv.Reader,
) (*BNIBank, error) {
	return &BNIBank{
		parser:    BNIBankParser,
		bank:      bank,
		csvReader: csvReader,
	}, nil
}

func (d *BNIBank) GetParser() BankParser {
	return d.parser
}

func (d *BNIBank) GetBank() string {
	return d.bank
}

func (d *BNIBank) ToBankTrxData(ctx context.Context, isHaveHeader bool) (returnData []*BankTrxData, err error) {
	var dec *csvutil.Decoder
	if isHaveHeader {
		dec, err = csvutil.NewDecoder(d.csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, er := csvutil.Header(BNIBankTrxData{}, "csv")
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
		originalData := &BNIBankTrxData{}
		bankTrxData := &BankTrxData{}
		if err := dec.Decode(originalData); err == io.EOF {
			break
		} else if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}

		err := deepcopier.Copy(originalData).To(bankTrxData)
		if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}

		bankTrxData.Bank = d.bank
		returnData = append(returnData, bankTrxData)
	}

	return
}

func (d *BNIBank) ToSql(ctx context.Context, isHaveHeader bool) (returnData string, err error) {
	data, err := d.ToBankTrxData(ctx, isHaveHeader)
	if err != nil {
		log.AddErr(ctx, err)
		return returnData, err
	}

	var buffer bytes.Buffer

	loparallel.ForEach(data, func(d *BankTrxData, _ int) {
		buffer.WriteString(d.UniqueIdentifier)
	})

	returnData = buffer.String()
	return
}
