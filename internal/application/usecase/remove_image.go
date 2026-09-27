package usecase

import (
	"context"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type RemoveImage struct {
	StorageKey string
}

type RemoveImageResult struct {
}

type RemoveImageHandler struct {
	imageRepo  repo.Image
	unitOfWork port.UnitOfWork
	storage    port.Storage
}

func NewRemoveImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
) *RemoveImageHandler {
	return &RemoveImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
		storage:    storage,
	}
}

func (h *RemoveImageHandler) Execute(
	ctx context.Context,
	uc *RemoveImage,
) (*RemoveImageResult, error) {
	now := time.Now()

	storageKey, err := vo.NewStorageKey(uc.StorageKey)
	if err != nil {
		return nil, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		image, err := h.imageRepo.FindByStorageKey(ctx, storageKey)
		if err != nil {
			return err
		}

		if image == nil || image.IsDeleted() {
			return nil
		}

		if err := image.Delete(now); err != nil {
			return err
		}

		if err := h.storage.Delete(ctx, storageKey); err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return nil, err
	}

	return &RemoveImageResult{}, nil
}
