package core

import (
	"simple-reconciliation-service/internal/pkg/driver/sql"
	"time"

	"github.com/ulule/deepcopier"
)

// PostgresParameters ..
type PostgresParameters struct {
	Password          string        `mapstructure:"password"`
	Host              string        `mapstructure:"host"`
	DB                string        `mapstructure:"db"`
	Username          string        `mapstructure:"username"`
	SSLMode           string        `mapstructure:"sslmode"`
	MaxPoolSize       int           `mapstructure:"max_pool_size"`
	Port              int           `mapstructure:"port"`
	ConnMaxLifetime   time.Duration `mapstructure:"conn_max_lifetime"`
	MaxIdleConnection int           `mapstructure:"max_idle_connection"`
	IsEnabled         bool          `default:"false"                    mapstructure:"is_enabled"`
	IsMigrationEnable bool          `mapstructure:"is_migration_enable"`
}

func (pp *PostgresParameters) Options(logPrefix string) (returnData sql.DBPostgresOption) {
	_ = deepcopier.Copy(pp).To(&returnData)
	returnData.LogPrefix = logPrefix
	return
}

// Postgres ..
type Postgres struct {
	Write     PostgresParameters `mapstructure:"write"`
	Read      PostgresParameters `mapstructure:"read"`
	IsEnabled bool               `default:"true"       mapstructure:"is_enabled"`
}
