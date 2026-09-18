package usecase

import (
	"context"
	"io"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type UploadImage struct {
	Token string
	Image io.Reader
}

type UploadImageResult struct {
	StorageKey string
}

type UploadImageHandler struct {
	imageRepo        repo.Image
	unitOfWork       port.UnitOfWork
	storage          port.Storage
	imageProcessor   port.ImageProcessor
	uploadImageToken port.UploadImageToken
}

func NewUploadImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
	imageProcessor port.ImageProcessor,
	uploadImageToken port.UploadImageToken,
) *UploadImageHandler {
	return &UploadImageHandler{
		imageRepo:        imageRepo,
		unitOfWork:       unitOfWork,
		storage:          storage,
		imageProcessor:   imageProcessor,
		uploadImageToken: uploadImageToken,
	}
}

func (h *UploadImageHandler) Execute(
	ctx context.Context,
	uc *UploadImage,
) (*UploadImageResult, error) {
	payload, err := h.uploadImageToken.Verify(uc.Token)
	if err != nil {
		return nil, err
	}

	imageObject := uc.Image

	if payload.Transform != nil {
		imageObject, err = h.imageProcessor.Process(imageObject, payload.Transform)
		if err != nil {
			return nil, err
		}
	}

	storageKey, err := vo.NewStorageKey(payload.StorageKey)
	if err != nil {
		return nil, err
	}

	createdImage := entity.NewImage(storageKey)

	if err := h.storage.Put(ctx, storageKey.Value(), imageObject); err != nil {
		return nil, err
	}

	createdImage.Finalize()

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.imageRepo.Save(ctx, createdImage)
	}); err != nil {
		return nil, err
	}

	return &UploadImageResult{
		StorageKey: payload.StorageKey,
	}, nil
}
