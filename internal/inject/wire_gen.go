// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"context"
	"embed"
	"simple-reconciliation-service/internal/app/appcontext"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/cerror"
	"simple-reconciliation-service/internal/app/component/clogger"
	"simple-reconciliation-service/internal/app/err/core"
	"simple-reconciliation-service/internal/app/handler/hcli"
	"simple-reconciliation-service/internal/app/repository"
	"simple-reconciliation-service/internal/app/server"
	"simple-reconciliation-service/internal/app/server/cli"
	"simple-reconciliation-service/internal/app/service"
)

// Injectors from inject.go:

func WireApp(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType) (*appcontext.AppContext, func(), error) {
	repositories := repository.NewRepositories()
	services := service.NewServices()
	configPaths := _wireConfigPathsValue
	config, err := cconfig.NewConfig(ctx, embedFS, configPaths, appName, tz)
	if err != nil {
		return nil, nil, err
	}
	logger := clogger.NewLogger(ctx)
	erType := cerror.ProvideErType(errType)
	cerrorError := cerror.NewError(erType)
	components := component.NewComponents(config, logger, cerrorError)
	v := hcli.ProviderHandlers()
	cliCli, err := cli.NewCli(components, services, repositories, v)
	if err != nil {
		return nil, nil, err
	}
	serverServer := server.NewServer(cliCli)
	appContext := appcontext.NewAppContext(ctx, embedFS, repositories, services, components, serverServer)
	return appContext, func() {
	}, nil
}

var (
	_wireConfigPathsValue = cconfig.ConfigPaths([]string{
		"./*.toml",
		"./params/*.toml",
	})
)
