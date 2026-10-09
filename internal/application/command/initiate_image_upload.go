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

type InitiateImageUpload struct {
	ImageID string
	TTL     time.Duration
}

type InitiateImageUploadResult struct {
	Image     common.ImageResult
	UploadURL string
}

func (InitiateImageUpload) resultType() InitiateImageUploadResult {
	return InitiateImageUploadResult{}
}

type InitiateImageUploadHandler Handler[InitiateImageUpload, InitiateImageUploadResult]

type presignImageHandler struct {
	imageRepo  repo.Image
	unitOfWork port.UnitOfWork
	storage    port.Storage
}

func NewInitiateImageUploadHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
) InitiateImageUploadHandler {
	return &presignImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
		storage:    storage,
	}
}

func (h *presignImageHandler) Handle(
	ctx context.Context,
	cmd InitiateImageUpload,
) (InitiateImageUploadResult, error) {
	now := time.Now()

	imageID, _ := vo.NewImageID(cmd.ImageID)

	var image *entity.Image
	var uploadURL string

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		var err error

		image, err = h.imageRepo.FindByID(ctx, imageID)
		if err != nil {
			return err
		}
		if image == nil {
			return common.ErrImageNotFound
		}

		if err := image.CanStartUploading(); err != nil {
			return err
		}

		uploadURL, err = h.storage.PresignPut(ctx, image.SourceObjectKey(), cmd.TTL)
		if err != nil {
			return err
		}

		if err := image.StartUploading(now); err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return InitiateImageUploadResult{}, err
	}

	return InitiateImageUploadResult{
		Image: common.ImageResult{
			ID:         image.ID().String(),
			Tags:       image.Tags().Strings(),
			ObjectKey:  image.ObjectKey().String(),
			State:      image.State().String(),
			CreateTime: image.CreateTime(),
			UpdateTime: image.UpdateTime(),
		},
		UploadURL: uploadURL,
	}, nil
}
