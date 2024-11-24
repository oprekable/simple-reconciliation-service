package service

import (
	"simple-reconciliation-service/internal/app/service/sample"

	"github.com/google/wire"
)

func NewServices(
	svcSample sample.Service,
) *Services {
	return &Services{
		SvcSample: svcSample,
	}
}

var Set = wire.NewSet(
	sample.Set,
	NewServices,
)
