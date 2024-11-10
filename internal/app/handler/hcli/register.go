package hcli

import (
	"simple-reconciliation-service/internal/app/handler/hcli/noop"
	"simple-reconciliation-service/internal/app/handler/hcli/process"
	"simple-reconciliation-service/internal/app/handler/hcli/sample"
)

var (
	Handlers = append(commonHandlers, applicationHandlers...)

	commonHandlers = []Handler{
		noop.NewHandler(),
	}

	applicationHandlers = []Handler{
		process.NewHandler(),
		sample.NewHandler(),
	}
)
