package application

import (
	"github.com/velonyapp/asset/internal/application/domainevent"
	"github.com/velonyapp/asset/internal/application/usecase"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	usecase.NewUploadImageHandler,
	usecase.NewPresignImageHandler,
	usecase.NewRemoveImageHandler,
	domainevent.NewDispatcher,
	domainevent.NewImageCreatedHandler,
	domainevent.NewImageDeletedHandler,
)
