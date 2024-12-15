package sample

import (
	"context"
	"database/sql"
	"encoding/json"
	"simple-reconciliation-service/internal/app/err/core"
	"simple-reconciliation-service/internal/pkg/utils/log"
	"strconv"
	"strings"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/blockloop/scan/v2"
	"github.com/pkg/errors"
)

type DB struct {
	db                           *sql.DB
	stmtDropTableArguments       *sql.Stmt
	stmtDropTableBanks           *sql.Stmt
	stmtDropTableBaseData        *sql.Stmt
	stmtCreateTableArguments     *sql.Stmt
	stmtCreateTableBanks         *sql.Stmt
	stmtCreateTableBaseData      *sql.Stmt
	stmtCreateIndexTableBaseData *sql.Stmt
	stmtGetTrxData               *sql.Stmt
}

var _ Repository = (*DB)(nil)

func NewDB(
	db *sql.DB,
) (*DB, error) {
	return &DB{
		db: db,
	}, nil
}

func (d *DB) dropTables(ctx context.Context, tx *sql.Tx) (err error) {
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtDropTableArguments == nil {
				d.stmtDropTableArguments, e = d.db.PrepareContext(c, QueryDropTableArguments)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtDropTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtDropTableBanks == nil {
				d.stmtDropTableBanks, e = d.db.PrepareContext(c, QueryDropTableBanks)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtDropTableBanks).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtDropTableBaseData == nil {
				d.stmtDropTableBaseData, e = d.db.PrepareContext(c, QueryDropTableBaseData)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtDropTableBaseData).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
	)

	return
}

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time, limitTrxData int64, matchPercentage int) (err error) {
	defer func() {
		log.Err(
			ctx,
			"[sample.NewDB] Exec Pre method in db",
			err,
		)
	}()

	var tx *sql.Tx
	tx, err = d.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, d.dropTables(c, tx)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateTableArguments == nil {
				d.stmtCreateTableArguments, e = d.db.PrepareContext(c, QueryCreateTableArguments)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtCreateTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
				startDate.Format("2006-01-02"),
				toDate.Format("2006-01-02"),
				strconv.FormatInt(limitTrxData, 10),
				strconv.Itoa(matchPercentage),
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateTableBanks == nil {
				d.stmtCreateTableBanks, e = d.db.PrepareContext(c, QueryCreateTableBanks)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			b := new(strings.Builder)
			err := json.NewEncoder(b).Encode(listBank)
			if err != nil {
				return nil, err
			}

			return tx.StmtContext(c, d.stmtCreateTableBanks).ExecContext( //nolint:sqlclosecheck
				c,
				b.String(),
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateTableBaseData == nil {
				d.stmtCreateTableBaseData, e = d.db.PrepareContext(c, QueryCreateTableBaseData)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			b := new(strings.Builder)
			err := json.NewEncoder(b).Encode(listBank)
			if err != nil {
				return nil, err
			}

			return tx.StmtContext(c, d.stmtCreateTableBaseData).ExecContext( //nolint:sqlclosecheck
				c,
				b.String(),
			)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtCreateIndexTableBaseData == nil {
				d.stmtCreateIndexTableBaseData, e = d.db.PrepareContext(c, QueryCreateIndexTableBaseData)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return tx.StmtContext(c, d.stmtCreateIndexTableBaseData).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

	return
}

func (d *DB) GetTrx(ctx context.Context) (returnData []TrxData, err error) {
	defer func() {
		log.Err(
			ctx,
			"[sample.NewDB] Exec GetData method in db",
			err,
		)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if d.stmtGetTrxData == nil {
				d.stmtGetTrxData, e = d.db.PrepareContext(c, QueryGetTrxData)
			}

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return d.stmtGetTrxData.QueryContext(
				c,
			)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			rows := i.(*sql.Rows)
			return nil, scan.RowsStrict(&returnData, rows)
		},
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = core.CErrDBConn.Error()
	}

	return
}

func (d *DB) Post(ctx context.Context) (err error) {
	defer func() {
		log.Err(
			ctx,
			"[sample.NewDB] Exec Post method in db",
			err,
		)
	}()

	var tx *sql.Tx
	tx, err = d.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, d.dropTables(c, tx)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

	return
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}
