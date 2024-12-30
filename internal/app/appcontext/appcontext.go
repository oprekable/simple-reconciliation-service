package appcontext

import (
	"context"
	"embed"
	"fmt"
	"os"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/server"
	"simple-reconciliation-service/internal/app/service"
	"simple-reconciliation-service/internal/pkg/shutdown"
	"simple-reconciliation-service/internal/pkg/utils/atexit"
	"simple-reconciliation-service/internal/pkg/utils/log"

	"golang.org/x/sync/errgroup"
)

type AppContext struct {
	ctx          context.Context
	ctxCancel    context.CancelFunc
	eg           *errgroup.Group
	embedFS      *embed.FS
	repositories *repository.Repositories
	services     *service.Services
	components   *component.Components
	servers      *server.Server
}

var _ IAppContext = (*AppContext)(nil)

// NewAppContext initiate AppContext object
func NewAppContext(
	ctx context.Context,
	embedFS *embed.FS,
	repository *repository.Repositories,
	services *service.Services,
	components *component.Components,
	servers *server.Server,
) (*AppContext, func()) {
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)

	return &AppContext{
		ctx:          ctx,
		ctxCancel:    cancel,
		eg:           eg,
		embedFS:      embedFS,
		repositories: repository,
		services:     services,
		components:   components,
		servers:      servers,
	}, cancel
}

func (a *AppContext) GetCtx() context.Context {
	if a.components != nil || a.components.Logger != nil {
		return a.components.Logger.GetCtx()
	}

	return a.ctx
}

func (a *AppContext) GetComponents() *component.Components {
	return a.components
}

func (a *AppContext) Start() {
	atexit.Add(a.Shutdown)
	a.eg.Go(func() error {
		log.Msg(a.GetCtx(), "[start] application")
		return shutdown.TermSignalTrap().Wait(a.ctx, func() {
			defer func() {
				if r := recover(); r != nil {
					errRecovery := fmt.Errorf("recovered from panic: %s", r)
					log.AddErr(context.Background(), errRecovery)
					return
				}
			}()

			atexit.AtExit()

			if context.Cause(a.ctx).Error() == "done" {
				os.Exit(0)
			}

			if a.ctx.Err() != nil {
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		})
	})

	if a.servers != nil {
		a.servers.Run(a.eg)
	}

	if err := a.eg.Wait(); err != nil {
		log.Err(a.GetCtx(), "[shutdown] application", err)
	}
}

func (a *AppContext) Shutdown() {
	log.Msg(a.GetCtx(), "[shutdown] application")
}
