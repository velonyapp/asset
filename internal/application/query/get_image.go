package query

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
	Image common.ImageResult
}

func (GetImage) resultType() GetImageResult {
	return GetImageResult{}
}

type GetImageHandler Handler[GetImage, GetImageResult]

type getImageHandler struct {
	imageRepo repo.Image
}

func NewGetImageHandler(
	imageRepo repo.Image,
) GetImageHandler {
	return &getImageHandler{
		imageRepo: imageRepo,
	}
}

func (h *getImageHandler) Handle(
	ctx context.Context,
	qry GetImage,
) (GetImageResult, error) {
	imageID, err := vo.NewImageID(qry.ImageID)
	if err != nil {
		return GetImageResult{}, err
	}

	image, err := h.imageRepo.GetByID(ctx, imageID)
	if err != nil {
		return GetImageResult{}, err
	}

	return GetImageResult{
		Image: common.ImageResult{
			ID:           image.ID().String(),
			Tags:         image.Tags().Strings(),
			ObjectKey:    image.ObjectKey().String(),
			ObjectExists: image.ObjectExists(),
			CreateTime:   image.CreateTime(),
		},
	}, nil
}
