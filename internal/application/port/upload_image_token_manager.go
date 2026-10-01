package port

import (
	"errors"
	"time"
)

var (
	ErrInvalidUploadToken = errors.New("invalid upload token")
	ErrExpiredUploadToken = errors.New("upload token expired")
)

type UploadImageTokenPayload struct {
	Tags []string

	StorageKey string

	Transform *ImageTransform

	ExpireTime time.Time
}

type UploadImageTokenManager interface {
	Sign(payload UploadImageTokenPayload) (string, error)
	Verify(token string) (UploadImageTokenPayload, error)
}
