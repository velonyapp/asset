package port

import (
	"context"
	"io"
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

type Storage interface {
	Exists(ctx context.Context, key vo.ObjectKey) (bool, error)

	Get(ctx context.Context, key vo.ObjectKey) (io.ReadCloser, error)
	Put(ctx context.Context, key vo.ObjectKey, body io.Reader) error

	Delete(ctx context.Context, key vo.ObjectKey) error
	DeleteMany(ctx context.Context, keys []vo.ObjectKey) error

	PresignGet(ctx context.Context, key vo.ObjectKey, ttl time.Duration) (string, error)
	PresignPut(ctx context.Context, key vo.ObjectKey, ttl time.Duration) (string, error)
}
