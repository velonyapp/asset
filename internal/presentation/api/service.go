package api

import (
	"bytes"
	"context"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/application/usecase"
)

type Service struct {
	v1.UnimplementedAssetServiceServer

	uploadImageHandler  *usecase.UploadImageHandler
	presignImageHandler *usecase.PresignImageHandler
	removeImageHandler  *usecase.RemoveImageHandler
}

func NewService(
	uploadImageHandler *usecase.UploadImageHandler,
	presignImageHandler *usecase.PresignImageHandler,
	removeImageHandler *usecase.RemoveImageHandler,
) *Service {
	return &Service{
		uploadImageHandler:  uploadImageHandler,
		presignImageHandler: presignImageHandler,
		removeImageHandler:  removeImageHandler,
	}
}

func (s *Service) PresignImage(ctx context.Context, req *v1.PresignImageRequest) (*v1.PresignImageResponse, error) {
	var transform *port.ImageTransform

	if req.Transform != nil {
		transform = &port.ImageTransform{}

		if req.Transform.Resize != nil {
			var fit port.ImageResizeFit

			switch req.Transform.Resize.Fit {
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_UNSPECIFIED:
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_CONTAIN:
				fit = port.ImageResizeFitContain
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_COVER:
				fit = port.ImageResizeFitCover
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_PAD:
				fit = port.ImageResizeFitPad
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_STRETCH:
				fit = port.ImageResizeFitStretch
			default:
				return nil, port.ErrUnsupportedResizeFit
			}

			var gravity port.ImageGravity

			switch req.Transform.Resize.Gravity {
			case v1.ImageGravity_IMAGE_GRAVITY_UNSPECIFIED:
			case v1.ImageGravity_IMAGE_GRAVITY_CENTER:
				gravity = port.ImageGravityCenter
			case v1.ImageGravity_IMAGE_GRAVITY_TOP:
				gravity = port.ImageGravityTop
			case v1.ImageGravity_IMAGE_GRAVITY_TOP_RIGHT:
				gravity = port.ImageGravityTopRight
			case v1.ImageGravity_IMAGE_GRAVITY_RIGHT:
				gravity = port.ImageGravityRight
			case v1.ImageGravity_IMAGE_GRAVITY_BOTTOM_RIGHT:
				gravity = port.ImageGravityBottomRight
			case v1.ImageGravity_IMAGE_GRAVITY_BOTTOM:
				gravity = port.ImageGravityBottom
			case v1.ImageGravity_IMAGE_GRAVITY_BOTTOM_LEFT:
				gravity = port.ImageGravityBottomLeft
			case v1.ImageGravity_IMAGE_GRAVITY_LEFT:
				gravity = port.ImageGravityLeft
			case v1.ImageGravity_IMAGE_GRAVITY_TOP_LEFT:
				gravity = port.ImageGravityTopLeft
			default:
				return nil, port.ErrUnsupportedImageGravity
			}

			transform.Resize = &port.ImageResize{
				Width:           req.Transform.Resize.Width,
				Height:          req.Transform.Resize.Height,
				Fit:             fit,
				Gravity:         gravity,
				BackgroundColor: req.Transform.Resize.BackgroundColor,
				AllowUpscale:    req.Transform.Resize.AllowUpscale,
			}
		}

		if req.Transform.Encoding != nil {
			var format port.ImageFormat

			switch req.Transform.Encoding.Format {
			case v1.ImageFormat_IMAGE_FORMAT_UNSPECIFIED:
			case v1.ImageFormat_IMAGE_FORMAT_JPEG:
				format = port.ImageFormatJPEG
			case v1.ImageFormat_IMAGE_FORMAT_PNG:
				format = port.ImageFormatPNG
			case v1.ImageFormat_IMAGE_FORMAT_WEBP:
				format = port.ImageFormatWebP
			case v1.ImageFormat_IMAGE_FORMAT_AVIF:
				format = port.ImageFormatAVIF
			default:
				return nil, port.ErrUnsupportedImageFormat
			}

			transform.Encoding = &port.ImageEncoding{
				Format:  format,
				Quality: req.Transform.Encoding.Quality,
			}
		}
	}

	result, err := s.presignImageHandler.Execute(ctx, &usecase.PresignImage{
		StorageKey: req.StorageKey,
		Transform:  transform,
		ExpireTime: req.ExpireTime.AsTime(),
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

func (s *Service) RemoveImage(ctx context.Context, req *v1.RemoveImageRequest) (*v1.RemoveImageResponse, error) {
	if err := s.removeImageHandler.Execute(ctx, &usecase.RemoveImage{
		StorageKey: req.StorageKey,
	}); err != nil {
		return nil, mapError(err)
	}

	return &v1.RemoveImageResponse{}, nil
}
