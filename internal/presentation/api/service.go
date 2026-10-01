package api

import (
	"bytes"
	"context"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/application/usecase"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	imageResourcePattern = "images/{image}"
)

type Service struct {
	v1.UnimplementedAssetServiceServer

	getImageHandler     *usecase.GetImageHandler
	uploadImageHandler  *usecase.UploadImageHandler
	presignImageHandler *usecase.PresignImageHandler
	deleteImageHandler  *usecase.DeleteImageHandler
}

func NewService(
	getImageHandler *usecase.GetImageHandler,
	uploadImageHandler *usecase.UploadImageHandler,
	presignImageHandler *usecase.PresignImageHandler,
	deleteImageHandler *usecase.DeleteImageHandler,
) *Service {
	return &Service{
		getImageHandler:     getImageHandler,
		uploadImageHandler:  uploadImageHandler,
		presignImageHandler: presignImageHandler,
		deleteImageHandler:  deleteImageHandler,
	}
}

func (s *Service) GetImage(ctx context.Context, req *v1.GetImageRequest) (*v1.Image, error) {
	var imageID string

	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := s.getImageHandler.Execute(ctx, &usecase.GetImage{ImageID: imageID})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.Image{
		Name:       resourcename.Sprint(imageResourcePattern, result.Image.ID),
		StorageKey: result.Image.StorageKey,
		Ready:      result.Image.Ready,
		CreateTime: timestamppb.New(result.Image.CreateTime),
	}, nil
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
		Image: resourcename.Sprint(imageResourcePattern, result.ImageID),
	}, nil
}

func (s *Service) DeleteImage(ctx context.Context, req *v1.DeleteImageRequest) (*emptypb.Empty, error) {
	var imageID string

	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := s.deleteImageHandler.Execute(ctx, &usecase.DeleteImage{
		ImageID: imageID,
	}); err != nil {
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
