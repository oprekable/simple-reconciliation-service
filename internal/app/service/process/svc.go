package process

import (
	"context"
	"encoding/csv"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/pkg/reconcile/parser"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"simple-reconciliation-service/internal/pkg/utils/progressbarhelper"
	"slices"
	"strings"
	"time"

	"github.com/ulule/deepcopier"

	"github.com/samber/lo"

	"github.com/aaronjan/hunch"
	"github.com/samber/lo/parallel"
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

func (s *Svc) parseSystemTrxFile(ctx context.Context, afs afero.Fs, filePath string) (returnData []*parser.SystemTrxData, err error) {
	var f afero.File
	f, err = afs.Open(filePath)

	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()

	var systemParser *parser.DefaultSystem
	systemParser, err = parser.NewDefaultSystem(
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

func (s *Svc) parseSystemTrxFiles(ctx context.Context, afs afero.Fs) (returnData []*parser.SystemTrxData, err error) {
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

	parallel.ForEach(filePathSystemTrx, func(item string, _ int) {
		data, _ := s.parseSystemTrxFile(ctx, afs, item)
		returnData = append(returnData, data...)
	})

	return
}

func (s *Svc) importReconcileSystemDataToDB(ctx context.Context, data []*parser.SystemTrxData) (err error) {
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

func (s *Svc) parseBankTrxFile(ctx context.Context, afs afero.Fs, item FilePathBankTrx) (returnData []*parser.BankTrxData, err error) {
	var bankParser parser.ReconcileBankData
	var f afero.File
	f, err = afs.Open(item.FilePath)

	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()

	bank := strings.ToUpper(item.Bank)

	switch bank {
	case string(parser.BCABankParser):
		{
			bankParser, err = parser.NewBCABank(
				bank,
				csv.NewReader(f),
				true,
			)
		}
	case string(parser.BNIBankParser):
		{
			bankParser, err = parser.NewBNIBank(
				bank,
				csv.NewReader(f),
				true,
			)
		}
	default:
		{
			bankParser, err = parser.NewDefaultBank(
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

func (s *Svc) parseBankTrxFiles(ctx context.Context, afs afero.Fs) (returnData []*parser.BankTrxData, err error) {
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

	parallel.ForEach(filePathBankTrx, func(item FilePathBankTrx, _ int) {
		data, _ := s.parseBankTrxFile(ctx, afs, item)
		returnData = append(returnData, data...)
	})

	return
}

func (s *Svc) importReconcileBankDataToDB(ctx context.Context, data []*parser.BankTrxData) (err error) {
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
	isOKSystemTrx := func(timeToCheck string) bool {
		t, e := time.Parse("2006-01-02 15:04:05", timeToCheck)
		if e != nil {
			return false
		}

		minDate := s.comp.Config.Data.Reconciliation.FromDate
		maxDate := s.comp.Config.Data.Reconciliation.ToDate.AddDate(0, 0, 1)

		return (t.Equal(minDate) || t.After(minDate)) && t.Before(maxDate)
	}

	isOKBankTrx := func(timeToCheck string) bool {
		t, e := time.Parse("2006-01-02", timeToCheck)
		if e != nil {
			return false
		}

		minDate := s.comp.Config.Data.Reconciliation.FromDate
		maxDate := s.comp.Config.Data.Reconciliation.ToDate.AddDate(0, 0, 1)

		return (t.Equal(minDate) || t.After(minDate)) && t.Before(maxDate)
	}

	_, err = hunch.All(
		ctx,
		func(ct context.Context) (d interface{}, e error) {
			defer func() {
				log.Err(ct, "[process.NewSvc] GenerateReconciliation parseSystemTrxFiles executed", e)
			}()

			var data []*parser.SystemTrxData
			data, e = s.parseSystemTrxFiles(ct, afs)

			if e != nil {
				return
			}

			trxData.SystemTrx = lo.Filter(data, func(item *parser.SystemTrxData, index int) bool {
				if !isOKSystemTrx(item.TransactionTime) {
					return false
				}

				if trxData.MinSystemAmount > item.Amount {
					trxData.MinSystemAmount = item.Amount
				}

				if trxData.MaxSystemAmount < item.Amount {
					trxData.MaxSystemAmount = item.Amount
				}

				return true
			})

			return
		},
		func(ct context.Context) (d interface{}, e error) {
			defer func() {
				log.Err(ct, "[process.NewSvc] GenerateReconciliation parseBankTrxFiles executed", e)
			}()

			var data []*parser.BankTrxData
			data, e = s.parseBankTrxFiles(ct, afs)

			if e != nil {
				return
			}

			trxData.BankTrx = lo.Filter(data, func(item *parser.BankTrxData, index int) bool {
				return isOKBankTrx(item.Date)
			})

			return
		},
	)

	return
}

func (s *Svc) generateReconciliationSummaryAndFiles(ctx context.Context) (returnData ReconciliationSummary, err error) {
	defer func() {
		log.Err(ctx, "[process.NewSvc] GenerateReconciliation RepoProcess.GetReconciliationSummary executed", err)
	}()

	summary, er := s.repo.RepoProcess.GetReconciliationSummary(ctx)
	if er != nil {
		err = er
		return
	}

	err = deepcopier.Copy(&summary).To(&returnData)
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

			returnData, e = s.generateReconciliationSummaryAndFiles(c)
			return
		},
		func(c context.Context, i interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][7/7] Post Process Generate Reconciliation...")
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
