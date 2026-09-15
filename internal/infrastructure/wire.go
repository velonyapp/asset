package infrastructure

import (
	"github.com/velonyapp/asset/internal/infrastructure/data/s3"
	"github.com/velonyapp/asset/internal/infrastructure/image"
	"github.com/velonyapp/asset/internal/infrastructure/observability"
	"github.com/velonyapp/asset/internal/infrastructure/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	s3.NewConnection,
	s3.NewStorage,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewTracesMiddleware,
	transport.NewMetricsMiddleware,
	transport.NewValidationMiddleware,
	observability.NewServerMetrics,
	observability.NewOpenTelemetry,
	image.NewProcessor,
	image.NewUploadToken,
)
