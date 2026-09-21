package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

type ImageCreatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageCreatedHandler(
	eventPublisher port.EventPublisher,
) *ImageCreatedHandler {
	return &ImageCreatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageCreatedHandler) Execute(ctx context.Context, domainEvent event.ImageCreated) error {
	integrationEvent := integrationevent.NewImageCreated(
		domainEvent.AggregateID(),
		domainEvent.StorageKey.Value(),
		domainEvent.CreateTime.Value(),
	)

	return h.eventPublisher.Publish(ctx, &integrationEvent)
}
