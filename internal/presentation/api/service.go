package api

import (
	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/usecase"
)

const (
	imageResourcePattern = "images/{image}"
)

type Service struct {
	v1.UnimplementedAssetServiceServer

	prepareImageHandler  *usecase.PrepareImageHandler
	finalizeImageHandler *usecase.FinalizeImageHandler
}

func NewService(
	prepareImageHandler *usecase.PrepareImageHandler,
	finalizeImageHandler *usecase.FinalizeImageHandler,
) *Service {
	return &Service{
		prepareImageHandler:  prepareImageHandler,
		finalizeImageHandler: finalizeImageHandler,
	}
}
