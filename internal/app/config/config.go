package config

import (
	"simple-reconciliation-service/internal/app/config/core"
	"simple-reconciliation-service/internal/app/config/reconciliation"
)

type Data struct {
	core.App
	reconciliation.Reconciliation
	core.Postgres
}
