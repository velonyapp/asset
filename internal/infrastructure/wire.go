package infrastructure

import (
	"github.com/velonyapp/asset/internal/infrastructure/data/mysql"
	"github.com/velonyapp/asset/internal/infrastructure/data/redis"
	"github.com/velonyapp/asset/internal/infrastructure/data/s3"
	"github.com/velonyapp/asset/internal/infrastructure/observability"
	"github.com/velonyapp/asset/internal/infrastructure/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewImageRepo,
	redis.NewConnection,
	redis.NewCache,
	s3.NewConnection,
	s3.NewStorage,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewTracesMiddleware,
	transport.NewMetricsMiddleware,
	transport.NewValidationMiddleware,
	observability.NewServerMetrics,
	observability.NewOpenTelemetry,
)
