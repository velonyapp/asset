package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

type ImageCreatedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewImageCreatedHandler(
	outboxPublisher port.OutboxPublisher,
) *ImageCreatedHandler {
	return &ImageCreatedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *ImageCreatedHandler) Execute(ctx context.Context, domainEvent event.ImageCreated) error {
	integrationEvent := integrationevent.NewImageCreated(
		domainEvent.AggregateID(),
		domainEvent.StorageKey.Value(),
		domainEvent.CreateTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx,
		port.OutboxMessage{
			PartitionKey: domainEvent.AggregateID(),
			Event:        &integrationEvent,
		},
	)
}
