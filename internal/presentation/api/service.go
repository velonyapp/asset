package api

import (
	"context"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/command"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/application/query"

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
}

func NewService(
	_ *command.HandlerRegistry,
	_ *query.HandlerRegistry,
) *Service {
	return &Service{}
}

func (s *Service) GetImage(ctx context.Context, req *v1.GetImageRequest) (*v1.Image, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := query.Send(ctx, &query.GetImage{ImageID: imageID})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.Image{
		Name:         resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:         result.Image.Tags,
		ObjectKey:    result.Image.ObjectKey,
		ObjectExists: result.Image.ObjectExists,
		CreateTime:   timestamppb.New(result.Image.CreateTime),
	}, nil
}

func (s *Service) CreateImage(ctx context.Context, req *v1.CreateImageRequest) (*v1.Image, error) {
	result, err := command.Send(ctx, &command.CreateImage{
		Tags:      req.Image.Tags,
		ObjectKey: req.Image.ObjectKey,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.Image{
		Name:         resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:         result.Image.Tags,
		ObjectKey:    result.Image.ObjectKey,
		ObjectExists: result.Image.ObjectExists,
		CreateTime:   timestamppb.New(result.Image.CreateTime),
	}, nil
}

func (s *Service) PresignImage(ctx context.Context, req *v1.PresignImageRequest) (*v1.PresignImageResponse, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := query.Send(ctx, &query.PresignImage{
		ImageID: imageID,
		TTL:     req.Ttl.AsDuration(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.PresignImageResponse{
		UploadUrl: result.UploadURL,
	}, nil
}

func (s *Service) ProcessImage(ctx context.Context, req *v1.ProcessImageRequest) (*v1.ProcessImageResponse, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var resize *port.ImageResize
	if req.Resize != nil {
		var fit port.ImageResizeFit
		switch req.Resize.Fit {
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
		switch req.Resize.Gravity {
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

		resize = &port.ImageResize{
			Width:           req.Resize.Width,
			Height:          req.Resize.Height,
			Fit:             fit,
			Gravity:         gravity,
			BackgroundColor: req.Resize.BackgroundColor,
			AllowUpscale:    req.Resize.AllowUpscale,
		}
	}

	var encoding *port.ImageEncoding
	if req.Encoding != nil {
		var format port.ImageFormat
		switch req.Encoding.Format {
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

		encoding = &port.ImageEncoding{
			Format:  format,
			Quality: req.Encoding.Quality,
		}
	}

	if _, err := command.Send(ctx, &command.ProcessImage{
		ImageID: imageID,
		Options: port.ImageProcessOptions{
			Resize:         resize,
			Encoding:       encoding,
			AutoRotate:     req.AutoRotate,
			RemoveMetadata: req.RemoveMetadata,
		},
	}); err != nil {
		return nil, mapError(err)
	}

	return &v1.ProcessImageResponse{}, nil
}

func (s *Service) ReconcileImage(ctx context.Context, req *v1.ReconcileImageRequest) (*v1.ReconcileImageResponse, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := command.Send(ctx, &command.ReconcileImage{
		ImageID: imageID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.ReconcileImageResponse{
		Image: &v1.Image{
			Name:         resourcename.Sprint(imageResourcePattern, result.Image.ID),
			Tags:         result.Image.Tags,
			ObjectKey:    result.Image.ObjectKey,
			ObjectExists: result.Image.ObjectExists,
			CreateTime:   timestamppb.New(result.Image.CreateTime),
		},
	}, nil
}

func (s *Service) DeleteImage(ctx context.Context, req *v1.DeleteImageRequest) (*emptypb.Empty, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := command.Send(ctx, &command.DeleteImage{
		ImageID: imageID,
	}); err != nil {
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
