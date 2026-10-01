package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[*event.ImageFinalized] = (*ImageFinalizedHandler)(nil)

type ImageFinalizedHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageFinalizedHandler(
	eventPublisher port.EventPublisher,
) *ImageFinalizedHandler {
	return &ImageFinalizedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageFinalizedHandler) Execute(ctx context.Context, domainEvent *event.ImageFinalized) error {
	tags := domainEvent.Tags()
	tagValues := make([]string, len(tags))
	for i, tag := range tags {
		tagValues[i] = tag.Value()
	}

	integrationEvent := integrationevent.NewImageFinalized(
		domainEvent.AggregateID(),
		tagValues,
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
