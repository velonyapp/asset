package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[event.ImageConfirmed] = (*ImageConfirmedHandler)(nil)

type ImageConfirmedHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageConfirmedHandler(
	eventPublisher port.EventPublisher,
) *ImageConfirmedHandler {
	return &ImageConfirmedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageConfirmedHandler) Execute(ctx context.Context, domainEvent event.ImageConfirmed) error {
	integrationEvent := integrationevent.NewImageConfirmed(
		domainEvent.AggregateID(),
		domainEvent.Tags().Strings(),
		domainEvent.ConfirmTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
