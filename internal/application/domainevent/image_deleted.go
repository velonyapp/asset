package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

type ImageDeletedHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageDeletedHandler(
	eventPublisher port.EventPublisher,
) *ImageDeletedHandler {
	return &ImageDeletedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageDeletedHandler) Execute(ctx context.Context, domainEvent event.ImageDeleted) error {
	integrationEvent := integrationevent.NewImageDeleted(
		domainEvent.AggregateID(),
		domainEvent.DeleteTime.Value(),
	)

	return h.eventPublisher.Publish(ctx, &integrationEvent)
}
