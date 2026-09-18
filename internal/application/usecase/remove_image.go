package usecase

import (
	"context"

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
	storageKey, err := vo.NewStorageKey(uc.StorageKey)
	if err != nil {
		return nil, err
	}

	image, err := h.imageRepo.FindByStorageKey(ctx, storageKey)
	if err != nil {
		return nil, err
	}

	if image == nil {
		return &RemoveImageResult{}, nil
	}

	if err := image.Delete(); err != nil {
		return nil, err
	}

	if err := h.storage.Delete(ctx, uc.StorageKey); err != nil {
		return nil, err
	}

	return &RemoveImageResult{}, nil
}
