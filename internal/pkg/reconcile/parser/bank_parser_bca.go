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

type BCABankTrxData struct {
	BCAUniqueIdentifier string  `csv:"BCAUniqueIdentifier"`
	BCADate             string  `csv:"BCADate"`
	BCAAmount           float64 `csv:"BCAAmount"`
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
	csvReader    *csv.Reader
	parser       BankParser
	bank         string
	isHaveHeader bool
}

var _ ReconcileBankData = (*BCABank)(nil)

func NewBCABank(
	bank string,
	csvReader *csv.Reader,
	isHaveHeader bool,
) (*BCABank, error) {
	return &BCABank{
		parser:       BCABankParser,
		bank:         bank,
		csvReader:    csvReader,
		isHaveHeader: isHaveHeader,
	}, nil
}

func (d *BCABank) GetParser() BankParser {
	return d.parser
}

func (d *BCABank) GetBank() string {
	return d.bank
}

func (d *BCABank) ToBankTrxData(ctx context.Context, filePath string) (returnData []*BankTrxData, err error) {
	var dec *csvutil.Decoder
	if d.isHaveHeader {
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

		bankTrxData.Bank = d.bank
		bankTrxData.FilePath = filePath
		returnData = append(returnData, bankTrxData)
	}

	return
}

func (d *BCABank) ToSql(ctx context.Context, filePath string, sqlPattern string) (returnData string, err error) {
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
