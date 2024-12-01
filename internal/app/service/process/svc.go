package process

import (
	"context"
	"encoding/csv"
	"io/fs"
	"os"
	"path/filepath"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/pkg/reconcile/parser"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"slices"
	"strings"

	"github.com/aaronjan/hunch"

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

func (s *Svc) GenerateReconciliation(ctx context.Context, afs afero.Fs) (returnSummary ReconciliationSummary, err error) {
	ctx = s.comp.Logger.GetLogger().With().Str("component", "Process Service").Ctx(ctx).Logger().WithContext(s.comp.Logger.GetCtx())
	//var reconciliationData []process.ReconciliationData

	defer func() {
		_ = s.repo.RepoProcess.Close()
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			er := s.repo.RepoProcess.Pre(
				c,
				s.comp.Config.Reconciliation.ListBank,
				s.comp.Config.Reconciliation.FromDate,
				s.comp.Config.Reconciliation.ToDate,
			)

			log.AddErr(c, er)
			log.Msg(c, "[process.NewSvc] RepoProcess.Pre executed")

			return nil, err
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
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

			if er != nil {
				return nil, er
			}

			for k := range filePathSystemTrx {
				var systemParser parser.ReconcileSystemData
				f, er := afs.Open(filePathSystemTrx[k])
				if er != nil {
					return nil, err
				}

				systemParser, err = parser.NewDefaultSystem(
					csv.NewReader(f),
				)

				er = s.repo.RepoProcess.ImportSystemTrx(
					c,
					systemParser,
				)

				if er != nil {
					return nil, er
				}
			}

			return nil, nil
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			filePathBankTrx := make(map[string][]string)
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
								filePathBankTrx[bank] = append(
									filePathBankTrx[bank],
									path,
								)
							}
						}
					}
				}

				return nil
			})

			if er != nil {
				return nil, er
			}

			for bank, fileBankPath := range filePathBankTrx {
				bank = strings.ToUpper(bank)
				for k := range fileBankPath {
					var bankParser parser.ReconcileBankData
					f, er := afs.Open(fileBankPath[k])
					if er != nil {
						return nil, err
					}

					switch bank {
					case string(parser.BCA):
						{
							bankParser, err = parser.NewBCABank(
								bank,
								csv.NewReader(f),
							)
						}
					case string(parser.BNI):
						{
							bankParser, err = parser.NewBNIBank(
								bank,
								csv.NewReader(f),
							)
						}
					default:
						{
							bankParser, err = parser.NewDefaultBank(
								bank,
								csv.NewReader(f),
							)
						}
					}

					er = s.repo.RepoProcess.ImportBankTrx(
						c,
						bank,
						bankParser,
					)

					if er != nil {
						return nil, er
					}
				}
			}

			return nil, er
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			er := s.repo.RepoProcess.Post(
				c,
			)

			log.AddErr(c, er)
			log.Msg(c, "[process.NewSvc] RepoProcess.Post executed")

			return nil, er
		},
	)

	return
}
