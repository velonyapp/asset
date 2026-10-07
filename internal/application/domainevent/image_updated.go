package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[event.ImageUpdated] = (*ImageUpdatedHandler)(nil)

type ImageUpdatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageUpdatedHandler(
	eventPublisher port.EventPublisher,
) *ImageUpdatedHandler {
	return &ImageUpdatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageUpdatedHandler) Execute(ctx context.Context, domainEvent event.ImageUpdated) error {
	integrationEvent := integrationevent.NewImageUpdated(
		domainEvent.AggregateID(),
		domainEvent.Tags().Strings(),
		domainEvent.UpdateTime(),
		domainEvent.State().String(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
