//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"log/slog"

	"github.com/velonyapp/asset/internal/application"
	"github.com/velonyapp/asset/internal/conf"
	"github.com/velonyapp/asset/internal/info"
	"github.com/velonyapp/asset/internal/infrastructure"
	"github.com/velonyapp/asset/internal/presentation"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	context.Context,
	*info.Service,
	*conf.Data,
	*conf.Transport,
	*conf.Observability,
	*slog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		presentation.ProviderSet,
		infrastructure.ProviderSet,
		application.ProviderSet,
		newApp,
	))
}
