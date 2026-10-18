package service

import (
	"context"

	v1 "github.com/velonyapp/asset/gen/job/v1"
	"github.com/velonyapp/asset/internal/application/command"
)

type Job struct {
	commandBus *command.Bus
}

func NewJob(
	commandBus *command.Bus,
) *Job {
	return &Job{
		commandBus: commandBus,
	}
}

func (s *Job) CreateImage(ctx context.Context, req *v1.CreateImage) error {
	if _, err := command.Send(ctx, s.commandBus, command.CreateImage{
		Tags:      req.Tags,
		ObjectKey: req.ObjectKey,
		Width:     int(req.Width),
		Height:    int(req.Height),
		Format:    parseCreateImageFormatJobV1(req.Format),
	}); err != nil {
		return err
	}

	return nil
}

func (s *Job) ConfirmImage(ctx context.Context, req *v1.ConfirmImage) error {
	imageID, err := parseImageName(req.Name)
	if err != nil {
		return err
	}

	if _, err := command.Send(ctx, s.commandBus, command.ConfirmImage{
		ImageID: imageID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Job) ProcessImage(ctx context.Context, req *v1.ProcessImage) error {
	imageID, err := parseImageName(req.Name)
	if err != nil {
		return err
	}

	if _, err := command.Send(ctx, s.commandBus, command.ProcessImage{
		ImageID: imageID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Job) DeleteImage(ctx context.Context, req *v1.DeleteImage) error {
	imageID, err := parseImageName(req.Name)
	if err != nil {
		return err
	}

	if _, err := command.Send(ctx, s.commandBus, command.DeleteImage{
		ImageID: imageID,
	}); err != nil {
		return err
	}

	return nil
}

func parseCreateImageFormatJobV1(format v1.CreateImage_Format) string {
	return format.String()
}
