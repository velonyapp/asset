package infrastructure

import (
	"github.com/velonyapp/asset/internal/infrastructure/data/mysql"
	"github.com/velonyapp/asset/internal/infrastructure/data/s3"
	"github.com/velonyapp/asset/internal/infrastructure/event"
	"github.com/velonyapp/asset/internal/infrastructure/image"
	"github.com/velonyapp/asset/internal/infrastructure/observability"
	"github.com/velonyapp/asset/internal/infrastructure/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	image.NewProcessor,
	image.NewUploadTokenManager,
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewImageRepo,
	mysql.NewEventPublisher,
	s3.NewConnection,
	s3.NewStorage,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewRabbitMQConsumer,
	transport.NewTracesMiddleware,
	transport.NewMetricsMiddleware,
	transport.NewValidationMiddleware,
	observability.NewServerMetrics,
	observability.NewOpenTelemetry,
	event.NewEncoder,
)
