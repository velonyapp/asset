package port

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Exists(ctx context.Context, key string) (bool, error)

	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Put(ctx context.Context, key string, body io.Reader) error

	Delete(ctx context.Context, key string) error
	DeleteMany(ctx context.Context, keys []string) error

	PresignGet(ctx context.Context, key string, expiresIn time.Duration) (string, error)
	PresignPut(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}