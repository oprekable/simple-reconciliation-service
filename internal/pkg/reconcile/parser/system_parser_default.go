package parser

import (
	"context"
	"encoding/csv"
	"io"
	"math"
	"simple-reconciliation-service/internal/pkg/utils/log"

	"github.com/ulule/deepcopier"

	"github.com/jszwec/csvutil"
)

type DefaultSystemTrxData struct {
	DefaultTrxID           string  `csv:"TrxID"`
	DefaultTransactionTime string  `csv:"TransactionTime"`
	DefaultType            string  `csv:"Type"`
	DefaultAmount          float64 `csv:"Amount"`
}

func (u *DefaultSystemTrxData) TrxID() string {
	return u.DefaultTrxID
}

func (u *DefaultSystemTrxData) TransactionTime() string {
	return u.DefaultTransactionTime
}

func (u *DefaultSystemTrxData) Amount() float64 {
	return math.Abs(u.DefaultAmount)
}

func (u *DefaultSystemTrxData) Type() TrxType {
	if u.DefaultAmount <= 0 {
		return DEBIT
	}

	return CREDIT
}

type DefaultSystem struct {
	csvReader *csv.Reader
	parser    SystemParser
}

var _ ReconcileSystemData = (*DefaultSystem)(nil)

func NewDefaultSystem(
	csvReader *csv.Reader,
) (*DefaultSystem, error) {
	return &DefaultSystem{
		parser:    DEFAULT_SYSTEM,
		csvReader: csvReader,
	}, nil
}

func (d *DefaultSystem) GetParser() SystemParser {
	return d.parser
}

func (d *DefaultSystem) ToSystemTrxData(ctx context.Context, isHaveHeader bool) (returnData []*SystemTrxData, err error) {
	var dec *csvutil.Decoder
	if isHaveHeader {
		dec, err = csvutil.NewDecoder(d.csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, er := csvutil.Header(DefaultSystemTrxData{}, "csv")
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
		originalData := &DefaultSystemTrxData{}
		systemTrxData := &SystemTrxData{}
		if err := dec.Decode(originalData); err == io.EOF {
			break
		} else if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}

		err := deepcopier.Copy(originalData).To(systemTrxData)
		if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}

		returnData = append(returnData, systemTrxData)
	}

	return
}
