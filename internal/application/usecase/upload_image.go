package usecase

import (
	"context"
	"io"

	"github.com/velonyapp/asset/internal/application/port"
)

type UploadImage struct {
	Token string
	Image io.Reader
}

type UploadImageResult struct {
	StorageKey string
}

type UploadImageHandler struct {
	storage          port.Storage
	imageProcessor   port.ImageProcessor
	uploadImageToken port.UploadImageToken
}

func NewUploadImageHandler(
	storage port.Storage,
	imageProcessor port.ImageProcessor,
	uploadImageToken port.UploadImageToken,
) *UploadImageHandler {
	return &UploadImageHandler{
		storage:          storage,
		imageProcessor:   imageProcessor,
		uploadImageToken: uploadImageToken,
	}
}

func (h *UploadImageHandler) Execute(
	ctx context.Context,
	uc *UploadImage,
) (*UploadImageResult, error) {
	payload, err := h.uploadImageToken.Verify(uc.Token)
	if err != nil {
		return nil, err
	}

	image := uc.Image

	if payload.Transform != nil {
		image, err = h.imageProcessor.Process(image, payload.Transform)
		if err != nil {
			return nil, err
		}
	}

	if err := h.storage.Put(ctx, payload.StorageKey, image); err != nil {
		return nil, err
	}

	return &UploadImageResult{
		StorageKey: payload.StorageKey,
	}, nil
}
