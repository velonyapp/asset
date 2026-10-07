package command

import (
	"context"
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/service"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrObjectAlreadyExists       = errors.New("object already exists")
	ErrSourceObjectAlreadyExists = errors.New("source object already exists")
)

type CreateImage struct {
	Tags      []string
	ObjectKey string
}

type CreateImageResult struct {
	Image common.ImageResult
}

func (CreateImage) resultType() CreateImageResult {
	return CreateImageResult{}
}

type CreateImageHandler Handler[CreateImage, CreateImageResult]

type createImageHandler struct {
	imageRepo       repo.Image
	objectKeyPolicy *service.ObjectKeyPolicy
	unitOfWork      port.UnitOfWork
	storage         port.Storage
}

func NewCreateImageHandler(
	imageRepo repo.Image,
	objectKeyPolicy *service.ObjectKeyPolicy,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
) CreateImageHandler {
	return &createImageHandler{
		imageRepo:       imageRepo,
		objectKeyPolicy: objectKeyPolicy,
		unitOfWork:      unitOfWork,
		storage:         storage,
	}
}

func (h *createImageHandler) Handle(
	ctx context.Context,
	cmd CreateImage,
) (CreateImageResult, error) {
	now := time.Now().UTC()

	tags, err := vo.NewTags(cmd.Tags)
	if err != nil {
		return CreateImageResult{}, err
	}
	objectKey, err := vo.NewObjectKey(cmd.ObjectKey)
	if err != nil {
		return CreateImageResult{}, err
	}

	// make generation port later
	sourceObjectKey := vo.ObjectKey{}

	var image *entity.Image

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.objectKeyPolicy.CanUse(ctx, objectKey); err != nil {
			if errors.Is(err, service.ErrObjectKeyAlreadyExists) {
				return ErrObjectAlreadyExists
			}
			return err
		}

		objectPresent, err := h.storage.Exists(ctx, objectKey)
		if err != nil {
			return err
		}
		if objectPresent {
			return ErrObjectAlreadyExists
		}

		if err := h.objectKeyPolicy.CanUse(ctx, sourceObjectKey); err != nil {
			if errors.Is(err, service.ErrObjectKeyAlreadyExists) {
				return ErrSourceObjectAlreadyExists
			}
			return err
		}

		sourceObjectPresent, err := h.storage.Exists(ctx, sourceObjectKey)
		if err != nil {
			return err
		}
		if sourceObjectPresent {
			return ErrSourceObjectAlreadyExists
		}

		image, err = entity.NewImage(
			tags,
			sourceObjectKey,
			objectKey,
			now,
		)
		if err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return CreateImageResult{}, err
	}

	return CreateImageResult{
		Image: common.ImageResult{
			ID:         image.ID().String(),
			Tags:       image.Tags().Strings(),
			ObjectKey:  image.ObjectKey().String(),
			State:      image.State().String(),
			CreateTime: image.CreateTime(),
		},
	}, nil
}
