package usecase

import (
	"context"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type CreateImage struct {
	Tags      []string
	ObjectKey string
}

type CreateImageResult struct {
	Image *common.ImageResult
}

type CreateImageHandler struct {
	imageRepo  repo.Image
	unitOfWork port.UnitOfWork
}

func NewCreateImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
) *CreateImageHandler {
	return &CreateImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
	}
}

func (h *CreateImageHandler) Execute(
	ctx context.Context,
	uc *CreateImage,
) (*CreateImageResult, error) {
	now := time.Now()

	tags, err := vo.NewTags(uc.Tags)
	if err != nil {
		return nil, err
	}
	objectKey, err := vo.NewObjectKey(uc.ObjectKey)
	if err != nil {
		return nil, err
	}

	image := entity.NewImage(tags, objectKey, now)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return nil, err
	}

	return &CreateImageResult{
		Image: &common.ImageResult{
			ID:           image.ID().String(),
			ObjectKey:    image.ObjectKey().String(),
			ObjectExists: image.ObjectExists(),
			CreateTime:   image.CreateTime(),
		},
	}, nil
}
