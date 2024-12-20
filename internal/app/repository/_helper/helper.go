package _helper

import (
	"context"
	"database/sql"
	"github.com/aaronjan/hunch"
	"github.com/pkg/errors"
)

type StmtData struct {
	Name  string
	Query string
	Args  []any
}

func ExecTxQueries(ctx context.Context, db *sql.DB, tx *sql.Tx, stmtMap map[string]*sql.Stmt, stmtData []StmtData) (err error) {
	var executableInSequence []hunch.ExecutableInSequence
	for k := range stmtData {
		executableInSequence = append(
			executableInSequence,
			func(c context.Context, _ interface{}) (r interface{}, e error) {
				if _, ok := stmtMap[stmtData[k].Name]; !ok {
					stmtMap[stmtData[k].Name], e = db.PrepareContext(
						c,
						stmtData[k].Query,
					)

					if e != nil {
						return nil, e
					}
				}

				return tx.StmtContext(c, stmtMap[stmtData[k].Name]).ExecContext( //nolint:sqlclosecheck
					c,
					stmtData[k].Args...,
				)
			},
		)
	}

	_, err = hunch.Waterfall(
		ctx,
		executableInSequence...,
	)

	return
}

func CommitOrRollback(tx *sql.Tx, er error) (err error) {
	if er != nil {
		err = errors.Wrap(tx.Rollback(), er.Error())
	} else {
		err = tx.Commit()
	}

	return
}
