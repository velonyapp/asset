package domainevent

import (
	"context"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/event"
)

var _ Handler[event.ImageReconciled] = (*ImageReconciledHandler)(nil)

type ImageReconciledHandler struct {
	eventPublisher port.EventPublisher
}

func NewImageReconciledHandler(
	eventPublisher port.EventPublisher,
) *ImageReconciledHandler {
	return &ImageReconciledHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageReconciledHandler) Execute(ctx context.Context, domainEvent event.ImageReconciled) error {
	integrationEvent := integrationevent.NewImageReconciled(
		domainEvent.AggregateID(),
		domainEvent.Tags().Strings(),
		domainEvent.UpdateTime(),
		domainEvent.ObjectExists(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
