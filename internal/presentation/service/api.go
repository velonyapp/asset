package service

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

type API struct {
	v1.UnimplementedAssetServiceServer

	commandBus *command.Bus
	queryBus   *query.Bus
}

func NewAPI(
	commandBus *command.Bus,
	queryBus *query.Bus,
) *API {
	return &API{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

func (s *API) GetImage(ctx context.Context, req *v1.GetImageRequest) (*v1.Image, error) {
	imageID, err := parseImageName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := query.Send(ctx, s.queryBus, query.GetImage{ImageID: imageID})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:        formatImageName(result.Image.ID),
		Tags:        result.Image.Tags,
		ObjectKey:   result.Image.ObjectKey,
		State:       formatImageStateAPIV1(result.Image.State),
		Width:       result.Image.Width,
		Height:      result.Image.Height,
		Format:      formatImageFormatAPIV1(result.Image.Format),
		CreateTime:  timestamppb.New(result.Image.CreateTime),
		UploadUrl:   result.UploadURL,
		DownloadUrl: result.DownloadURL,
	}, nil
}

func (s *API) CreateImage(ctx context.Context, req *v1.CreateImageRequest) (*v1.Image, error) {
	result, err := command.Send(ctx, s.commandBus, command.CreateImage{
		Tags:      req.Image.Tags,
		ObjectKey: req.Image.ObjectKey,
		Width:     int(req.Image.Width),
		Height:    int(req.Image.Height),
		Format:    parseImageFormatAPIV1(req.Image.Format),
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:       formatImageName(result.Image.ID),
		Tags:       result.Image.Tags,
		ObjectKey:  result.Image.ObjectKey,
		State:      formatImageStateAPIV1(result.Image.State),
		Width:      result.Image.Width,
		Height:     result.Image.Height,
		Format:     formatImageFormatAPIV1(result.Image.Format),
		CreateTime: timestamppb.New(result.Image.CreateTime),
		UploadUrl:  &result.UploadURL,
	}, nil
}

func (s *API) ConfirmImage(ctx context.Context, req *v1.ConfirmImageRequest) (*v1.Image, error) {
	imageID, err := parseImageName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.ConfirmImage{
		ImageID: imageID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:       formatImageName(result.Image.ID),
		Tags:       result.Image.Tags,
		ObjectKey:  result.Image.ObjectKey,
		State:      formatImageStateAPIV1(result.Image.State),
		Width:      result.Image.Width,
		Height:     result.Image.Height,
		Format:     formatImageFormatAPIV1(result.Image.Format),
		CreateTime: timestamppb.New(result.Image.CreateTime),
	}, nil
}

func (s *API) ProcessImage(ctx context.Context, req *v1.ProcessImageRequest) (*v1.Image, error) {
	imageID, err := parseImageName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.ProcessImage{
		ImageID: imageID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Image{
		Name:        formatImageName(result.Image.ID),
		Tags:        result.Image.Tags,
		ObjectKey:   result.Image.ObjectKey,
		State:       formatImageStateAPIV1(result.Image.State),
		Width:       result.Image.Width,
		Height:      result.Image.Height,
		Format:      formatImageFormatAPIV1(result.Image.Format),
		CreateTime:  timestamppb.New(result.Image.CreateTime),
		DownloadUrl: &result.DownloadURL,
	}, nil
}

func (s *API) DeleteImage(ctx context.Context, req *v1.DeleteImageRequest) (*emptypb.Empty, error) {
	imageID, err := parseImageName(req.Name)
	if err != nil {
		return nil, err
	}

	if _, err := command.Send(ctx, s.commandBus, command.DeleteImage{
		ImageID: imageID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func formatImageStateAPIV1(state string) v1.Image_State {
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

func parseImageFormatAPIV1(format v1.Image_Format) string {
	return format.String()
}

func formatImageFormatAPIV1(format string) v1.Image_Format {
	switch format {
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

func parseImageName(name string) (string, error) {
	var imageID string
	if err := resourcename.Sscan(name, imageResourcePattern, &imageID); err != nil {
		return "", ErrInvalidImageResourceName
	}

	return imageID, nil
}

func formatImageName(imageID string) string {
	return resourcename.Sprint(imageResourcePattern, imageID)
}
