package api

import (
	"bytes"
	"context"
	"errors"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/application/usecase"
)

type Service struct {
	v1.UnimplementedAssetServiceServer

	uploadImageHandler  *usecase.UploadImageHandler
	presignImageHandler *usecase.PresignImageHandler
}

func NewService(
	uploadImageHandler *usecase.UploadImageHandler,
	presignImageHandler *usecase.PresignImageHandler,
) *Service {
	return &Service{
		uploadImageHandler:  uploadImageHandler,
		presignImageHandler: presignImageHandler,
	}
}

func (s *Service) PresignImage(
	ctx context.Context,
	req *v1.PresignImageRequest,
) (*v1.PresignImageResponse, error) {
	var resizeOptions *port.ImageResizeOptions

	if req.ResizeOptions != nil {
		var fit port.ImageResizeFit

		switch req.ResizeOptions.Fit {
		case v1.ImageResizeFit_IMAGE_RESIZE_FIT_CONTAIN:
			fit = port.ImageResizeFitContain
		case v1.ImageResizeFit_IMAGE_RESIZE_FIT_COVER:
			fit = port.ImageResizeFitCover
		default:
			return nil, errors.New("unsupported image resize fit")
		}

		resizeOptions = &port.ImageResizeOptions{
			Width:  req.ResizeOptions.Width,
			Height: req.ResizeOptions.Height,
			Fit:    fit,
		}
	}

	result, err := s.presignImageHandler.Execute(ctx, &usecase.PresignImage{
		StorageKey:    req.StorageKey,
		ResizeOptions: resizeOptions,
		ExpireTime:    req.ExpireTime.AsTime(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.PresignImageResponse{
		UploadUrl: result.UploadURL,
	}, nil
}

func (s *Service) UploadImage(ctx context.Context, req *v1.UploadImageRequest) (*v1.UploadImageResponse, error) {
	result, err := s.uploadImageHandler.Execute(ctx, &usecase.UploadImage{
		Token: req.Token,
		Image: bytes.NewReader(req.Image.Data),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.UploadImageResponse{
		StorageKey: result.StorageKey,
	}, nil
}
