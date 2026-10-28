package query

import (
	"context"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type GetImage struct {
	ImageID string
}

type GetImageResult struct {
	Image       common.ImageResult
	UploadURL   *string
	DownloadURL *string
}

func (GetImage) resultType() GetImageResult {
	return GetImageResult{}
}

type GetImageHandler Handler[GetImage, GetImageResult]

type getImageHandler struct {
	imageRepo repo.Image
	storage   port.Storage
}

func NewGetImageHandler(
	imageRepo repo.Image,
	storage port.Storage,
) GetImageHandler {
	return &getImageHandler{
		imageRepo: imageRepo,
		storage:   storage,
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

	var uploadURL, downloadURL *string

	if image.IsAwaitingUpload() {
		url, err := h.storage.PresignPut(ctx,
			image.SourceObjectKey(),
			common.PresignTTL,
		)
		if err != nil {
			return GetImageResult{}, err
		}

		uploadURL = &url
	}
	if image.IsActive() {
		url, err := h.storage.PresignGet(ctx,
			image.ObjectKey(),
			common.PresignTTL,
		)
		if err != nil {
			return GetImageResult{}, err
		}

		downloadURL = &url
	}

	return GetImageResult{
		Image: common.ImageResult{
			ID:         image.ID().String(),
			Tags:       image.Tags().Strings(),
			Metadata:   image.Metadata().Map(),
			ObjectKey:  image.ObjectKey().String(),
			State:      image.State().String(),
			Width:      image.Width().UInt32(),
			Height:     image.Height().UInt32(),
			Format:     image.Format().String(),
			CreateTime: image.CreateTime(),
		},
		UploadURL:   uploadURL,
		DownloadURL: downloadURL,
	}, nil
}
