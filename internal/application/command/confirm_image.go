package command

import (
	"context"
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrSourceObjectNotFound = errors.New("source object not found")
)

type ConfirmImage struct {
	ImageID string
}

type ConfirmImageResult struct {
	Image common.ImageResult
}

func (ConfirmImage) resultType() ConfirmImageResult {
	return ConfirmImageResult{}
}

type ConfirmImageHandler Handler[ConfirmImage, ConfirmImageResult]

type confirmImageHandler struct {
	imageRepo  repo.Image
	unitOfWork port.UnitOfWork
	storage    port.Storage
}

func NewConfirmImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
) ConfirmImageHandler {
	return &confirmImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
		storage:    storage,
	}
}

func (h *confirmImageHandler) Handle(
	ctx context.Context,
	cmd ConfirmImage,
) (ConfirmImageResult, error) {
	now := time.Now().UTC()

	imageID, _ := vo.NewImageID(cmd.ImageID)

	var image *entity.Image

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		var err error

		image, err = h.imageRepo.FindByID(ctx, imageID)
		if err != nil {
			return err
		}
		if image == nil {
			return common.ErrImageNotFound
		}

		if err := image.CanConfirm(); err != nil {
			return err
		}

		present, err := h.storage.Exists(ctx, image.SourceObjectKey())
		if err != nil {
			return err
		}
		if !present {
			return ErrSourceObjectNotFound
		}

		if err := image.Confirm(now); err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return ConfirmImageResult{}, err
	}

	return ConfirmImageResult{
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
