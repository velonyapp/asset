package api

import (
	"context"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/usecase"
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

func (s *Service) PrepareImage(ctx context.Context, req *v1.PrepareImageRequest) (*v1.PrepareImageResponse, error) {
	result, err := s.prepareImageHandler.Execute(ctx, &usecase.PrepareImage{
		Key: req.Key,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.PrepareImageResponse{
		UploadUrl: result.UploadURL,
	}, nil
}
