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
	command.NewConfirmImageUploadHandler,
	command.NewProcessImageHandler,
	command.NewDeleteImageHandler,
	query.NewBus,
	query.NewGetImageHandler,
	query.NewPresignImageHandler,
	domainevent.NewDispatcher,
	domainevent.NewImageCreatedHandler,
	domainevent.NewImageUpdatedHandler,
	domainevent.NewImageDeletedHandler,
)
