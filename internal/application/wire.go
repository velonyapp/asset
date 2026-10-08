package application

import (
	"github.com/velonyapp/asset/internal/application/command"
	"github.com/velonyapp/asset/internal/application/domainevent"
	"github.com/velonyapp/asset/internal/application/query"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	command.NewBus,
	command.NewCreateImageHandler,
	command.NewPresignImageHandler,
	command.NewConfirmImageUploadHandler,
	command.NewProcessImageHandler,
	command.NewDeleteImageHandler,
	query.NewBus,
	query.NewGetImageHandler,
	domainevent.NewDispatcher,
	domainevent.NewImageCreatedHandler,
	domainevent.NewImageUploadedHandler,
	domainevent.NewImageProcessedHandler,
	domainevent.NewImageDeletedHandler,
)
