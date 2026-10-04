package transport

import (
	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/conf"
	"github.com/velonyapp/asset/internal/presentation/api"
	"github.com/velonyapp/asset/internal/presentation/middleware"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

func NewHTTPServer(
	c *conf.Transport,
	service *api.Service,
	tracing middleware.Tracing,
	metrics middleware.Metrics,
	validation middleware.Validation,
) *http.Server {
	opts := []http.ServerOption{
		http.Address(c.Http.Address),
		http.Middleware(
			recovery.Recovery(),
			tracing.Middleware(),
			metrics.Middleware(),
			validation.Middleware(),
		),
	}

	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	v1.RegisterAssetServiceHTTPServer(srv, service)

	return srv
}
