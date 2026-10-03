package query

import (
	"context"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type PresignImage struct {
	ImageID string
	TTL     time.Duration
}

type PresignImageResult struct {
	UploadURL string
}

func (*PresignImage) resultType() *PresignImageResult {
	return nil
}

type PresignImageHandler Handler[*PresignImage, *PresignImageResult]

type presignImageHandler struct {
	imageRepo repo.Image
	storage   port.Storage
}

func NewPresignImageHandler(
	imageRepo repo.Image,
	storage port.Storage,
) PresignImageHandler {
	return &presignImageHandler{
		imageRepo: imageRepo,
		storage:   storage,
	}
}

func (h *presignImageHandler) Handle(
	ctx context.Context,
	qry *PresignImage,
) (*PresignImageResult, error) {
	imageID, err := vo.NewImageID(qry.ImageID)
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

	uploadURL, err := h.storage.PresignPut(ctx, image.ObjectKey(), qry.TTL)
	if err != nil {
		return nil, err
	}

	return &PresignImageResult{
		UploadURL: uploadURL,
	}, nil
}
