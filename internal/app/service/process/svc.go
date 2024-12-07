package process

import (
	"context"
	"encoding/csv"
	"io/fs"
	"os"
	"path/filepath"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/repository/process"
	"simple-reconciliation-service/internal/pkg/reconcile/parser"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"slices"
	"strings"

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

func (s *Svc) parseSystemTrxFiles(ctx context.Context, afs afero.Fs) (returnData []*parser.SystemTrxData, err error) {
	var filePathSystemTrx []string
	cleanPath := filepath.Clean(s.comp.Config.Reconciliation.SystemTRXPath)
	er := afero.Walk(afs, cleanPath, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".csv" {
			filePathSystemTrx = append(
				filePathSystemTrx,
				path,
			)
		}

		return nil
	})
	log.Err(ctx, "[process.NewSvc] parseSystemTrxFiles afero.Walk SystemTRXPath executed", er)
	if er != nil {
		return nil, er
	}
	parallel.ForEach(filePathSystemTrx, func(item string, _ int) {
		f, er := afs.Open(item)
		log.Err(ctx, "[process.NewSvc] parseSystemTrxFiles fs.Open - '"+item+"'", er)
		if er != nil {
			if f != nil {
				_ = f.Close()
			}
			return
		}
		systemParser, er := parser.NewDefaultSystem(
			csv.NewReader(f),
			true,
		)
		log.Err(ctx, "[process.NewSvc] parseSystemTrxFiles - '"+item+"' parse", er)
		if er != nil {
			if f != nil {
				_ = f.Close()
			}
			return
		}
		data, er := systemParser.ToSystemTrxData(ctx, item)
		log.Err(ctx, "[process.NewSvc] parseSystemTrxFiles parse.ToSystemTrxData", er)
		if er != nil {
			if f != nil {
				_ = f.Close()
			}
			return
		}
		returnData = append(returnData, data...)
		if f != nil {
			_ = f.Close()
		}
	})

	return
}

func (s *Svc) importReconcileSystemDataToDB(ctx context.Context, data []*parser.SystemTrxData) (err error) {
	defSize := len(data) / s.comp.Config.Reconciliation.NumberWorker
	numBigger := len(data) - defSize*s.comp.Config.Reconciliation.NumberWorker
	size := defSize + 1

	for i, idx := 0, 0; i < s.comp.Config.Reconciliation.NumberWorker; i++ {
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

func (s *Svc) parseBankTrxFiles(ctx context.Context, afs afero.Fs) (returnData []*parser.BankTrxData, err error) {
	type FilePathBankTrx struct {
		Bank     string
		FilePath string
	}
	var filePathBankTrx []FilePathBankTrx
	cleanPath := filepath.Clean(s.comp.Config.Reconciliation.BankTRXPath)
	// scan only csv file with first folder as bank name, bank should in the list of accepted bank name
	er := afero.Walk(afs, cleanPath, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".csv" {
			splitPath := strings.Split(path, cleanPath)
			if len(splitPath) > 1 {
				pathSuffix := strings.Split(splitPath[1], string(os.PathSeparator))
				if len(pathSuffix) > 1 {
					bank := pathSuffix[1]
					if slices.Contains(s.comp.Config.Reconciliation.ListBank, bank) {
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
		bank := strings.ToUpper(item.Bank)
		var bankParser parser.ReconcileBankData
		f, er := afs.Open(item.FilePath)
		log.Err(ctx, "[process.NewSvc] parseBankTrxFiles fs.Open - '"+item.FilePath+"'", er)
		if er != nil {
			if f != nil {
				_ = f.Close()
			}
			return
		}

		switch bank {
		case string(parser.BCABankParser):
			{
				bankParser, er = parser.NewBCABank(
					bank,
					csv.NewReader(f),
					true,
				)
			}
		case string(parser.BNIBankParser):
			{
				bankParser, er = parser.NewBNIBank(
					bank,
					csv.NewReader(f),
					true,
				)
			}
		default:
			{
				bankParser, er = parser.NewDefaultBank(
					bank,
					csv.NewReader(f),
					true,
				)
			}
		}
		log.Err(ctx, "[process.NewSvc] parseBankTrxFiles parse ("+bank+") - '"+item.Bank+"' executed", er)
		if er != nil {
			if f != nil {
				_ = f.Close()
			}
			return
		}
		data, er := bankParser.ToBankTrxData(ctx, item.FilePath)
		log.Err(ctx, "[process.NewSvc] parseBankTrxFiles parse.ToBankTrxData ("+bank+") executed", er)
		if er != nil {
			if f != nil {
				_ = f.Close()
			}
			return
		}
		returnData = append(returnData, data...)
		if f != nil {
			_ = f.Close()
		}
	})

	return
}

func (s *Svc) importReconcileBankDataToDB(ctx context.Context, data []*parser.BankTrxData) (err error) {
	defSize := len(data) / s.comp.Config.Reconciliation.NumberWorker
	numBigger := len(data) - defSize*s.comp.Config.Reconciliation.NumberWorker
	size := defSize + 1

	for i, idx := 0, 0; i < s.comp.Config.Reconciliation.NumberWorker; i++ {
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

func (s *Svc) generateReconciliationSummaryAndFiles(ctx context.Context, data []process.ReconciliationData) (returnData ReconciliationSummary, err error) {
	return
}

func (s *Svc) GenerateReconciliation(ctx context.Context, afs afero.Fs, bar *progressbar.ProgressBar) (returnData ReconciliationSummary, err error) {
	ctx = s.comp.Logger.GetLogger().With().Str("component", "Process Service").Ctx(ctx).Logger().WithContext(s.comp.Logger.GetCtx())
	defer func() {
		_ = s.repo.RepoProcess.Close()
		if bar != nil {
			_ = bar.Clear()
		}
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][1/8] Pre Process Generate Reconciliation...")
			}

			er := s.repo.RepoProcess.Pre(
				c,
				s.comp.Config.Reconciliation.ListBank,
				s.comp.Config.Reconciliation.FromDate,
				s.comp.Config.Reconciliation.ToDate,
			)

			log.Err(c, "[process.NewSvc] GenerateReconciliation RepoProcess.Pre executed", er)
			return nil, err
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][2/8] Parse System Trx Files...")
			}

			data, er := s.parseSystemTrxFiles(c, afs)
			log.Err(c, "[process.NewSvc] GenerateReconciliation parseSystemTrxFiles executed", er)
			if er != nil {
				return nil, er
			}

			if bar != nil {
				bar.Describe("[cyan][3/8] Import System Trx to DB...")
			}

			er = s.importReconcileSystemDataToDB(c, data)
			log.Err(c, "[process.NewSvc] GenerateReconciliation importReconcileSystemDataToDB executed", er)
			return nil, er
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][4/8] Parse Bank Trx Files...")
			}

			data, er := s.parseBankTrxFiles(c, afs)
			log.Err(c, "[process.NewSvc] GenerateReconciliation parseBankTrxFiles executed", er)
			if er != nil {
				return nil, er
			}

			if bar != nil {
				bar.Describe("[cyan][5/8] Import Bank Trx to DB...")
			}

			er = s.importReconcileBankDataToDB(c, data)
			log.Err(c, "[process.NewSvc] GenerateReconciliation importReconcileBankDataToDB executed", er)

			return nil, er
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][6/8] Calculate Reconciliation Data...")
			}

			data, er := s.repo.RepoProcess.GetReconciliation(
				c,
			)
			log.Err(c, "[process.NewSvc] GenerateReconciliation RepoProcess.GetReconciliation executed", er)

			if bar != nil {
				bar.Describe("[cyan][7/8] Generate Reconciliation Report Files...")
			}

			returnData, er = s.generateReconciliationSummaryAndFiles(c, data)
			if er != nil {
				return nil, er
			}
			log.Err(c, "[process.NewSvc] GenerateReconciliation generateReconciliationSummaryAndFiles executed", er)
			return nil, er
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if bar != nil {
				bar.Describe("[cyan][8/8] Post Process Generate Reconciliation...")
			}

			er := s.repo.RepoProcess.Post(
				c,
			)
			log.Err(c, "[process.NewSvc] GenerateReconciliation RepoProcess.Post executed", er)
			return nil, er
		},
	)

	return
}
