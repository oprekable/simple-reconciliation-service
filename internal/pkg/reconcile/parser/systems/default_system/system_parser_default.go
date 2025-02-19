package default_system

import (
	"context"
	"encoding/csv"
	"io"
	"math"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"time"

	"github.com/jszwec/csvutil"
)

type CSVSystemTrxData struct {
	TrxID           string  `csv:"TrxID"`
	TransactionTime string  `csv:"TransactionTime"`
	Type            string  `csv:"Type"`
	Amount          float64 `csv:"Amount"`
}

func (u *CSVSystemTrxData) GetTrxID() string {
	return u.TrxID
}

func (u *CSVSystemTrxData) GetTransactionTime() string {
	return u.TransactionTime
}

func (u *CSVSystemTrxData) GetAmount() float64 {
	return math.Abs(u.Amount)
}

func (u *CSVSystemTrxData) GetType() systems.TrxType {
	return systems.TrxType(u.Type)
}
func (u *CSVSystemTrxData) ToSystemTrxData() (returnData *systems.SystemTrxData, err error) {
	t, e := time.Parse("2006-01-02 15:04:05", u.TransactionTime)
	if e != nil {
		return nil, e
	}

	return &systems.SystemTrxData{
		TrxID:           u.TrxID,
		TransactionTime: t,
		Type:            systems.TrxType(u.Type),
		FilePath:        "",
		Amount:          u.Amount,
	}, nil
}

type SystemParser struct {
	csvReader    *csv.Reader
	parser       systems.SystemParserType
	isHaveHeader bool
}

var _ systems.ReconcileSystemData = (*SystemParser)(nil)

func NewSystemParser(
	csvReader *csv.Reader,
	isHaveHeader bool,
) (*SystemParser, error) {
	return &SystemParser{
		parser:       systems.DefaultSystemParser,
		csvReader:    csvReader,
		isHaveHeader: isHaveHeader,
	}, nil
}

func (d *SystemParser) ToSystemTrxData(ctx context.Context, filePath string) (returnData []*systems.SystemTrxData, err error) {
	var dec *csvutil.Decoder
	if d.isHaveHeader {
		dec, err = csvutil.NewDecoder(d.csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, er := csvutil.Header(CSVSystemTrxData{}, "csv")
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
		originalData := &CSVSystemTrxData{}
		if err := dec.Decode(originalData); err == io.EOF {
			break
		} else if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}

		systemTrxData, er := originalData.ToSystemTrxData()
		if er != nil {
			log.AddErr(ctx, er)
			continue
		}

		systemTrxData.FilePath = filePath
		returnData = append(returnData, systemTrxData)
	}

	return
}
