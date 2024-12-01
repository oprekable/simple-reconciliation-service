package process

import (
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"

	"github.com/google/wire"
)

func ProviderSvc(
	comp *component.Components,
	repo *repository.Repositories,
) *Svc {
	return NewSvc(comp, repo)
}

var Set = wire.NewSet(
	ProviderSvc,
	wire.Bind(new(Service), new(*Svc)),
)
