package command

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
	Image common.ImageResult
}

func (CreateImage) resultType() CreateImageResult {
	return CreateImageResult{}
}

type CreateImageHandler Handler[CreateImage, CreateImageResult]

type createImageHandler struct {
	imageRepo  repo.Image
	unitOfWork port.UnitOfWork
}

func NewCreateImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
) CreateImageHandler {
	return &createImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
	}
}

func (h *createImageHandler) Handle(
	ctx context.Context,
	cmd CreateImage,
) (CreateImageResult, error) {
	now := time.Now()

	tags, err := vo.NewTags(cmd.Tags)
	if err != nil {
		return CreateImageResult{}, err
	}
	objectKey, err := vo.NewObjectKey(cmd.ObjectKey)
	if err != nil {
		return CreateImageResult{}, err
	}

	image := entity.NewImage(tags, objectKey, now)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return CreateImageResult{}, err
	}

	return CreateImageResult{
		Image: common.ImageResult{
			ID:           image.ID().String(),
			Tags:         image.Tags().Strings(),
			ObjectKey:    image.ObjectKey().String(),
			ObjectExists: image.ObjectExists(),
			CreateTime:   image.CreateTime(),
		},
	}, nil
}
