package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"

	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// DBPostgresOption options for postgres connection
type DBPostgresOption struct {
	LogPrefix         string        `deepcopier:"skip"`
	Host              string        `deepcopier:"field:Host"`
	DB                string        `deepcopier:"field:DB"`
	Username          string        `deepcopier:"field:Username"`
	Password          string        `deepcopier:"field:Password"`
	SSLMode           string        `deepcopier:"field:sslmode"`
	Port              int           `deepcopier:"field:Port"`
	MaxPoolSize       int           `deepcopier:"field:MaxPoolSize"`
	ConnMaxLifetime   time.Duration `deepcopier:"field:ConnMaxLifetime"`
	MaxIdleConnection int           `deepcopier:"field:MaxIdleConnection"`
}

// NewPostgresDatabase return gorp dbmap object with postgres options param
func NewPostgresDatabase(option DBPostgresOption, logger zerolog.Logger, isDoLogging bool) (db *sql.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		option.Host,
		option.Port,
		option.Username,
		option.DB,
		option.Password,
		option.SSLMode,
	)

	loggerAdapter := NewNoopLog()

	if isDoLogging {
		loggerAdapter = zerologadapter.New(logger)
	}

	driverName := "postgres"
	dbOtel, err := otelsql.Open(
		driverName,
		dsn,
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
		),
	)

	if err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(dbOtel, driverName)
	db = sqldblogger.OpenDriver(dsn, dbx.Driver(), loggerAdapter)
	db.SetMaxOpenConns(option.MaxPoolSize)
	db.SetConnMaxLifetime(option.ConnMaxLifetime)
	db.SetMaxIdleConns(option.MaxIdleConnection)

	// Register DB stats to meter
	err = otelsql.RegisterDBStatsMetrics(
		db,
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
		),
	)

	return db, err
}
