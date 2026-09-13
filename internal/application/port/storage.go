package port

import (
	"context"
	"time"
)

type Storage interface {
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	DeleteMany(ctx context.Context, keys []string) error
	PresignGet(ctx context.Context, key string, expiresIn time.Duration) (string, error)
	PresignPut(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}
