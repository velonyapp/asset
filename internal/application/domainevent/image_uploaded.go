package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[event.ImageUploaded] = (*ImageUploadedHandler)(nil)

type ImageUploadedHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageUploadedHandler(
	eventPublisher port.EventPublisher,
) *ImageUploadedHandler {
	return &ImageUploadedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageUploadedHandler) Execute(ctx context.Context, domainEvent event.ImageUploaded) error {
	integrationEvent := integrationevent.NewImageUploaded(
		domainEvent.AggregateID(),
		domainEvent.Tags().Strings(),
		domainEvent.UploadTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
