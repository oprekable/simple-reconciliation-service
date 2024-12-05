package process

import (
	"context"
	"database/sql"
	"simple-reconciliation-service/internal/pkg/reconcile/parser"
	"time"
)

type DB struct {
	db *sql.DB
}

var _ Repository = (*DB)(nil)

func NewDB(
	db *sql.DB,
) (*DB, error) {
	return &DB{
		db: db,
	}, nil
}

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time) (err error) {
	return
}

func (d *DB) ImportSystemTrx(ctx context.Context, systemParser parser.ReconcileSystemData) (err error) {
	return
}

func (d *DB) ImportBankTrx(ctx context.Context, bank string, bankParser parser.ReconcileBankData, numWorkers int) (err error) {
	// TODO: import data to DB
	return
}

func (d *DB) GetReconciliation(ctx context.Context) (returnData []ReconciliationData, err error) {
	return
}

func (d *DB) Post(ctx context.Context) (err error) {
	return
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}
