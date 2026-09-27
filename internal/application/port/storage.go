package port

import (
	"context"
	"io"
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

type Storage interface {
	Exists(ctx context.Context, key vo.StorageKey) (bool, error)

	Get(ctx context.Context, key vo.StorageKey) (io.ReadCloser, error)
	Put(ctx context.Context, key vo.StorageKey, body io.Reader) error

	Delete(ctx context.Context, key vo.StorageKey) error
	DeleteMany(ctx context.Context, keys []vo.StorageKey) error

	PresignGet(ctx context.Context, key vo.StorageKey, expiresIn time.Duration) (string, error)
	PresignPut(ctx context.Context, key vo.StorageKey, expiresIn time.Duration) (string, error)
}
