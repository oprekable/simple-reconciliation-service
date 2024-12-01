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

type BCABankTrxData struct {
	BCAUniqueIdentifier string  `csv:"uniqueIdentifier"`
	BCADate             string  `csv:"date"`
	BCAAmount           float64 `csv:"amount"`
}

func (u *BCABankTrxData) UniqueIdentifier() string {
	return u.BCAUniqueIdentifier
}

func (u *BCABankTrxData) Date() string {
	return u.BCADate
}

func (u *BCABankTrxData) Amount() float64 {
	return math.Abs(u.BCAAmount)
}

func (u *BCABankTrxData) Type() TrxType {
	if u.BCAAmount <= 0 {
		return DEBIT
	}

	return CREDIT
}

type BCABank struct {
	csvReader *csv.Reader
	parser    Parser
	bank      string
}

var _ ReconcileData = (*BCABank)(nil)

func NewBCABank(
	bank string,
	csvReader *csv.Reader,
) (*BCABank, error) {
	return &BCABank{
		parser:    BCA,
		bank:      bank,
		csvReader: csvReader,
	}, nil
}

func (d *BCABank) GetParser() Parser {
	return d.parser
}

func (d *BCABank) GetBank() string {
	return d.bank
}

func (d *BCABank) ToBankTrxData(ctx context.Context, isHaveHeader bool) (returnData []*BankTrxData, err error) {
	var dec *csvutil.Decoder
	if isHaveHeader {
		dec, err = csvutil.NewDecoder(d.csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, er := csvutil.Header(BCABankTrxData{}, "csv")
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
		originalData := &BCABankTrxData{}
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

		returnData = append(returnData, bankTrxData)
	}

	return
}

func (d *BCABank) ToSql(ctx context.Context, isHaveHeader bool) (returnData string, err error) {
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
