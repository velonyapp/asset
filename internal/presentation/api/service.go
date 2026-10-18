package api

import (
	"context"
	"errors"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/command"
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
		Name:        resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:        result.Image.Tags,
		ObjectKey:   result.Image.ObjectKey,
		State:       mapImageState(result.Image.State),
		Width:       result.Image.Width,
		Height:      result.Image.Height,
		Format:      mapImageFormat(result.Image.Format),
		CreateTime:  timestamppb.New(result.Image.CreateTime),
		UploadUrl:   result.UploadURL,
		DownloadUrl: result.DownloadURL,
	}, nil
}

func (s *Service) CreateImage(ctx context.Context, req *v1.CreateImageRequest) (*v1.Image, error) {
	result, err := command.Send(ctx, s.commandBus, command.CreateImage{
		Tags:      req.Image.Tags,
		ObjectKey: req.Image.ObjectKey,
		Width:     int(req.Image.Width),
		Height:    int(req.Image.Height),
		Format:    req.Image.Format.String(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:       resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:       result.Image.Tags,
		ObjectKey:  result.Image.ObjectKey,
		State:      mapImageState(result.Image.State),
		Width:      result.Image.Width,
		Height:     result.Image.Height,
		Format:     mapImageFormat(result.Image.Format),
		CreateTime: timestamppb.New(result.Image.CreateTime),
		UploadUrl:  &result.UploadURL,
	}, nil
}

func (s *Service) ConfirmImage(ctx context.Context, req *v1.ConfirmImageRequest) (*v1.Image, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, ErrInvalidImageResourceName
	}

	result, err := command.Send(ctx, s.commandBus, command.ConfirmImage{
		ImageID: imageID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:       resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:       result.Image.Tags,
		ObjectKey:  result.Image.ObjectKey,
		State:      mapImageState(result.Image.State),
		Width:      result.Image.Width,
		Height:     result.Image.Height,
		Format:     mapImageFormat(result.Image.Format),
		CreateTime: timestamppb.New(result.Image.CreateTime),
	}, nil
}

func (s *Service) ProcessImage(ctx context.Context, req *v1.ProcessImageRequest) (*v1.Image, error) {
	var imageID string
	if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
		return nil, ErrInvalidImageResourceName
	}

	result, err := command.Send(ctx, s.commandBus, command.ProcessImage{
		ImageID: imageID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:        resourcename.Sprint(imageResourcePattern, result.Image.ID),
		Tags:        result.Image.Tags,
		ObjectKey:   result.Image.ObjectKey,
		State:       mapImageState(result.Image.State),
		Width:       result.Image.Width,
		Height:      result.Image.Height,
		Format:      mapImageFormat(result.Image.Format),
		CreateTime:  timestamppb.New(result.Image.CreateTime),
		DownloadUrl: &result.DownloadURL,
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

func mapImageState(state string) v1.Image_State {
	switch state {
	case "AWAITING_UPLOAD":
		return v1.Image_AWAITING_UPLOAD
	case "PENDING":
		return v1.Image_PENDING
	case "PROCESSING":
		return v1.Image_PROCESSING
	case "ACTIVE":
		return v1.Image_ACTIVE
	default:
		return v1.Image_STATE_UNSPECIFIED
	}
}

func mapImageFormat(state string) v1.Image_Format {
	switch state {
	case "JPEG":
		return v1.Image_JPEG
	case "PNG":
		return v1.Image_PNG
	case "WEBP":
		return v1.Image_WEBP
	case "AVIF":
		return v1.Image_AVIF
	default:
		return v1.Image_FORMAT_UNSPECIFIED
	}
}
