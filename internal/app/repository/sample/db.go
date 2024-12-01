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

	"github.com/blockloop/scan/v2"
	"github.com/pkg/errors"

	"github.com/aaronjan/hunch"
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

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time, limitTrxData int64, matchPercentage int) (err error) {
	var tx *sql.Tx

	e := d.db.Ping()
	if e != nil {
		log.AddErr(ctx, e)
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return d.db.BeginTx(ctx, nil)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				tx = i.(*sql.Tx)
			}

			return nil, nil
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableArguments == nil {
				return d.db.PrepareContext(c, QueryDropTableArguments)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableArguments = i.(*sql.Stmt)
			}

			return tx.StmtContext(ctx, d.stmtDropTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableBanks == nil {
				return d.db.PrepareContext(c, QueryDropTableBanks)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableBanks = i.(*sql.Stmt)
			}

			return tx.StmtContext(ctx, d.stmtDropTableBanks).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableBaseData == nil {
				return d.db.PrepareContext(c, QueryDropTableBaseData)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableBaseData = i.(*sql.Stmt)
			}

			return tx.StmtContext(ctx, d.stmtDropTableBaseData).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateTableArguments == nil {
				return d.db.PrepareContext(c, QueryCreateTableArguments)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateTableArguments = i.(*sql.Stmt)
			}

			return tx.StmtContext(ctx, d.stmtCreateTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
				startDate.Format("2006-01-02"),
				toDate.Format("2006-01-02"),
				strconv.FormatInt(limitTrxData, 10),
				strconv.Itoa(matchPercentage),
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateTableBanks == nil {
				return d.db.PrepareContext(c, QueryCreateTableBanks)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateTableBanks = i.(*sql.Stmt)
			}

			b := new(strings.Builder)
			err := json.NewEncoder(b).Encode(listBank)
			if err != nil {
				return nil, err
			}

			return tx.StmtContext(ctx, d.stmtCreateTableBanks).ExecContext( //nolint:sqlclosecheck
				c,
				b.String(),
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateTableBaseData == nil {
				return d.db.PrepareContext(c, QueryCreateTableBaseData)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateTableBaseData = i.(*sql.Stmt)
			}

			b := new(strings.Builder)
			err := json.NewEncoder(b).Encode(listBank)
			if err != nil {
				return nil, err
			}

			return tx.StmtContext(ctx, d.stmtCreateTableBaseData).ExecContext( //nolint:sqlclosecheck
				c,
				b.String(),
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtCreateIndexTableBaseData == nil {
				return d.db.PrepareContext(c, QueryCreateIndexTableBaseData)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtCreateIndexTableBaseData = i.(*sql.Stmt)
			}

			return tx.StmtContext(ctx, d.stmtCreateIndexTableBaseData).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

	log.AddErr(ctx, err)
	log.Msg(
		ctx,
		"[sample.NewDB] Exec Pre method from db",
	)

	return
}

func (d *DB) GetTrx(ctx context.Context) (returnData []TrxData, err error) {
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtGetTrxData == nil {
				return d.db.PrepareContext(c, QueryGetTrxData)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtGetTrxData = i.(*sql.Stmt)
			}

			return d.stmtGetTrxData.QueryContext(
				c,
			)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			rows := i.(*sql.Rows)
			return nil, scan.RowsStrict(&returnData, rows)
		},
	)

	log.AddErr(ctx, err)
	log.Msg(
		ctx,
		"[sample.NewDB] Exec GetData method from db",
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = core.CErrDBConn.Error()
	}

	return
}

func (d *DB) Post(ctx context.Context) (err error) {
	var tx *sql.Tx

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return d.db.BeginTx(ctx, nil)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				tx = i.(*sql.Tx)
			}

			return nil, nil
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableArguments == nil {
				return d.db.PrepareContext(c, QueryDropTableArguments)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableArguments = i.(*sql.Stmt)
			}

			return tx.StmtContext(ctx, d.stmtDropTableArguments).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableBanks == nil {
				return d.db.PrepareContext(c, QueryDropTableBanks)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableBanks = i.(*sql.Stmt)
			}

			return tx.StmtContext(ctx, d.stmtDropTableBanks).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			if d.stmtDropTableBaseData == nil {
				return d.db.PrepareContext(c, QueryDropTableBaseData)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				d.stmtDropTableBaseData = i.(*sql.Stmt)
			}

			return tx.StmtContext(ctx, d.stmtDropTableBaseData).ExecContext( //nolint:sqlclosecheck
				c,
			)
		},
	)

	if err != nil {
		err = errors.Wrap(tx.Rollback(), err.Error())
	} else {
		err = tx.Commit()
	}

	log.AddErr(ctx, err)
	log.Msg(
		ctx,
		"[sample.NewDB] Exec Post method from db",
	)

	return
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}
