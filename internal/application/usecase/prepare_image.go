package usecase

import (
	"context"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
)

type PrepareImage struct {
	Key string
}

type PrepareImageResult struct {
	UploadURL string
}

type PrepareImageHandler struct {
	imageRepo repo.Image
	storage   port.Storage
}

func NewPrepareImageHandler(
	imageRepo repo.Image,
	storage port.Storage,
) *PrepareImageHandler {
	return &PrepareImageHandler{
		imageRepo: imageRepo,
		storage:   storage,
	}
}

func (h *PrepareImageHandler) Execute(
	ctx context.Context,
	uc *PrepareImage,
) (*PrepareImageResult, error) {
	uploadURL, err := h.storage.PresignPut(ctx, uc.Key, time.Minute*5)
	if err != nil {
		return nil, err
	}

	return &PrepareImageResult{UploadURL: uploadURL}, nil
}
