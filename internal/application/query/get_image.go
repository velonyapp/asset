package query

import (
	"context"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
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
	imageQuery Image
	storage    port.Storage
}

func NewGetImageHandler(
	imageQuery Image,
	storage port.Storage,
) GetImageHandler {
	return &getImageHandler{
		imageQuery: imageQuery,
		storage:    storage,
	}
}

func (h *getImageHandler) Handle(
	ctx context.Context,
	qry GetImage,
) (GetImageResult, error) {
	imageID, _ := vo.NewImageID(qry.ImageID)

	image, err := h.imageQuery.Find(ctx, imageID)
	if err != nil {
		return GetImageResult{}, err
	}
	if image == nil {
		return GetImageResult{}, common.ErrImageNotFound
	}

	var uploadURL, downloadURL *string

	if image.State == vo.ImageStateAwaitingUpload.String() {
		sourceObjectKey, _ := vo.NewObjectKey(image.SourceObjectKey)

		url, err := h.storage.PresignPut(ctx,
			sourceObjectKey,
			common.PresignTTL,
		)
		if err != nil {
			return GetImageResult{}, err
		}

		uploadURL = &url
	}
	if image.State == vo.ImageStateActive.String() {
		objectKey, _ := vo.NewObjectKey(image.ObjectKey)

		url, err := h.storage.PresignGet(ctx,
			objectKey,
			common.PresignTTL,
		)
		if err != nil {
			return GetImageResult{}, err
		}

		downloadURL = &url
	}

	return GetImageResult{
		Image:       *image,
		UploadURL:   uploadURL,
		DownloadURL: downloadURL,
	}, nil
}
