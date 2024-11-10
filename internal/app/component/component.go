package component

import (
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/cerror"
	"simple-reconciliation-service/internal/app/component/clogger"
)

type Components struct {
	Config *cconfig.Config
	Logger *clogger.Logger
	Error  *cerror.Error
}
