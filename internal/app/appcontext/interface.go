package appcontext

import (
	"context"
	"embed"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/server"
	"simple-reconciliation-service/internal/app/service"

	"golang.org/x/sync/errgroup"
)

type IAppContext interface {
	AddToEg(in func())
	GetEmbedFS() *embed.FS
	GetCtx() context.Context
	GetCtxCancel() context.CancelFunc
	GetEg() *errgroup.Group
	GetRepositories() *repository.Repositories
	GetServices() *service.Services
	GetComponents() *component.Components
	GetServers() *server.Server
	Start()
	Shutdown()
}
