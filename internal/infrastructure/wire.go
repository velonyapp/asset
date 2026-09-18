package infrastructure

import (
	"github.com/velonyapp/asset/internal/infrastructure/data/mysql"
	"github.com/velonyapp/asset/internal/infrastructure/data/s3"
	"github.com/velonyapp/asset/internal/infrastructure/messaging/event"
	"github.com/velonyapp/asset/internal/infrastructure/observability"
	"github.com/velonyapp/asset/internal/infrastructure/service"
	"github.com/velonyapp/asset/internal/infrastructure/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewImageProcessor,
	service.NewUploadImageToken,
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewImageRepo,
	mysql.NewOutboxPublisher,
	s3.NewConnection,
	s3.NewStorage,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewTracesMiddleware,
	transport.NewMetricsMiddleware,
	transport.NewValidationMiddleware,
	observability.NewServerMetrics,
	observability.NewOpenTelemetry,
	event.NewEncoder,
)
