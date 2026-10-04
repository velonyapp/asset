//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/velonyapp/asset/internal/application"
	"github.com/velonyapp/asset/internal/conf"
	"github.com/velonyapp/asset/internal/infrastructure"
	"github.com/velonyapp/asset/internal/presentation"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	*conf.Data,
	*conf.Transport,
	*slog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		presentation.ProviderSet,
		infrastructure.ProviderSet,
		application.ProviderSet,
		newApp,
	))
}
