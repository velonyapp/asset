package middleware

import (
	"github.com/velonyapp/asset/internal/presentation/observability"

	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	"github.com/go-kratos/kratos/v3/middleware"
)

type Metrics middleware.Middleware

func NewMetricsMiddleware(serverMetrics *observability.ServerMetrics) Metrics {
	return Metrics(
		metrics.Server(
			metrics.WithSeconds(serverMetrics.Seconds),
			metrics.WithRequests(serverMetrics.Requests),
		),
	)
}

func (m Metrics) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
