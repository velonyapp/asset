package usecase

import (
	"context"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type GetImage struct {
	ImageID string
}

type GetImageResult struct {
	Image *common.ImageResult
}

type GetImageHandler struct {
	imageRepo repo.Image
}

func NewGetImageHandler(
	imageRepo repo.Image,
) *GetImageHandler {
	return &GetImageHandler{
		imageRepo: imageRepo,
	}
}

func (h *GetImageHandler) Execute(
	ctx context.Context,
	uc *GetImage,
) (*GetImageResult, error) {
	imageID, err := vo.NewImageID(uc.ImageID)
	if err != nil {
		return nil, err
	}

	image, err := h.imageRepo.FindByID(ctx, imageID)
	if err != nil {
		return nil, err
	}
	if image == nil || image.IsDeleted() {
		return nil, common.ErrImageNotFound
	}

	return &GetImageResult{
		Image: &common.ImageResult{
			ID:           image.ID().String(),
			ObjectKey:    image.ObjectKey().String(),
			ObjectExists: image.ObjectExists(),
			CreateTime:   image.CreateTime(),
		},
	}, nil
}
