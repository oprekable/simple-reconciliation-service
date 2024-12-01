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

	"github.com/spf13/afero"

	"github.com/samber/lo"

	"github.com/aaronjan/hunch"
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

func (s *Svc) GenerateReport(ctx context.Context, fs afero.Fs) (returnSummary Summary, err error) {
	ctx = s.comp.Logger.GetLogger().With().Str("component", "Sample Service").Ctx(ctx).Logger().WithContext(s.comp.Logger.GetCtx())

	var trxData []sample.TrxData
	defer func() {
		_ = s.repo.RepoSample.Close()
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			er := s.repo.RepoSample.Pre(
				c,
				s.comp.Config.Reconciliation.ListBank,
				s.comp.Config.Reconciliation.FromDate,
				s.comp.Config.Reconciliation.ToDate,
				s.comp.Config.Reconciliation.TotalData,
				s.comp.Config.Reconciliation.PercentageMatch,
			)

			log.AddErr(c, er)
			log.Msg(c, "[sample.NewSvc] RepoSample.Pre executed")

			return nil, err
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			r, er := s.repo.RepoSample.GetTrx(
				c,
			)

			log.AddErr(c, er)
			log.Msg(c, "[sample.NewSvc] RepoSample.GetTrx executed")

			return r, er
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			trxData = i.([]sample.TrxData)
			er := s.repo.RepoSample.Post(
				c,
			)

			log.AddErr(c, er)
			log.Msg(c, "[sample.NewSvc] RepoSample.Post executed")

			return nil, er
		},
		func(c context.Context, i interface{}) (interface{}, error) {
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
					if data.Type == "DEBIT" {
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
						true,
					)

					log.AddErr(ct, er)
					log.Msg(c, "[sample.NewSvc] save csv file "+returnSummary.FileSystemTrx+" executed")

					return nil, er
				},
			)

			returnSummary.TotalBankTrx = make(map[string]int64)
			returnSummary.FileBankTrx = make(map[string]string)

			for bankName, trxData := range bankTrxData {
				returnSummary.FileBankTrx[bankName] = fmt.Sprintf("%s/%s/%s_%s.csv", s.comp.Config.Reconciliation.BankTRXPath, bankName, bankName, fileNameSuffix)

				switch strings.ToUpper(bankName) {
				case "BCA":
					{
						bd := make([]BCABankTrxData, 0, len(trxData))
						lo.ForEach(trxData, func(data interface{}, _ int) {
							bd = append(bd, data.(BCABankTrxData))
						})

						returnSummary.TotalBankTrx[bankName] = int64(len(bd))

						executor = append(
							executor,
							func(ct context.Context) (interface{}, error) {
								er := csvhelper.StructToCSVFile(
									ct,
									fs,
									returnSummary.FileBankTrx[bankName],
									bd,
									true,
								)

								log.AddErr(ct, er)
								log.Msg(c, "[sample.NewSvc] save csv file "+returnSummary.FileBankTrx[bankName]+" executed")

								return nil, er
							},
						)
					}
				case "BNI":
					{
						bd := make([]BNIBankTrxData, 0, len(trxData))
						lo.ForEach(trxData, func(data interface{}, _ int) {
							bd = append(bd, data.(BNIBankTrxData))
						})

						returnSummary.TotalBankTrx[bankName] = int64(len(bd))

						executor = append(
							executor,
							func(ct context.Context) (interface{}, error) {
								er := csvhelper.StructToCSVFile(
									ct,
									fs,
									returnSummary.FileBankTrx[bankName],
									bd,
									true,
								)

								log.AddErr(ct, er)
								log.Msg(c, "[sample.NewSvc] save csv file "+returnSummary.FileBankTrx[bankName]+" executed")

								return nil, er
							},
						)
					}
				default:
					{
						bd := make([]DefaultBankTrxData, 0, len(trxData))
						lo.ForEach(trxData, func(data interface{}, _ int) {
							bd = append(bd, data.(DefaultBankTrxData))
						})

						returnSummary.TotalBankTrx[bankName] = int64(len(bd))

						executor = append(
							executor,
							func(ct context.Context) (interface{}, error) {
								er := csvhelper.StructToCSVFile(
									ct,
									fs,
									returnSummary.FileBankTrx[bankName],
									bd,
									true,
								)

								log.AddErr(ct, er)
								log.Msg(c, "[sample.NewSvc] save csv file "+returnSummary.FileBankTrx[bankName]+" executed")

								return nil, er
							},
						)
					}
				}
			}

			return hunch.All(
				c,
				executor...,
			)
		},
	)

	return
}
