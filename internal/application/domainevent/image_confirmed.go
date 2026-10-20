package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/job"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[event.ImageConfirmed] = (*ImageConfirmedHandler)(nil)

type ImageConfirmedHandler struct {
	jobPublisher port.JobPublisher
}

func NewImageConfirmedHandler(
	jobPublisher port.JobPublisher,
) *ImageConfirmedHandler {
	return &ImageConfirmedHandler{
		jobPublisher: jobPublisher,
	}
}

func (h *ImageConfirmedHandler) Execute(ctx context.Context, domainEvent event.ImageConfirmed) error {
	jobb := job.NewProcessImage(
		domainEvent.AggregateID(),
		domainEvent.ConfirmTime(),
	)

	return h.jobPublisher.Publish(ctx, jobb)
}
