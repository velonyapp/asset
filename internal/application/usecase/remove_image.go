package usecase

import (
	"context"

	"github.com/velonyapp/asset/internal/application/port"
)

type RemoveImage struct {
	StorageKey string
}

type RemoveImageHandler struct {
	storage port.Storage
}

func NewRemoveImageHandler(
	storage port.Storage,
) *RemoveImageHandler {
	return &RemoveImageHandler{
		storage: storage,
	}
}

func (h *RemoveImageHandler) Execute(
	ctx context.Context,
	uc *RemoveImage,
) error {
	if err := h.storage.Delete(ctx, uc.StorageKey); err != nil {
		return err
	}

	return nil
}
