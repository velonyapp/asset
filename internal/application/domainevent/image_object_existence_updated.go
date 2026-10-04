package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[event.ImageObjectExistenceUpdated] = (*ImageObjectExistenceUpdatedHandler)(nil)

type ImageObjectExistenceUpdatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageObjectExistenceUpdatedHandler(
	eventPublisher port.EventPublisher,
) *ImageObjectExistenceUpdatedHandler {
	return &ImageObjectExistenceUpdatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageObjectExistenceUpdatedHandler) Execute(ctx context.Context, domainEvent event.ImageObjectExistenceUpdated) error {
	integrationEvent := integrationevent.NewImageObjectExistenceUpdated(
		domainEvent.AggregateID(),
		domainEvent.Tags().Strings(),
		domainEvent.UpdateTime(),
		domainEvent.ObjectExists(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
