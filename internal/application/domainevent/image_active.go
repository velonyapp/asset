package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[event.ImageActive] = (*ImageActiveHandler)(nil)

type ImageActiveHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageActiveHandler(
	eventPublisher port.EventPublisher,
) *ImageActiveHandler {
	return &ImageActiveHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageActiveHandler) Execute(ctx context.Context, domainEvent event.ImageActive) error {
	integrationEvent := integrationevent.NewImageActive(
		domainEvent.AggregateID(),
		domainEvent.Tags().Strings(),
		domainEvent.ActiveTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
