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
	imageID, _ := vo.NewImageID(qry.ImageID)

	image, err := h.imageRepo.FindByID(ctx, imageID)
	if err != nil {
		return GetImageResult{}, err
	}
	if image == nil {
		return GetImageResult{}, common.ErrImageNotFound
	}

	return GetImageResult{
		Image: common.ImageResult{
			ID:         image.ID().String(),
			Tags:       image.Tags().Strings(),
			ObjectKey:  image.ObjectKey().String(),
			State:      image.State().String(),
			CreateTime: image.CreateTime(),
			UpdateTime: image.UpdateTime(),
		},
	}, nil
}
