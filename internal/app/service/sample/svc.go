package sample

import (
	"context"
	"fmt"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/repository/sample"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/bca"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/bni"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/default_bank"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems/default_system"
	"simple-reconciliation-service/internal/pkg/utils/csvhelper"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"simple-reconciliation-service/internal/pkg/utils/progressbarhelper"
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

func (s *Svc) deleteDirectorySystemTrxBankTrx(ctx context.Context, fs afero.Fs, isDeleteDirectory bool) (err error) {
	if !isDeleteDirectory {
		return
	}

	if er := csvhelper.DeleteDirectory(ctx, fs, s.comp.Config.Data.Reconciliation.SystemTRXPath); er != nil {
		log.Err(ctx, "[sample.NewSvc] DeleteDirectory SystemTRXPath", er)
		return er
	}

	if er := csvhelper.DeleteDirectory(ctx, fs, s.comp.Config.Data.Reconciliation.BankTRXPath); er != nil {
		log.Err(ctx, "[sample.NewSvc] DeleteDirectory BankTRXPath", er)
		return er
	}

	return
}

func (s *Svc) parse(data sample.TrxData) (systemTrxData systems.SystemTrxDataInterface, bankTrxData banks.BankTrxDataInterface) {
	if data.IsSystemTrx {
		systemTrxData = &default_system.CSVSystemTrxData{
			TrxID:           data.TrxID,
			TransactionTime: data.TransactionTime,
			Type:            data.Type,
			Amount:          data.Amount,
		}
	}

	if data.IsBankTrx || (!data.IsBankTrx && !data.IsSystemTrx) {
		bank := strings.ToLower(data.Bank)
		multiplier := float64(1)
		if data.Type == DEBIT {
			multiplier = float64(-1)
		}

		switch strings.ToUpper(bank) {
		case "BCA":
			{
				bankTrxData = &bca.CSVBankTrxData{
					UniqueIdentifier: data.UniqueIdentifier,
					Date:             data.Date,
					Amount:           data.Amount * multiplier,
					Bank:             bank,
				}
			}
		case "BNI":
			{
				bankTrxData = &bni.CSVBankTrxData{
					UniqueIdentifier: data.UniqueIdentifier,
					Date:             data.Date,
					Amount:           data.Amount * multiplier,
					Bank:             bank,
				}
			}
		default:
			{
				bankTrxData = &default_bank.CSVBankTrxData{
					UniqueIdentifier: data.UniqueIdentifier,
					Date:             data.Date,
					Amount:           data.Amount * multiplier,
					Bank:             bank,
				}
			}
		}
	}

	return
}

func (s *Svc) GenerateSample(ctx context.Context, fs afero.Fs, bar *progressbar.ProgressBar, isDeleteDirectory bool) (returnSummary Summary, err error) {
	ctx = s.comp.Logger.GetLogger().With().Str("component", "Sample Service").Ctx(ctx).Logger().WithContext(s.comp.Logger.GetCtx())

	var trxData []sample.TrxData
	defer func() {
		_ = s.repo.RepoSample.Close()
		progressbarhelper.BarClear(bar)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return nil, s.deleteDirectorySystemTrxBankTrx(c, fs, isDeleteDirectory)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][1/5] Pre Process Generate Sample...")
			defer func() {
				log.Err(c, "[sample.NewSvc] RepoSample.Pre executed", e)
			}()

			e = s.repo.RepoSample.Pre(
				c,
				s.comp.Config.Data.Reconciliation.ListBank,
				s.comp.Config.Data.Reconciliation.FromDate,
				s.comp.Config.Data.Reconciliation.ToDate,
				s.comp.Config.Data.Reconciliation.TotalData,
				s.comp.Config.Data.Reconciliation.PercentageMatch,
			)

			return nil, e
		},
		func(c context.Context, _ interface{}) (r interface{}, er error) {
			progressbarhelper.BarDescribe(bar, "[cyan][2/5] Populate Trx Data...")

			trxData, er = s.repo.RepoSample.GetTrx(
				c,
			)

			log.Err(c, "[sample.NewSvc] RepoSample.GetTrx executed", er)
			return nil, er
		},
		func(c context.Context, i interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][3/5] Post Process Generate Sample...")
			if !s.comp.Config.Data.IsDebug {
				e = s.repo.RepoSample.Post(
					c,
				)
			}

			log.Err(c, "[sample.NewSvc] RepoSample.Post executed", e)

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			progressbarhelper.BarDescribe(bar, "[cyan][4/5] Parse Sample Data...")
			systemTrxData := make([]systems.SystemTrxDataInterface, 0, len(trxData))
			bankTrxData := make(map[string][]banks.BankTrxDataInterface)

			lo.ForEach(trxData, func(data sample.TrxData, _ int) {
				systemTrx, bankTrx := s.parse(data)
				if systemTrx != nil {
					systemTrxData = append(systemTrxData, systemTrx)
				}

				if bankTrx != nil {
					bankTrxData[bankTrx.GetBank()] = append(bankTrxData[bankTrx.GetBank()], bankTrx)
				}
			})

			trxData = nil

			log.Msg(c, "[sample.NewSvc] populate systemTrxData & bankTrxData executed")
			progressbarhelper.BarDescribe(bar, "[cyan][5/5] Export Sample Data to CSV files...")

			lengthSystemTrxData := len(systemTrxData)
			fileNameSuffix := strconv.FormatInt(time.Now().Unix(), 10)
			returnSummary.FileSystemTrx = fmt.Sprintf("%s/%s.csv", s.comp.Config.Data.Reconciliation.SystemTRXPath, fileNameSuffix)
			returnSummary.TotalSystemTrx = int64(lengthSystemTrxData)

			sd := make([]*default_system.CSVSystemTrxData, 0, lengthSystemTrxData)

			lo.ForEach(systemTrxData, func(data systems.SystemTrxDataInterface, _ int) {
				sd = append(sd, data.(*default_system.CSVSystemTrxData))
			})

			systemTrxData = nil

			executor := make([]hunch.Executable, 0, len(bankTrxData)+1)
			executor = append(
				executor,
				func(ct context.Context) (interface{}, error) {
					er := csvhelper.StructToCSVFile(
						ct,
						fs,
						returnSummary.FileSystemTrx,
						sd,
						isDeleteDirectory,
					)

					log.Err(c, "[sample.NewSvc] save csv file "+returnSummary.FileSystemTrx+" executed", er)

					return nil, er
				},
			)

			returnSummary.TotalBankTrx = make(map[string]int64)
			returnSummary.FileBankTrx = make(map[string]string)

			for bankName, bankTrx := range bankTrxData {
				returnSummary.FileBankTrx[bankName] = fmt.Sprintf("%s/%s/%s_%s.csv", s.comp.Config.Data.Reconciliation.BankTRXPath, bankName, bankName, fileNameSuffix)
				totalBankTrx, exec := s.appendExecutor(
					fs,
					returnSummary.FileBankTrx[bankName],
					bankTrx,
					isDeleteDirectory,
				)

				if exec == nil || totalBankTrx == 0 {
					continue
				}

				returnSummary.TotalBankTrx[bankName] = totalBankTrx
				executor = append(executor, exec)
			}

			bankTrxData = nil

			return hunch.All(
				c,
				executor...,
			)
		},
	)

	return
}

func (s *Svc) appendExecutor(fs afero.Fs, filePath string, trxDataSlice []banks.BankTrxDataInterface, isDeleteDirectory bool) (totalData int64, executor hunch.Executable) {
	if len(trxDataSlice) == 0 {
		return 0, nil
	}

	formatText := "[sample.NewSvc] save csv file %s executed"

	switch trxDataSlice[0].(type) {
	case *bca.CSVBankTrxData:
		{
			bd := make([]*bca.CSVBankTrxData, 0, len(trxDataSlice))

			lo.ForEach(trxDataSlice, func(data banks.BankTrxDataInterface, _ int) {
				bd = append(bd, data.(*bca.CSVBankTrxData))
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

				log.Err(ct, fmt.Sprintf(formatText, filePath), er)
				return nil, er
			}
		}
	case *bni.CSVBankTrxData:
		{
			bd := make([]*bni.CSVBankTrxData, 0, len(trxDataSlice))
			lo.ForEach(trxDataSlice, func(data banks.BankTrxDataInterface, _ int) {
				bd = append(bd, data.(*bni.CSVBankTrxData))
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

				log.Err(ct, fmt.Sprintf(formatText, filePath), er)
				return nil, er
			}
		}
	default:
		{
			bd := make([]*default_bank.CSVBankTrxData, 0, len(trxDataSlice))
			lo.ForEach(trxDataSlice, func(data banks.BankTrxDataInterface, _ int) {
				bd = append(bd, data.(*default_bank.CSVBankTrxData))
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

				log.Err(ct, fmt.Sprintf(formatText, filePath), er)
				return nil, er
			}
		}
	}

	return
}
