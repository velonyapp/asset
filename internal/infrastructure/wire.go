package infrastructure

import (
	"github.com/velonyapp/asset/internal/infrastructure/data/mysql"
	"github.com/velonyapp/asset/internal/infrastructure/data/s3"
	"github.com/velonyapp/asset/internal/infrastructure/messaging/protobuf"
	"github.com/velonyapp/asset/internal/infrastructure/observability"
	"github.com/velonyapp/asset/internal/infrastructure/processing"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	processing.NewImageProcessor,
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewImageRepo,
	mysql.NewEventPublisher,
	s3.NewConnection,
	s3.NewStorage,
	observability.NewOpenTelemetry,
	observability.NewServerMetrics,
	protobuf.NewEncoder,
)
