package usecase

import (
	"context"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type DeleteImage struct {
	ImageID string
}

type DeleteImageResult struct {
}

type DeleteImageHandler struct {
	imageRepo  repo.Image
	unitOfWork port.UnitOfWork
	storage    port.Storage
}

func NewDeleteImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
) *DeleteImageHandler {
	return &DeleteImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
		storage:    storage,
	}
}

func (h *DeleteImageHandler) Execute(
	ctx context.Context,
	uc *DeleteImage,
) (*DeleteImageResult, error) {
	now := time.Now()

	imageID := vo.NewImageID(uc.ImageID)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		image, err := h.imageRepo.FindByID(ctx, imageID)
		if err != nil {
			return err
		}
		if image == nil || image.IsDeleted() {
			return nil
		}

		image.Delete(now)

		if err := h.storage.Delete(ctx, image.ObjectKey()); err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return nil, err
	}

	return &DeleteImageResult{}, nil
}
