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

type DefaultBankTrxData struct {
	DefaultUniqueIdentifier string  `csv:"uniqueIdentifier"`
	DefaultDate             string  `csv:"date"`
	DefaultAmount           float64 `csv:"amount"`
}

func (u *DefaultBankTrxData) UniqueIdentifier() string {
	return u.DefaultUniqueIdentifier
}

func (u *DefaultBankTrxData) Date() string {
	return u.DefaultDate
}

func (u *DefaultBankTrxData) Amount() float64 {
	return math.Abs(u.DefaultAmount)
}

func (u *DefaultBankTrxData) Type() TrxType {
	if u.DefaultAmount <= 0 {
		return DEBIT
	}

	return CREDIT
}

type DefaultBank struct {
	csvReader *csv.Reader
	parser    Parser
	bank      string
}

var _ ReconcileData = (*DefaultBank)(nil)

func NewDefaultBank(
	bank string,
	csvReader *csv.Reader,
) (*DefaultBank, error) {
	return &DefaultBank{
		parser:    DEFAULT,
		bank:      bank,
		csvReader: csvReader,
	}, nil
}

func (d *DefaultBank) GetParser() Parser {
	return d.parser
}

func (d *DefaultBank) GetBank() string {
	return d.bank
}

func (d *DefaultBank) ToBankTrxData(ctx context.Context, isHaveHeader bool) (returnData []*BankTrxData, err error) {
	var dec *csvutil.Decoder
	if isHaveHeader {
		dec, err = csvutil.NewDecoder(d.csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, er := csvutil.Header(DefaultBankTrxData{}, "csv")
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
		originalData := &DefaultBankTrxData{}
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

func (d *DefaultBank) ToSql(ctx context.Context, isHaveHeader bool) (returnData string, err error) {
	data, err := d.ToBankTrxData(ctx, isHaveHeader)
	if err != nil {
		log.AddErr(ctx, err)
		return returnData, err
	}

	var buffer bytes.Buffer

	loparallel.ForEach(data, func(d *BankTrxData, _ int) {
		buffer.WriteString(d.UniqueIdentifier + "\n")
	})

	returnData = buffer.String()
	return
}
