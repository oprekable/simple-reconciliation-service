package cli

import (
	"context"
	"errors"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/handler/hcli"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/service"
	"simple-reconciliation-service/internal/pkg/utils/log"

	"golang.org/x/sync/errgroup"
)

const name = "cli"

type Cli struct {
	ctx      context.Context
	comp     *component.Components
	svc      *service.Services
	repo     *repository.Repositories
	handlers []hcli.Handler
}

func NewCli(
	comp *component.Components,
	svc *service.Services,
	repo *repository.Repositories,
	handlers []hcli.Handler,
) (*Cli, error) {
	returnData := &Cli{
		ctx:      comp.Logger.GetCtx(),
		comp:     comp,
		svc:      svc,
		repo:     repo,
		handlers: handlers,
	}

	for k := range handlers {
		handlers[k].SetComponents(comp)
		handlers[k].SetServices(svc)
		handlers[k].SetRepositories(repo)
	}

	return returnData, nil
}

func (c *Cli) Name() string {
	return name
}

func (c *Cli) Start(eg *errgroup.Group) {
	eg.Go(func() (err error) {
		ctx := c.ctx

		for k, v := range c.handlers {
			if v.Name() == c.comp.Config.Action {
				err = c.handlers[k].Exec()
				if err != nil {
					log.Err(ctx, "error", err)
				} else {
					err = errors.New("done")
				}

				return err
			}
		}

		return err
	})
}

func (c *Cli) Shutdown() {
	log.Msg(c.ctx, "[shutdown] "+name)
}
