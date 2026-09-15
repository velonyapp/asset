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
	StorageKey string

	ResizeOptions *ImageResizeOptions

	ExpireTime time.Time
}

type UploadImageToken interface {
	Sign(payload UploadImageTokenPayload) (string, error)
	Verify(token string) (UploadImageTokenPayload, error)
}
