package usecase

import (
	"context"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/repo"
)

type FinalizeImage struct {
	FullName string
}

type FinalizeImageResult struct {
	AccessToken  string
	RefreshToken string
}

type FinalizeImageHandler struct {
	imageRepo repo.Image
}

func NewFinalizeImageHandler(
	imageRepo repo.Image,
	cache port.Cache,
) *FinalizeImageHandler {
	return &FinalizeImageHandler{
		imageRepo: imageRepo,
	}
}

func (h *FinalizeImageHandler) Execute(
	ctx context.Context,
	cmd *FinalizeImage,
) (*FinalizeImageResult, error) {

	return &FinalizeImageResult{}, nil
}
