package usecase

import (
	"context"
	"net/url"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/conf"
)

type PresignImage struct {
	StorageKey string

	Transform *port.ImageTransform

	ExpireTime time.Time
}

type PresignImageResult struct {
	UploadURL string
}

type PresignImageHandler struct {
	c                *conf.Service
	uploadImageToken port.UploadImageToken
}

func NewPresignImageHandler(
	c *conf.Service,
	uploadImageToken port.UploadImageToken,
) *PresignImageHandler {
	return &PresignImageHandler{
		c:                c,
		uploadImageToken: uploadImageToken,
	}
}

func (h *PresignImageHandler) Execute(
	ctx context.Context,
	uc *PresignImage,
) (*PresignImageResult, error) {
	publicURL, err := url.Parse(h.c.PublicUrl)
	if err != nil {
		return nil, err
	}

	token, err := h.uploadImageToken.Sign(port.UploadImageTokenPayload{
		StorageKey: uc.StorageKey,
		Transform:  uc.Transform,
		ExpireTime: uc.ExpireTime,
	})
	if err != nil {
		return nil, err
	}

	uploadURL := publicURL.JoinPath("v1:uploadImage")

	query := uploadURL.Query()
	query.Set("token", token)
	uploadURL.RawQuery = query.Encode()

	return &PresignImageResult{
		UploadURL: uploadURL.String(),
	}, nil
}
