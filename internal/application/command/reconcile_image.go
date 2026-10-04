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

type ReconcileImage struct {
	ImageID string
}

type ReconcileImageResult struct {
	Image common.ImageResult
}

func (ReconcileImage) resultType() ReconcileImageResult {
	return ReconcileImageResult{}
}

type ReconcileImageHandler Handler[ReconcileImage, ReconcileImageResult]

type reconcileImageHandler struct {
	imageRepo  repo.Image
	unitOfWork port.UnitOfWork
	storage    port.Storage
}

func NewReconcileImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
) ReconcileImageHandler {
	return &reconcileImageHandler{
		imageRepo:  imageRepo,
		unitOfWork: unitOfWork,
		storage:    storage,
	}
}

func (h *reconcileImageHandler) Handle(
	ctx context.Context,
	cmd ReconcileImage,
) (ReconcileImageResult, error) {
	now := time.Now()

	imageID, err := vo.NewImageID(cmd.ImageID)
	if err != nil {
		return ReconcileImageResult{}, err
	}

	image, err := h.imageRepo.FindByID(ctx, imageID)
	if err != nil {
		return ReconcileImageResult{}, err
	}
	if image.IsDeleted() {
		return ReconcileImageResult{}, entity.ErrImageDeleted
	}

	imageObjectExists, err := h.storage.Exists(ctx, image.ObjectKey())
	if err != nil {
		return ReconcileImageResult{}, err
	}

	if err := image.UpdateObjectExistence(imageObjectExists, now); err != nil {
		return ReconcileImageResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return ReconcileImageResult{}, err
	}

	return ReconcileImageResult{
		Image: common.ImageResult{
			ID:           image.ID().String(),
			Tags:         image.Tags().Strings(),
			ObjectKey:    image.ObjectKey().String(),
			ObjectExists: image.ObjectExists(),
			CreateTime:   image.CreateTime(),
		},
	}, nil
}
