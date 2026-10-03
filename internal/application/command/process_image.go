package command

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var ErrImageObjectDoesntExist = errors.New("image object doesn't exist")

type ProcessImage struct {
	ImageID string
	Options port.ImageProcessOptions
}

type ProcessImageResult struct{}

func (*ProcessImage) resultType() *ProcessImageResult {
	return nil
}

type ProcessImageHandler Handler[*ProcessImage, *ProcessImageResult]

type processImageHandler struct {
	imageRepo      repo.Image
	eventPublisher port.EventPublisher
	storage        port.Storage
	imageProcessor port.ImageProcessor
}

func NewProcessImageHandler(
	imageRepo repo.Image,
	eventPublisher port.EventPublisher,
	storage port.Storage,
	imageProcessor port.ImageProcessor,
) ProcessImageHandler {
	return &processImageHandler{
		imageRepo:      imageRepo,
		eventPublisher: eventPublisher,
		storage:        storage,
		imageProcessor: imageProcessor,
	}
}

func (h *processImageHandler) Handle(ctx context.Context, cmd *ProcessImage) (*ProcessImageResult, error) {
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
	if !image.ObjectExists() {
		return nil, ErrImageObjectDoesntExist
	}

	src, err := h.storage.Get(ctx, image.ObjectKey())
	if err != nil {
		return nil, err
	}
	defer src.Close()

	pr, pw := io.Pipe()

	processErrCh := make(chan error, 1)

	go func() {
		err := h.imageProcessor.Process(src, pw, cmd.Options)
		_ = pw.CloseWithError(err)
		processErrCh <- err
	}()

	putErr := h.storage.Put(ctx, image.ObjectKey(), pr)
	if putErr != nil {
		_ = pr.CloseWithError(putErr)
	} else {
		_ = pr.Close()
	}

	processErr := <-processErrCh
	if processErr != nil {
		return nil, processErr
	}
	if putErr != nil {
		return nil, putErr
	}

	h.eventPublisher.Publish(ctx, integrationevent.NewImageProcessed(
		image.ID().String(),
		image.Tags().Strings(),
		now,
	))

	return &ProcessImageResult{}, nil
}
