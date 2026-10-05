package command

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

func (DeleteImage) resultType() DeleteImageResult {
	return DeleteImageResult{}
}

type DeleteImageHandler Handler[DeleteImage, DeleteImageResult]

type deleteImageHandler struct {
	imageRepo  repo.Image
	unitOfWork port.UnitOfWork
	storage    port.Storage
}

func NewDeleteImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
) DeleteImageHandler {
	return &deleteImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
		storage:    storage,
	}
}

func (h *deleteImageHandler) Handle(
	ctx context.Context,
	cmd DeleteImage,
) (DeleteImageResult, error) {
	now := time.Now()

	imageID, err := vo.NewImageID(cmd.ImageID)
	if err != nil {
		return DeleteImageResult{}, err
	}

	image, err := h.imageRepo.GetByID(ctx, imageID)
	if err != nil {
		return DeleteImageResult{}, err
	}

	if err := h.storage.Delete(ctx, image.ObjectKey()); err != nil {
		return DeleteImageResult{}, err
	}

	image.Delete(now)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return DeleteImageResult{}, err
	}

	return DeleteImageResult{}, nil
}
