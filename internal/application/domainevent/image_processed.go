package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[event.ImageProcessed] = (*ImageProcessedHandler)(nil)

type ImageProcessedHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageProcessedHandler(
	eventPublisher port.EventPublisher,
) *ImageProcessedHandler {
	return &ImageProcessedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageProcessedHandler) Execute(ctx context.Context, domainEvent event.ImageProcessed) error {
	integrationEvent := integrationevent.NewImageProcessed(
		domainEvent.AggregateID(),
		domainEvent.Tags().Strings(),
		domainEvent.ProcessTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
