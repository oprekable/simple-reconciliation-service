package _helper

import (
	"context"
	"encoding/csv"
	"github.com/jszwec/csvutil"
	"io"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/utils/log"
)

func ToBankTrxData(ctx context.Context, filePath string, isHaveHeader bool, bank string, csvReader *csv.Reader, originalData banks.BankTrxDataInterface) (returnData []*banks.BankTrxData, err error) {
	var dec *csvutil.Decoder
	if isHaveHeader {
		dec, err = csvutil.NewDecoder(csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, er := csvutil.Header(originalData, "csv")
		if er != nil {
			log.AddErr(ctx, er)
			return nil, er
		}

		dec, er = csvutil.NewDecoder(csvReader, header...)
		if er != nil {
			log.AddErr(ctx, er)
			return nil, er
		}
	}

	for {
		if err := dec.Decode(originalData); err == io.EOF {
			break
		} else if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}

		bankTrxData := originalData.ToBankTrxData()
		bankTrxData.Bank = bank
		bankTrxData.FilePath = filePath
		bankTrxData.Type = originalData.GetType()
		returnData = append(returnData, bankTrxData)
	}

	return
}
