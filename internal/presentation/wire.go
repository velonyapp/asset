package presentation

import (
	"github.com/velonyapp/asset/internal/presentation/api"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
)
