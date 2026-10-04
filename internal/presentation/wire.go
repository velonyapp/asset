package presentation

import (
	"github.com/velonyapp/asset/internal/presentation/api"
	"github.com/velonyapp/asset/internal/presentation/middleware"
	"github.com/velonyapp/asset/internal/presentation/observability"
	"github.com/velonyapp/asset/internal/presentation/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewRabbitMQConsumer,
	observability.NewServerMetrics,
	middleware.NewTracingMiddleware,
	middleware.NewMetricsMiddleware,
	middleware.NewErrorMapperMiddleware,
	middleware.NewValidationMiddleware,
)
