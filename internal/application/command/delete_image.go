package command

import (
	"context"
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/entity"
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
	now := time.Now().UTC()

	imageID, _ := vo.NewImageID(cmd.ImageID)

	image, err := h.imageRepo.FindByID(ctx, imageID)
	if err != nil {
		return DeleteImageResult{}, err
	}
	if image == nil {
		return DeleteImageResult{}, common.ErrImageNotFound
	}

	if err := image.CanDelete(); err != nil {
		if errors.Is(err, entity.ErrImageAlreadyDeleted) {
			return DeleteImageResult{}, nil
		}

		return DeleteImageResult{}, err
	}

	if err := h.storage.DeleteMany(ctx, []vo.ObjectKey{
		image.SourceObjectKey(),
		image.ObjectKey(),
	}); err != nil {
		return DeleteImageResult{}, err
	}

	if err := image.Delete(now); err != nil {
		return DeleteImageResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return DeleteImageResult{}, err
	}

	return DeleteImageResult{}, nil
}
