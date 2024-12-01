package service

import (
	"simple-reconciliation-service/internal/app/service/process"
	"simple-reconciliation-service/internal/app/service/sample"

	"github.com/google/wire"
)

func NewServices(
	svcSample sample.Service,
	svcProcess process.Service,
) *Services {
	return &Services{
		SvcSample:  svcSample,
		SvcProcess: svcProcess,
	}
}

var Set = wire.NewSet(
	sample.Set,
	process.Set,
	NewServices,
)
