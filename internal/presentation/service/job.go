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

func (s *Job) CreateImage(ctx context.Context, payload *v1.CreateImagePayload) error {
	if _, err := command.Send(ctx, s.commandBus, command.CreateImage{
		Tags:      payload.Tags,
		Metadata:  payload.Metadata,
		ObjectKey: payload.ObjectKey,
		Width:     int(payload.Width),
		Height:    int(payload.Height),
		Format:    parseCreateImageFormatJobV1(payload.Format),
	}); err != nil {
		return err
	}

	return nil
}

func (s *Job) ConfirmImage(ctx context.Context, payload *v1.ConfirmImagePayload) error {
	if _, err := command.Send(ctx, s.commandBus, command.ConfirmImage{
		ImageID: payload.ImageId,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Job) ProcessImage(ctx context.Context, payload *v1.ProcessImagePayload) error {
	if _, err := command.Send(ctx, s.commandBus, command.ProcessImage{
		ImageID: payload.ImageId,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Job) DeleteImage(ctx context.Context, payload *v1.DeleteImagePayload) error {
	if _, err := command.Send(ctx, s.commandBus, command.DeleteImage{
		ImageID: payload.ImageId,
	}); err != nil {
		return err
	}

	return nil
}

func parseCreateImageFormatJobV1(format v1.CreateImagePayload_Format) string {
	return format.String()
}
