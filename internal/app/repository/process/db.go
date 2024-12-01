package process

import (
	"context"
	"database/sql"
	"fmt"
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
	sqlInsertPattern := "INSERT INTO system_trx VALUES ('%s', '%s', '%s', %f);\n"
	s, er := systemParser.ToSql(
		ctx,
		true,
		sqlInsertPattern,
	)
	if er != nil {
		return er
	}

	fmt.Println(s)

	// TODO: import data to DB

	return
}

func (d *DB) ImportBankTrx(ctx context.Context, bank string, bankParser parser.ReconcileBankData) (err error) {
	sqlInsertPattern := "INSERT INTO bank_trx VALUES ('%s', '%s', '%s', '%s', %f);\n"
	s, er := bankParser.ToSql(
		ctx,
		true,
		sqlInsertPattern,
	)
	if er != nil {
		return er
	}
	fmt.Println(s)

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
