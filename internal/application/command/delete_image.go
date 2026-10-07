package command

import (
	"context"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
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
}

func NewDeleteImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
) DeleteImageHandler {
	return &deleteImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
	}
}

func (h *deleteImageHandler) Handle(
	ctx context.Context,
	cmd DeleteImage,
) (DeleteImageResult, error) {
	now := time.Now().UTC()

	imageID, _ := vo.NewImageID(cmd.ImageID)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		image, err := h.imageRepo.FindByID(ctx, imageID)
		if err != nil {
			return err
		}
		if image == nil {
			return common.ErrImageNotFound
		}

		if err := image.Delete(now); err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return DeleteImageResult{}, err
	}

	return DeleteImageResult{}, nil
}
