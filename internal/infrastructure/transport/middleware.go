package transport

import (
	"github.com/velonyapp/asset/internal/infrastructure/observability"

	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/protobuf/proto"
)

type MetricsMiddleware middleware.Middleware

func NewMetricsMiddleware(serverMetrics *observability.ServerMetrics) MetricsMiddleware {
	return MetricsMiddleware(
		metrics.Server(
			metrics.WithSeconds(serverMetrics.Seconds),
			metrics.WithRequests(serverMetrics.Requests),
		),
	)
}

type TracesMiddleware middleware.Middleware

func NewTracesMiddleware() TracesMiddleware {
	return TracesMiddleware(tracing.Server())
}

type ValidationMiddleware middleware.Middleware

func NewValidationMiddleware() ValidationMiddleware {
	return ValidationMiddleware(
		validate.Validator(func(req any) error {
			message, ok := req.(proto.Message)
			if !ok {
				return nil
			}

			return fieldbehavior.ValidateRequiredFields(message)
		}),
	)
}
