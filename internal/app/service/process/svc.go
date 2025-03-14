package process

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/repository/process"
	"simple-reconciliation-service/internal/pkg/reconcile/parser"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/bca"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/bni"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks/default_bank"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems/default_system"
	"simple-reconciliation-service/internal/pkg/utils/csvhelper"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"simple-reconciliation-service/internal/pkg/utils/progressbarhelper"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo/parallel"

	"github.com/ulule/deepcopier"

	"github.com/samber/lo"

	"github.com/aaronjan/hunch"
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

func (s *Svc) parseSystemTrxFile(ctx context.Context, afs afero.Fs, filePath string) (returnData []*systems.SystemTrxData, err error) {
	var f afero.File
	f, err = afs.Open(filePath)

	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()

	var systemParser *default_system.SystemParser
	systemParser, err = default_system.NewSystemParser(
		csv.NewReader(f),
		true,
	)

	log.Err(ctx, "[process.NewSvc] parseSystemTrxFile - '"+filePath+"'", err)

	if err != nil {
		return
	}

	returnData, err = systemParser.ToSystemTrxData(ctx, filePath)
	log.Err(ctx, "[process.NewSvc] parseSystemTrxFile parse.ToSystemTrxData executed", err)

	return
}

func (s *Svc) parseSystemTrxFiles(ctx context.Context, afs afero.Fs) (returnData []*systems.SystemTrxData, err error) {
	var filePathSystemTrx []string
	cleanPath := filepath.Clean(s.comp.Config.Data.Reconciliation.SystemTRXPath)
	err = afero.Walk(afs, cleanPath, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".csv" {
			filePathSystemTrx = append(
				filePathSystemTrx,
				path,
			)
		}

		return nil
	})

	log.Err(ctx, "[process.NewSvc] parseSystemTrxFiles afero.Walk SystemTRXPath executed", err)

	if err != nil {
		return
	}

	sliceMutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	parallel.ForEach(filePathSystemTrx, func(item string, _ int) {
		wg.Add(1)
		defer wg.Done()
		data, _ := s.parseSystemTrxFile(ctx, afs, item)
		sliceMutex.Lock()
		returnData = append(returnData, data...)
		sliceMutex.Unlock()
	})

	wg.Wait()

	return
}

func (s *Svc) importReconcileSystemDataToDB(ctx context.Context, data []*systems.SystemTrxData) (err error) {
	defSize := len(data) / s.comp.Config.Data.Reconciliation.NumberWorker
	numBigger := len(data) - defSize*s.comp.Config.Data.Reconciliation.NumberWorker
	size := defSize + 1

	for i, idx := 0, 0; i < s.comp.Config.Data.Reconciliation.NumberWorker; i++ {
		if i == numBigger {
			size--
			if size == 0 {
				break
			}
		}

		err = s.repo.RepoProcess.ImportSystemTrx(
			ctx,
			data[idx:idx+size],
		)

		if err != nil {
			return
		}

		idx += size
	}

	return
}

func (s *Svc) importReconcileMapToDB(ctx context.Context, min float64, max float64) (err error) {
	max = max + 1
	numberWorker := float64(s.comp.Config.Data.Reconciliation.NumberWorker * 2)
	defSize := max / numberWorker
	numBigger := max - defSize*numberWorker
	size := defSize + 1

	for i, idx := 0.0, min; i < numberWorker; i++ {
		if i == numBigger {
			size--
			if size == 0 {
				break
			}
		}

		err = s.repo.RepoProcess.GenerateReconciliationMap(
			ctx,
			idx,
			idx+size,
		)

		if err != nil {
			return
		}

		idx += size
	}

	return
}

func (s *Svc) parseBankTrxFile(ctx context.Context, afs afero.Fs, item FilePathBankTrx) (returnData []*banks.BankTrxData, err error) {
	var bankParser banks.ReconcileBankData
	var f afero.File
	f, err = afs.Open(item.FilePath)

	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()

	bank := strings.ToUpper(item.Bank)

	switch bank {
	case string(banks.BCABankParser):
		{
			bankParser, err = bca.NewBankParser(
				bank,
				csv.NewReader(f),
				true,
			)
		}
	case string(banks.BNIBankParser):
		{
			bankParser, err = bni.NewBankParser(
				bank,
				csv.NewReader(f),
				true,
			)
		}
	default:
		{
			bankParser, err = default_bank.NewBankParser(
				bank,
				csv.NewReader(f),
				true,
			)
		}
	}

	log.Err(ctx, "[process.NewSvc] parseBankTrxFiles parse ("+bank+") - '"+item.Bank+"' executed", err)

	if err != nil {
		return
	}

	returnData, err = bankParser.ToBankTrxData(ctx, item.FilePath)
	log.Err(ctx, "[process.NewSvc] parseBankTrxFiles parse.ToBankTrxData ("+bank+") executed", err)

	return
}

func (s *Svc) parseBankTrxFiles(ctx context.Context, afs afero.Fs) (returnData []*banks.BankTrxData, err error) {
	var filePathBankTrx []FilePathBankTrx
	cleanPath := filepath.Clean(s.comp.Config.Data.Reconciliation.BankTRXPath)
	// scan only csv file with first folder as bank name, bank should in the list of accepted bank name
	er := afero.Walk(afs, cleanPath, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".csv" {
			splitPath := strings.Split(path, cleanPath)
			if len(splitPath) > 1 {
				pathSuffix := strings.Split(splitPath[1], string(os.PathSeparator))
				if len(pathSuffix) > 1 {
					bank := pathSuffix[1]
					if slices.Contains(s.comp.Config.Data.Reconciliation.ListBank, bank) {
						filePathBankTrx = append(
							filePathBankTrx,
							FilePathBankTrx{
								Bank:     bank,
								FilePath: path,
							},
						)
					}
				}
			}
		}

		return nil
	})
	log.Err(ctx, "[process.NewSvc] parseBankTrxFiles afero.Walk BankTRXPath executed", er)
	if er != nil {
		return nil, er
	}

	sliceMutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	parallel.ForEach(filePathBankTrx, func(item FilePathBankTrx, _ int) {
		wg.Add(1)
		defer wg.Done()
		data, _ := s.parseBankTrxFile(ctx, afs, item)
		sliceMutex.Lock()
		returnData = append(returnData, data...)
		sliceMutex.Unlock()
	})

	wg.Wait()

	return
}

func (s *Svc) importReconcileBankDataToDB(ctx context.Context, data []*banks.BankTrxData) (err error) {
	numberWorker := s.comp.Config.Data.Reconciliation.NumberWorker * 2
	defSize := len(data) / numberWorker
	numBigger := len(data) - defSize*numberWorker
	size := defSize + 1

	for i, idx := 0, 0; i < numberWorker; i++ {
		if i == numBigger {
			size--
			if size == 0 {
				break
			}
		}

		err = s.repo.RepoProcess.ImportBankTrx(
			ctx,
			data[idx:idx+size],
		)

		if err != nil {
			return
		}

		idx += size
	}

	return
}

func (s *Svc) parse(ctx context.Context, afs afero.Fs) (trxData parser.TrxData, err error) {
	isOK := func(t time.Time, minDate time.Time, maxDate time.Time) bool {
		return (t.Equal(minDate) || t.After(minDate)) && t.Before(maxDate)
	}

	isOKCheck := func(timeToCheck time.Time) bool {
		return isOK(
			timeToCheck,
			s.comp.Config.Data.Reconciliation.FromDate,
			s.comp.Config.Data.Reconciliation.ToDate.AddDate(0, 0, 1),
		)
	}

	setMinMaxAmount := func(currentAmount float64) {
		if trxData.MinSystemAmount > currentAmount {
			trxData.MinSystemAmount = currentAmount
		}

		if trxData.MaxSystemAmount < currentAmount {
			trxData.MaxSystemAmount = currentAmount
		}
	}

	_, err = hunch.All(
		ctx,
		func(ct context.Context) (d interface{}, e error) {
			defer func() {
				log.Err(ct, "[process.NewSvc] GenerateReconciliation parseSystemTrxFiles executed", e)
			}()

			var data []*systems.SystemTrxData
			data, e = s.parseSystemTrxFiles(ct, afs)

			if e != nil {
				return
			}

			trxData.SystemTrx = lo.Filter(data, func(item *systems.SystemTrxData, index int) bool {
				if !isOKCheck(item.TransactionTime) {
					return false
				}

				setMinMaxAmount(item.Amount)

				return true
			})

			return
		},
		func(ct context.Context) (d interface{}, e error) {
			defer func() {
				log.Err(ct, "[process.NewSvc] GenerateReconciliation parseBankTrxFiles executed", e)
			}()

			var data []*banks.BankTrxData
			data, e = s.parseBankTrxFiles(ct, afs)

			if e != nil {
				return
			}

			trxData.BankTrx = lo.Filter(data, func(item *banks.BankTrxData, index int) bool {
				return isOKCheck(item.Date)
			})

			return
		},
	)

	return
}

func (s *Svc) generateReconciliationSummaryAndFiles(ctx context.Context, fs afero.Fs, isDeleteDirectory bool) (returnData ReconciliationSummary, err error) {
	defer func() {
		log.Err(ctx, "[process.NewSvc] GenerateReconciliation RepoProcess.GetReconciliationSummary executed", err)
	}()

	if summary, er := s.repo.RepoProcess.GetReconciliationSummary(ctx); er != nil {
		return
	} else {
		err = deepcopier.Copy(&summary).To(&returnData)
	}

	if err != nil {
		return
	}

	err = s.generateReconciliationFiles(ctx, &returnData, fs, isDeleteDirectory)
	return
}

func (s *Svc) generateReconciliationFiles(ctx context.Context, reconciliationSummary *ReconciliationSummary, fs afero.Fs, isDeleteDirectory bool) (err error) {
	if reconciliationSummary == nil {
		return
	}

	fileNameSuffix := strconv.FormatInt(time.Now().Unix(), 10)
	logTemplate := "[process.NewSvc] save csv file %s executed"

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (_ interface{}, e error) {
			reconciliationSummary.FileMatchedSystemTrx = fmt.Sprintf("%s/%s/%s/matched_%s.csv", s.comp.Config.Data.Reconciliation.ReportTRXPath, "system", "matched", fileNameSuffix)
			defer func() {
				log.Err(c, fmt.Sprintf(logTemplate, reconciliationSummary.FileMatchedSystemTrx), e)
			}()

			d, er := s.repo.RepoProcess.GetMatchedTrx(ctx)
			if er != nil || d == nil {
				return nil, er
			}

			return nil, csvhelper.StructToCSVFile(
				c,
				fs,
				reconciliationSummary.FileMatchedSystemTrx,
				d,
				isDeleteDirectory,
			)
		},
		func(c context.Context, _ interface{}) (_ interface{}, e error) {
			reconciliationSummary.FileMissingSystemTrx = fmt.Sprintf("%s/%s/%s/notmatched_%s.csv", s.comp.Config.Data.Reconciliation.ReportTRXPath, "system", "notmatched", fileNameSuffix)
			defer func() {
				log.Err(c, fmt.Sprintf(logTemplate, reconciliationSummary.FileMissingSystemTrx), e)
			}()

			d, er := s.repo.RepoProcess.GetNotMatchedSystemTrx(ctx)
			if er != nil || d == nil {
				return nil, er
			}

			return nil, csvhelper.StructToCSVFile(
				c,
				fs,
				reconciliationSummary.FileMissingSystemTrx,
				d,
				isDeleteDirectory,
			)
		},
		func(c context.Context, _ interface{}) (_ interface{}, e error) {
			d, er := s.repo.RepoProcess.GetNotMatchedBankTrx(ctx)
			if er != nil || d == nil {
				return nil, er
			}

			bankTrxData := make(map[string][]process.NotMatchedBankTrx)
			lo.ForEach(d, func(data process.NotMatchedBankTrx, _ int) {
				data.Bank = strings.ToLower(data.Bank)
				bankTrxData[data.Bank] = append(bankTrxData[data.Bank], data)
			})

			reconciliationSummary.FileMissingBankTrx = make(map[string]string)
			bankNames := lo.Keys(bankTrxData)
			if isDeleteDirectory {
				dirReportBankTrx := fmt.Sprintf("%s/%s/%s", s.comp.Config.Data.Reconciliation.ReportTRXPath, "bank", "notmatched")
				e = csvhelper.DeleteDirectory(c, fs, dirReportBankTrx)
			}

			if e != nil {
				return
			}

			parallel.ForEach(bankNames, func(item string, _ int) {
				fileReportBankTrx := fmt.Sprintf("%s/%s/%s/%s_%s.csv", s.comp.Config.Data.Reconciliation.ReportTRXPath, "bank", "notmatched", item, fileNameSuffix)
				e := csvhelper.StructToCSVFile(
					c,
					fs,
					fileReportBankTrx,
					bankTrxData[item],
					false,
				)
				reconciliationSummary.FileMissingBankTrx[item] = fileReportBankTrx
				log.Err(c, fmt.Sprintf(logTemplate, fileReportBankTrx), e)
			})

			return nil, e
		},
	)

	return
}

func (s *Svc) GenerateReconciliation(ctx context.Context, afs afero.Fs, bar *progressbar.ProgressBar) (returnData ReconciliationSummary, err error) {
	ctx = s.comp.Logger.GetLogger().With().Str("component", "Process Service").Ctx(ctx).Logger().WithContext(s.comp.Logger.GetCtx())
	defer func() {
		_ = s.repo.RepoProcess.Close()
		progressbarhelper.BarClear(bar)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			progressbarhelper.BarDescribe(bar, "[cyan][1/7] Pre Process Generate Reconciliation...")

			er := s.repo.RepoProcess.Pre(
				c,
				s.comp.Config.Data.Reconciliation.ListBank,
				s.comp.Config.Data.Reconciliation.FromDate,
				s.comp.Config.Data.Reconciliation.ToDate,
			)

			log.Err(c, "[process.NewSvc] GenerateReconciliation RepoProcess.Pre executed", er)
			return nil, err
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			progressbarhelper.BarDescribe(bar, "[cyan][2/7] Parse System/Bank Trx Files...")
			return s.parse(c, afs)
		},
		func(c context.Context, i interface{}) (d interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][3/7] Import System Trx to DB...")

			if i == nil {
				e = errors.New("empty parse data")
				return
			}

			data := i.(parser.TrxData)

			e = s.importReconcileSystemDataToDB(c, data.SystemTrx)
			log.Err(c, "[process.NewSvc] GenerateReconciliation importReconcileSystemDataToDB executed", e)

			progressbarhelper.BarDescribe(bar, "[cyan][4/7] Import Bank Trx to DB...")

			e = s.importReconcileBankDataToDB(c, data.BankTrx)
			log.Err(c, "[process.NewSvc] GenerateReconciliation importReconcileBankDataToDB executed", e)

			progressbarhelper.BarDescribe(bar, "[cyan][5/7] Mapping Reconciliation Data...")

			e = s.importReconcileMapToDB(c, data.MinSystemAmount, data.MaxSystemAmount)
			log.Err(c, "[process.NewSvc] GenerateReconciliation importReconcileMapToDB executed", e)

			return
		},
		func(c context.Context, i interface{}) (d interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][6/7] Generate Reconciliation Report Files...")
			defer func() {
				log.Err(c, "[process.NewSvc] GenerateReconciliation generateReconciliationSummaryAndFiles executed", e)
			}()

			returnData, e = s.generateReconciliationSummaryAndFiles(c, afs, true)
			return
		},
		func(c context.Context, i interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][7/8] Post Process Generate Reconciliation...")
			if !s.comp.Config.Data.IsDebug {
				e = s.repo.RepoProcess.Post(
					c,
				)
				log.Err(c, "[process.NewSvc] GenerateReconciliation RepoProcess.Post executed", e)
			}

			return nil, e
		},
	)

	return
}
