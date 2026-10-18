package presentation

import (
	"github.com/velonyapp/asset/internal/presentation/middleware"
	"github.com/velonyapp/asset/internal/presentation/service"
	"github.com/velonyapp/asset/internal/presentation/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewAPI,
	service.NewJob,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewRabbitMQConsumer,
	middleware.NewTracing,
	middleware.NewMetrics,
	middleware.NewErrorMapper,
	middleware.NewValidation,
)
