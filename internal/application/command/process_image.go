package command

import (
	"context"
	"io"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type ProcessImage struct {
	ImageID string
	Options port.ImageProcessOptions
}

type ProcessImageResult struct {
	Image common.ImageResult
}

func (ProcessImage) resultType() ProcessImageResult {
	return ProcessImageResult{}
}

type ProcessImageHandler Handler[ProcessImage, ProcessImageResult]

type processImageHandler struct {
	imageRepo      repo.Image
	unitOfWork     port.UnitOfWork
	storage        port.Storage
	imageProcessor port.ImageProcessor
}

func NewProcessImageHandler(
	imageRepo repo.Image,
	unitOfWork port.UnitOfWork,
	storage port.Storage,
	imageProcessor port.ImageProcessor,
) ProcessImageHandler {
	return &processImageHandler{
		imageRepo:      imageRepo,
		unitOfWork:     unitOfWork,
		storage:        storage,
		imageProcessor: imageProcessor,
	}
}

func (h *processImageHandler) Handle(
	ctx context.Context,
	cmd ProcessImage,
) (ProcessImageResult, error) {
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

		if err := image.CanProcess(); err != nil {
			return err
		}

		src, err := h.storage.Get(ctx, image.SourceObjectKey())
		if err != nil {
			return err
		}
		defer src.Close()

		pr, pw := io.Pipe()

		processErrChan := make(chan error, 1)

		go func() {
			err := h.imageProcessor.Process(src, pw, cmd.Options)
			_ = pw.CloseWithError(err)
			processErrChan <- err
		}()

		putErr := h.storage.Put(ctx, image.ObjectKey(), pr)
		if putErr != nil {
			_ = pr.CloseWithError(putErr)
		} else {
			_ = pr.Close()
		}

		processErr := <-processErrChan
		if processErr != nil {
			return processErr
		}
		if putErr != nil {
			return putErr
		}

		now := time.Now().UTC()

		if err := image.Process(now); err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return ProcessImageResult{}, err
	}

	return ProcessImageResult{
		Image: common.ImageResult{
			ID:         image.ID().String(),
			Tags:       image.Tags().Strings(),
			ObjectKey:  image.ObjectKey().String(),
			State:      image.State().String(),
			CreateTime: image.CreateTime(),
		},
	}, nil
}
