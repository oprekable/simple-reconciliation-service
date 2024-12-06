package sample

import (
	"context"
	"fmt"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/repository/sample"
	"simple-reconciliation-service/internal/pkg/utils/csvhelper"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strconv"
	"strings"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/samber/lo"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/afero"
)

type Svc struct {
	comp *component.Components
	repo *repository.Repositories
}

var _ Service = (*Svc)(nil)

func NewSvc(
	comp *component.Components,
	repo *repository.Repositories,
) *Svc {
	return &Svc{
		comp: comp,
		repo: repo,
	}
}

func (s *Svc) GenerateSample(ctx context.Context, fs afero.Fs, bar *progressbar.ProgressBar, isDeleteDirectory bool) (returnSummary Summary, err error) {
	ctx = s.comp.Logger.GetLogger().With().Str("component", "Sample Service").Ctx(ctx).Logger().WithContext(s.comp.Logger.GetCtx())

	var trxData []sample.TrxData
	defer func() {
		_ = s.repo.RepoSample.Close()
		if bar != nil {
			_ = bar.Clear()
		}
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][1/5] Pre Process Generate Sample...")
			}

			if isDeleteDirectory {
				if er := csvhelper.DeleteDirectory(c, fs, s.comp.Config.Reconciliation.SystemTRXPath); er != nil {
					log.Err(c, "[sample.NewSvc] DeleteDirectory SystemTRXPath", er)
					return nil, er
				}

				if er := csvhelper.DeleteDirectory(c, fs, s.comp.Config.Reconciliation.BankTRXPath); er != nil {
					log.Err(c, "[sample.NewSvc] DeleteDirectory BankTRXPath", er)
					return nil, er
				}
			}

			er := s.repo.RepoSample.Pre(
				c,
				s.comp.Config.Reconciliation.ListBank,
				s.comp.Config.Reconciliation.FromDate,
				s.comp.Config.Reconciliation.ToDate,
				s.comp.Config.Reconciliation.TotalData,
				s.comp.Config.Reconciliation.PercentageMatch,
			)

			log.Err(c, "[sample.NewSvc] RepoSample.Pre executed", er)
			return nil, err
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][2/5] Populate Trx Data...")
			}

			r, er := s.repo.RepoSample.GetTrx(
				c,
			)

			log.Err(c, "[sample.NewSvc] RepoSample.GetTrx executed", er)
			return r, er
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][3/5] Post Process Generate Sample...")
			}

			if i != nil {
				trxData = i.([]sample.TrxData)
			}

			er := s.repo.RepoSample.Post(
				c,
			)

			log.Err(c, "[sample.NewSvc] RepoSample.Post executed", er)

			return nil, er
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][4/5] Parse Sample Data...")
			}

			systemTrxData := make([]SystemTrxData, 0, len(trxData))
			bankTrxData := make(map[string][]interface{})

			lo.ForEach(trxData, func(data sample.TrxData, _ int) {
				if data.IsSystemTrx {
					item := SystemTrxData{
						TrxID:           data.TrxID,
						Type:            data.Type,
						TransactionTime: data.TransactionTime,
						Amount:          data.Amount,
					}

					systemTrxData = append(systemTrxData, item)
				}

				if data.IsBankTrx {
					bank := strings.ToLower(data.Bank)
					if _, ok := bankTrxData[bank]; !ok {
						bankTrxData[bank] = make([]interface{}, 0, len(trxData))
					}

					multiplier := float64(1)
					if data.Type == DEBIT {
						multiplier = float64(-1)
					}

					switch strings.ToUpper(bank) {
					case "BCA":
						{
							item := BCABankTrxData{
								BCAUniqueIdentifier: data.UniqueIdentifier,
								BCADate:             data.Date,
								BCAAmount:           data.Amount * multiplier,
							}

							bankTrxData[bank] = append(bankTrxData[bank], item)
						}
					case "BNI":
						{
							item := BNIBankTrxData{
								BNIUniqueIdentifier: data.UniqueIdentifier,
								BNIDate:             data.Date,
								BNIAmount:           data.Amount * multiplier,
							}

							bankTrxData[bank] = append(bankTrxData[bank], item)
						}
					default:
						{
							item := DefaultBankTrxData{
								UniqueIdentifier: data.UniqueIdentifier,
								Date:             data.Date,
								Amount:           data.Amount * multiplier,
							}

							bankTrxData[bank] = append(bankTrxData[bank], item)
						}
					}
				}
			})

			log.Msg(c, "[sample.NewSvc] populate systemTrxData & bankTrxData executed")

			if bar != nil {
				bar.Describe("[cyan][5/5] Export Sample Data to CSV files...")
			}

			fileNameSuffix := strconv.FormatInt(time.Now().Unix(), 10)
			returnSummary.FileSystemTrx = fmt.Sprintf("%s/%s.csv", s.comp.Config.Reconciliation.SystemTRXPath, fileNameSuffix)
			returnSummary.TotalSystemTrx = int64(len(systemTrxData))

			executor := make([]hunch.Executable, 0, len(bankTrxData)+1)
			executor = append(
				executor,
				func(ct context.Context) (interface{}, error) {
					er := csvhelper.StructToCSVFile(
						ct,
						fs,
						returnSummary.FileSystemTrx,
						systemTrxData,
						isDeleteDirectory,
					)

					log.Err(c, "[sample.NewSvc] save csv file "+returnSummary.FileSystemTrx+" executed", er)

					return nil, er
				},
			)

			returnSummary.TotalBankTrx = make(map[string]int64)
			returnSummary.FileBankTrx = make(map[string]string)

			for bankName, trxData := range bankTrxData {
				returnSummary.FileBankTrx[bankName] = fmt.Sprintf("%s/%s/%s_%s.csv", s.comp.Config.Reconciliation.BankTRXPath, bankName, bankName, fileNameSuffix)
				totalBankTrx, exec := s.appendExecutor(
					fs,
					returnSummary.FileBankTrx[bankName],
					trxData,
					isDeleteDirectory,
				)

				returnSummary.TotalBankTrx[bankName] = totalBankTrx
				executor = append(executor, exec)
			}

			return hunch.All(
				c,
				executor...,
			)
		},
	)

	return
}

func (s *Svc) appendExecutor(fs afero.Fs, filePath string, trxDataSlice interface{}, isDeleteDirectory bool) (totalData int64, executor hunch.Executable) {
	switch value := trxDataSlice.(type) {
	case []interface{}:
		{
			switch value[0].(type) {
			case BCABankTrxData:
				{
					bd := make([]BCABankTrxData, 0, len(value))
					lo.ForEach(trxDataSlice.([]interface{}), func(data interface{}, _ int) {
						bd = append(bd, data.(BCABankTrxData))
					})
					totalData = int64(len(bd))
					executor = func(ct context.Context) (interface{}, error) {
						er := csvhelper.StructToCSVFile(
							ct,
							fs,
							filePath,
							bd,
							isDeleteDirectory,
						)

						log.Err(ct, "[sample.NewSvc] save csv file "+filePath+" executed", er)
						return nil, er
					}
				}
			case BNIBankTrxData:
				{
					bd := make([]BNIBankTrxData, 0, len(value))
					lo.ForEach(trxDataSlice.([]interface{}), func(data interface{}, _ int) {
						bd = append(bd, data.(BNIBankTrxData))
					})
					totalData = int64(len(bd))
					executor = func(ct context.Context) (interface{}, error) {
						er := csvhelper.StructToCSVFile(
							ct,
							fs,
							filePath,
							bd,
							isDeleteDirectory,
						)

						log.Err(ct, "[sample.NewSvc] save csv file "+filePath+" executed", er)
						return nil, er
					}
				}
			default:
				{
					bd := make([]DefaultBankTrxData, 0, len(value))
					lo.ForEach(trxDataSlice.([]interface{}), func(data interface{}, _ int) {
						bd = append(bd, data.(DefaultBankTrxData))
					})
					totalData = int64(len(bd))
					executor = func(ct context.Context) (interface{}, error) {
						er := csvhelper.StructToCSVFile(
							ct,
							fs,
							filePath,
							bd,
							isDeleteDirectory,
						)

						log.Err(ct, "[sample.NewSvc] save csv file "+filePath+" executed", er)
						return nil, er
					}
				}
			}
		}
	}

	return
}
