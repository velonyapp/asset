package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

type ImageDeletedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewImageDeletedHandler(
	outboxPublisher port.OutboxPublisher,
) *ImageDeletedHandler {
	return &ImageDeletedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *ImageDeletedHandler) Execute(ctx context.Context, domainEvent event.ImageDeleted) error {
	integrationEvent := integrationevent.NewImageDeleted(
		domainEvent.AggregateID(),
		domainEvent.DeleteTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx,
		port.OutboxMessage{
			PartitionKey: domainEvent.AggregateID(),
			Event:        integrationEvent,
		},
	)
}
