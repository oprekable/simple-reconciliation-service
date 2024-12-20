package _helper

import (
	"context"
	"database/sql"
	"github.com/aaronjan/hunch"
)

type StmtData struct {
	Query string
	Args  []any
}

func ExecTxQueries(ctx context.Context, db *sql.DB, tx *sql.Tx, stmtData []StmtData) (err error) {
	var executableInSequence []hunch.ExecutableInSequence
	for k := range stmtData {
		executableInSequence = append(
			executableInSequence,
			func(c context.Context, _ interface{}) (interface{}, error) {
				i, e := db.PrepareContext(
					c,
					stmtData[k].Query,
				)

				if e != nil {
					return nil, e
				}

				return tx.StmtContext(c, i).ExecContext( //nolint:sqlclosecheck
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
