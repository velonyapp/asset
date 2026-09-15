package usecase

import (
	"context"
	"net/url"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
)

type PresignImage struct {
	StorageKey string

	ResizeOptions *port.ImageResizeOptions

	ExpireTime time.Time
}

type PresignImageResult struct {
	UploadURL string
}

type PresignImageHandler struct {
	uploadImageToken port.UploadImageToken
}

func NewPresignImageHandler(
	uploadImageToken port.UploadImageToken,
) *PresignImageHandler {
	return &PresignImageHandler{
		uploadImageToken: uploadImageToken,
	}
}

func (h *PresignImageHandler) Execute(
	ctx context.Context,
	uc *PresignImage,
) (*PresignImageResult, error) {
	token, err := h.uploadImageToken.Sign(port.UploadImageTokenPayload{
		StorageKey:    uc.StorageKey,
		ResizeOptions: uc.ResizeOptions,
		ExpireTime:    uc.ExpireTime,
	})
	if err != nil {
		return nil, err
	}

	uploadURL := url.URL{
		Scheme: "http",
		Host:   "localhost:8010",
		Path:   "/v1/images:upload",
	}

	query := uploadURL.Query()
	query.Set("token", token)
	uploadURL.RawQuery = query.Encode()

	return &PresignImageResult{
		UploadURL: uploadURL.String(),
	}, nil
}
