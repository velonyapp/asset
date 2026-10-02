package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[*event.ImageCreated] = (*ImageCreatedHandler)(nil)

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

func (h *ImageCreatedHandler) Execute(ctx context.Context, domainEvent *event.ImageCreated) error {
	tags := domainEvent.Tags()
	tagValues := make([]string, len(tags))
	for i, tag := range tags {
		tagValues[i] = tag.Value()
	}

	integrationEvent := integrationevent.NewImageCreated(
		domainEvent.AggregateID(),
		tagValues,
		domainEvent.ObjectKey().Value(),
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
