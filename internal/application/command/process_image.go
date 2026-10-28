package command

import (
	"context"
	"errors"
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
}

type ProcessImageResult struct {
	Image       common.ImageResult
	DownloadURL string
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
	var downloadURL string

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		now := time.Now()

		var err error

		image, err = h.imageRepo.FindByIDForUpdate(ctx, imageID)
		if err != nil {
			return err
		}
		if image == nil {
			return common.ErrImageNotFound
		}

		if err := image.StartProcessing(now); err != nil {
			if errors.Is(err, entity.ErrImageAlreadyProcessing) || errors.Is(err, entity.ErrImageActive) {
				return nil
			}

			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		return ProcessImageResult{}, err
	}

	if err := func() error {
		src, err := h.storage.Get(ctx, image.SourceObjectKey())
		if err != nil {
			return err
		}
		defer src.Close()

		pr, pw := io.Pipe()

		processErrChan := make(chan error, 1)

		go func() {
			err := h.imageProcessor.Process(src, pw, image.Width(), image.Height(), image.Format())
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

		return nil
	}(); err != nil {
		// Try rolling back, failed rollback will be handled by background worker
		h.tryRollback(ctx, imageID)

		return ProcessImageResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		now := time.Now()

		var err error

		image, err = h.imageRepo.FindByIDForUpdate(ctx, imageID)
		if err != nil {
			return err
		}
		if image == nil {
			return common.ErrImageNotFound
		}

		if err := image.CompleteProcessing(now); err != nil {
			if errors.Is(err, entity.ErrImageAlreadyActive) {
				return nil
			}

			return err
		}

		downloadURL, err = h.storage.PresignGet(ctx,
			image.ObjectKey(),
			common.PresignTTL,
		)
		if err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	}); err != nil {
		// Try rolling back, failed rollback will be handled by background worker
		h.tryRollback(ctx, imageID)

		return ProcessImageResult{}, err
	}

	return ProcessImageResult{
		Image: common.ImageResult{
			ID:         image.ID().String(),
			Tags:       image.Tags().Strings(),
			Metadata:   image.Metadata().Map(),
			ObjectKey:  image.ObjectKey().String(),
			State:      image.State().String(),
			Width:      image.Width().UInt32(),
			Height:     image.Height().UInt32(),
			Format:     image.Format().String(),
			CreateTime: image.CreateTime(),
		},
		DownloadURL: downloadURL,
	}, nil
}

func (h *processImageHandler) tryRollback(ctx context.Context, imageID vo.ImageID) error {
	rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()

	return h.unitOfWork.Do(rollbackCtx, func(ctx context.Context) error {
		now := time.Now()

		image, err := h.imageRepo.FindByIDForUpdate(ctx, imageID)
		if err != nil {
			return err
		}
		if image == nil {
			return common.ErrImageNotFound
		}

		if err := image.CancelProcessing(now); err != nil {
			return err
		}

		return h.imageRepo.Save(ctx, image)
	})
}
