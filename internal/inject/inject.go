//go:build wireinject
// +build wireinject

package inject

import (
	"embed"
	"simple-reconciliation-service/internal/app/appcontext"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/clogger"
	"simple-reconciliation-service/internal/app/err/core"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/server"
	"simple-reconciliation-service/internal/app/service"

	"context"

	"github.com/google/wire"
)

func WireApp(
	ctx context.Context,
	embedFS *embed.FS,
	appName cconfig.AppName,
	tz cconfig.TimeZone,
	errType []core.ErrorType,
	isShowLog clogger.IsShowLog,
) (*appcontext.AppContext, func(), error) {
	wire.Build(
		component.Set,
		repository.Set,
		service.Set,
		server.Set,
		appcontext.NewAppContext,
	)

	return nil, func() {
		// do nothing
	}, nil
}
