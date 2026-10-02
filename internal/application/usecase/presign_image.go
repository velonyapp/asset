package usecase

import (
	"context"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type PresignImage struct {
	ImageID string
	TTL     time.Duration
}

type PresignImageResult struct {
	UploadURL string
}

type PresignImageHandler struct {
	imageRepo repo.Image
	storage   port.Storage
}

func NewPresignImageHandler(
	imageRepo repo.Image,
	storage port.Storage,
) *PresignImageHandler {
	return &PresignImageHandler{
		imageRepo: imageRepo,
		storage:   storage,
	}
}

func (h *PresignImageHandler) Execute(
	ctx context.Context,
	uc *PresignImage,
) (*PresignImageResult, error) {
	imageID := vo.NewImageID(uc.ImageID)

	image, err := h.imageRepo.FindByID(ctx, imageID)
	if err != nil {
		return nil, err
	}
	if image == nil || image.IsDeleted() {
		return nil, common.ErrImageNotFound
	}

	uploadURL, err := h.storage.PresignPut(ctx, image.ObjectKey(), uc.TTL)
	if err != nil {
		return nil, err
	}

	return &PresignImageResult{
		UploadURL: uploadURL,
	}, nil
}
