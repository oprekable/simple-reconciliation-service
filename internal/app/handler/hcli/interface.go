package hcli

import (
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/service"
)

type Handler interface {
	SetComponents(c *component.Components)
	SetServices(s *service.Services)
	SetRepositories(r *repository.Repositories)
	Exec() error
	Name() string
}
