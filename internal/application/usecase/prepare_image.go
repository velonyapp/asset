package usecase

import (
	"context"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
)

type PrepareImage struct {
	FullName string
}

type PrepareImageResult struct {
	AccessToken  string
	RefreshToken string
}

type PrepareImageHandler struct {
	imageRepo repo.Image
}

func NewPrepareImageHandler(
	imageRepo repo.Image,
	cache port.Cache,
) *PrepareImageHandler {
	return &PrepareImageHandler{
		imageRepo: imageRepo,
	}
}

func (h *PrepareImageHandler) Execute(
	ctx context.Context,
	cmd *PrepareImage,
) (*PrepareImageResult, error) {

	return &PrepareImageResult{}, nil
}
