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

type CreateImage struct {
	Tags      []string
	Metadata  map[string]string
	ObjectKey string
	Width     int32
	Height    int32
	Format    string
}

type CreateImageResult struct {
	Image     common.ImageResult
	UploadURL string
}

func (CreateImage) resultType() CreateImageResult {
	return CreateImageResult{}
}

type CreateImageHandler Handler[CreateImage, CreateImageResult]

type createImageHandler struct {
	imageRepo                repo.Image
	unitOfWork               port.UnitOfWork
	storage                  port.Storage
	sourceObjectKeyGenerator port.SourceObjectKeyGenerator
}

func NewCreateImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
	sourceObjectKeyGenerator port.SourceObjectKeyGenerator,
) CreateImageHandler {
	return &createImageHandler{
		imageRepo:                imageRepo,
		unitOfWork:               unitOfWork,
		storage:                  storage,
		sourceObjectKeyGenerator: sourceObjectKeyGenerator,
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
	metadata, err := vo.NewMetadata(cmd.Metadata)
	if err != nil {
		return CreateImageResult{}, err
	}
	objectKey, err := vo.NewObjectKey(cmd.ObjectKey)
	if err != nil {
		return CreateImageResult{}, err
	}
	width, err := vo.NewImageWidth(cmd.Width)
	if err != nil {
		return CreateImageResult{}, err
	}
	height, err := vo.NewImageHeight(cmd.Height)
	if err != nil {
		return CreateImageResult{}, err
	}
	format, err := vo.NewImageFormat(cmd.Format)
	if err != nil {
		return CreateImageResult{}, err
	}

	sourceObjectKey, err := h.sourceObjectKeyGenerator.Generate()
	if err != nil {
		return CreateImageResult{}, err
	}

	image, err := entity.NewImage(
		tags,
		metadata,
		sourceObjectKey,
		objectKey,
		width,
		height,
		format,
		now,
	)
	if err != nil {
		return CreateImageResult{}, err
	}

	uploadURL, err := h.storage.PresignPut(ctx,
		sourceObjectKey,
		common.PresignTTL,
	)
	if err != nil {
		return CreateImageResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return CreateImageResult{}, err
	}

	return CreateImageResult{
		Image: common.ImageResult{
			ID:         image.ID().String(),
			Tags:       image.Tags().Strings(),
			Metadata:   image.Metadata().Map(),
			ObjectKey:  image.ObjectKey().String(),
			State:      image.State().String(),
			Width:      image.Width().Int32(),
			Height:     image.Height().Int32(),
			Format:     image.Format().String(),
			CreateTime: image.CreateTime(),
		},
		UploadURL: uploadURL,
	}, nil
}
