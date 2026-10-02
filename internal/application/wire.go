package application

import (
	"github.com/velonyapp/asset/internal/application/domainevent"
	"github.com/velonyapp/asset/internal/application/usecase"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	usecase.NewGetImageHandler,
	usecase.NewCreateImageHandler,
	usecase.NewPresignImageHandler,
	usecase.NewProcessImageHandler,
	usecase.NewReconcileImageHandler,
	usecase.NewDeleteImageHandler,
	domainevent.NewDispatcher,
	domainevent.NewImageCreatedHandler,
	domainevent.NewImageObjectExistenceUpdatedHandler,
	domainevent.NewImageDeletedHandler,
)
