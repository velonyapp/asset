package command

import (
	"context"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type ReconcileImage struct {
	ImageID string
}

type ReconcileImageResult struct {
	Image *common.ImageResult
}

func (*ReconcileImage) resultType() *ReconcileImageResult {
	return nil
}

type ReconcileImageHandler Handler[*ReconcileImage, *ReconcileImageResult]

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
	cmd *ReconcileImage,
) (*ReconcileImageResult, error) {
	now := time.Now()

	imageID, err := vo.NewImageID(cmd.ImageID)
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

	imageObjectExists, err := h.storage.Exists(ctx, image.ObjectKey())
	if err != nil {
		return nil, err
	}

	image.UpdateObjectExistence(imageObjectExists, now)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return nil, err
	}

	return &ReconcileImageResult{}, nil
}
