package api

import (
	"context"
	"errors"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/command"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/application/query"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrInvalidImageResourceName = errors.New("invalid image resource name")
)

const (
	imageResourcePattern = "images/{image}"
)

type Service struct {
	v1.UnimplementedAssetServiceServer

	commandBus *command.Bus
	queryBus   *query.Bus
}

func NewService(
	commandBus *command.Bus,
	queryBus *query.Bus,
) *Service {
	return &Service{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

func (s *Service) GetImage(ctx context.Context, req *v1.GetImageRequest) (*v1.Image, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, ErrInvalidImageResourceName
	}

	result, err := query.Send(ctx, s.queryBus, query.GetImage{ImageID: imageID})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:       resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:       result.Image.Tags,
		ObjectKey:  result.Image.ObjectKey,
		State:      mapState(result.Image.State),
		CreateTime: timestamppb.New(result.Image.CreateTime),
		UpdateTime: timestamppb.New(result.Image.UpdateTime),
	}, nil
}

func (s *Service) CreateImage(ctx context.Context, req *v1.CreateImageRequest) (*v1.Image, error) {
	result, err := command.Send(ctx, s.commandBus, command.CreateImage{
		Tags:      req.Image.Tags,
		ObjectKey: req.Image.ObjectKey,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:       resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:       result.Image.Tags,
		ObjectKey:  result.Image.ObjectKey,
		State:      mapState(result.Image.State),
		CreateTime: timestamppb.New(result.Image.CreateTime),
		UpdateTime: timestamppb.New(result.Image.UpdateTime),
	}, nil
}

func (s *Service) InitiateImageUpload(ctx context.Context, req *v1.InitiateImageUploadRequest) (*v1.InitiateImageUploadResponse, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, ErrInvalidImageResourceName
	}

	result, err := command.Send(ctx, s.commandBus, command.InitiateImageUpload{
		ImageID: imageID,
		TTL:     req.Ttl.AsDuration(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.InitiateImageUploadResponse{
		Image: &v1.Image{
			Name:       resourcename.Sprint(imageResourcePattern, result.Image.ID),
			Tags:       result.Image.Tags,
			ObjectKey:  result.Image.ObjectKey,
			State:      mapState(result.Image.State),
			CreateTime: timestamppb.New(result.Image.CreateTime),
			UpdateTime: timestamppb.New(result.Image.UpdateTime),
		},
		UploadUrl: result.UploadURL,
	}, nil
}

func (s *Service) ConfirmImageUpload(ctx context.Context, req *v1.ConfirmImageUploadRequest) (*v1.Image, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, ErrInvalidImageResourceName
	}

	result, err := command.Send(ctx, s.commandBus, command.ConfirmImageUpload{
		ImageID: imageID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:       resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:       result.Image.Tags,
		ObjectKey:  result.Image.ObjectKey,
		State:      mapState(result.Image.State),
		CreateTime: timestamppb.New(result.Image.CreateTime),
		UpdateTime: timestamppb.New(result.Image.UpdateTime),
	}, nil
}

func (s *Service) ProcessImage(ctx context.Context, req *v1.ProcessImageRequest) (*v1.Image, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, ErrInvalidImageResourceName
	}

	var resize *port.ImageResize
	if req.Resize != nil {
		var fit port.ImageResizeFit
		switch req.Resize.Fit {
		case v1.ImageResize_FIT_UNSPECIFIED:
		case v1.ImageResize_CONTAIN:
			fit = port.ImageResizeFitContain
		case v1.ImageResize_COVER:
			fit = port.ImageResizeFitCover
		case v1.ImageResize_PAD:
			fit = port.ImageResizeFitPad
		case v1.ImageResize_STRETCH:
			fit = port.ImageResizeFitStretch
		default:
			return nil, port.ErrUnsupportedResizeFit
		}

		var gravity port.ImageResizeGravity
		switch req.Resize.Gravity {
		case v1.ImageResize_GRAVITY_UNSPECIFIED:
		case v1.ImageResize_CENTER:
			gravity = port.ImageResizeGravityCenter
		case v1.ImageResize_TOP:
			gravity = port.ImageResizeGravityTop
		case v1.ImageResize_TOP_RIGHT:
			gravity = port.ImageResizeGravityTopRight
		case v1.ImageResize_RIGHT:
			gravity = port.ImageResizeGravityRight
		case v1.ImageResize_BOTTOM_RIGHT:
			gravity = port.ImageResizeGravityBottomRight
		case v1.ImageResize_BOTTOM:
			gravity = port.ImageResizeGravityBottom
		case v1.ImageResize_BOTTOM_LEFT:
			gravity = port.ImageResizeGravityBottomLeft
		case v1.ImageResize_LEFT:
			gravity = port.ImageResizeGravityLeft
		case v1.ImageResize_TOP_LEFT:
			gravity = port.ImageResizeGravityTopLeft
		default:
			return nil, port.ErrUnsupportedImageResizeGravity
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
		var format port.ImageEncodingFormat
		switch req.Encoding.Format {
		case v1.ImageEncoding_FORMAT_UNSPECIFIED:
		case v1.ImageEncoding_JPEG:
			format = port.ImageEncodingFormatJPEG
		case v1.ImageEncoding_PNG:
			format = port.ImageEncodingFormatPNG
		case v1.ImageEncoding_WEBP:
			format = port.ImageEncodingFormatWebP
		case v1.ImageEncoding_AVIF:
			format = port.ImageEncodingFormatAVIF
		default:
			return nil, port.ErrUnsupportedImageEncodingFormat
		}

		encoding = &port.ImageEncoding{
			Format:  format,
			Quality: req.Encoding.Quality,
		}
	}

	result, err := command.Send(ctx, s.commandBus, command.ProcessImage{
		ImageID: imageID,
		Options: port.ImageProcessOptions{
			Resize:         resize,
			Encoding:       encoding,
			AutoRotate:     req.AutoRotate,
			RemoveMetadata: req.RemoveMetadata,
		},
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:       resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:       result.Image.Tags,
		ObjectKey:  result.Image.ObjectKey,
		State:      mapState(result.Image.State),
		CreateTime: timestamppb.New(result.Image.CreateTime),
		UpdateTime: timestamppb.New(result.Image.UpdateTime),
	}, nil
}

func (s *Service) DeleteImage(ctx context.Context, req *v1.DeleteImageRequest) (*emptypb.Empty, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, ErrInvalidImageResourceName
	}

	if _, err := command.Send(ctx, s.commandBus, command.DeleteImage{
		ImageID: imageID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func mapState(state string) v1.Image_State {
	switch state {
	case "CREATED":
		return v1.Image_CREATED
	case "UPLOADING":
		return v1.Image_UPLOADING
	case "UPLOADED":
		return v1.Image_UPLOADED
	case "PROCESSING":
		return v1.Image_PROCESSING
	case "PROCESSED":
		return v1.Image_PROCESSED
	default:
		return v1.Image_STATE_UNSPECIFIED
	}
}
