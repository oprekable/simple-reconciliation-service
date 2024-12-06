package parser

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"simple-reconciliation-service/internal/pkg/utils/log"

	"github.com/jszwec/csvutil"
	"github.com/samber/lo"
	"github.com/ulule/deepcopier"
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
	csvReader    *csv.Reader
	parser       BankParser
	bank         string
	isHaveHeader bool
}

var _ ReconcileBankData = (*BNIBank)(nil)

func NewBNIBank(
	bank string,
	csvReader *csv.Reader,
	isHaveHeader bool,
) (*BNIBank, error) {
	return &BNIBank{
		parser:       BNIBankParser,
		bank:         bank,
		csvReader:    csvReader,
		isHaveHeader: isHaveHeader,
	}, nil
}

func (d *BNIBank) GetParser() BankParser {
	return d.parser
}

func (d *BNIBank) GetBank() string {
	return d.bank
}

func (d *BNIBank) ToBankTrxData(ctx context.Context, filePath string) (returnData []*BankTrxData, err error) {
	var dec *csvutil.Decoder
	if d.isHaveHeader {
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
		bankTrxData.FilePath = filePath
		returnData = append(returnData, bankTrxData)
	}

	return
}

func (d *BNIBank) ToSql(ctx context.Context, filePath string, sqlPattern string) (returnData string, err error) {
	data, err := d.ToBankTrxData(ctx, filePath)
	if err != nil {
		log.AddErr(ctx, err)
		return returnData, err
	}

	var buffer bytes.Buffer

	lo.ForEach(data, func(d *BankTrxData, _ int) {
		buffer.WriteString(
			fmt.Sprintf(
				sqlPattern,
				d.UniqueIdentifier,
				d.Date,
				d.Bank,
				d.Type,
				d.Amount,
				d.FilePath,
			),
		)
	})

	returnData = buffer.String()
	return
}
