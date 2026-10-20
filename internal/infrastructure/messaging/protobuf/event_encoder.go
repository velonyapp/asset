package protobuf

import (
	"fmt"

	"github.com/velonyapp/asset/internal/application/integrationevent"

	"google.golang.org/protobuf/proto"
)

type EventEncoder struct{}

func NewEventEncoder() *EventEncoder {
	return &EventEncoder{}
}

func (e *EventEncoder) Encode(eventPayload integrationevent.IntegrationEvent) (proto.Message, error) {
	switch eventPayload := eventPayload.(type) {
	case integrationevent.ImageCreated:
		return imageCreatedPayload(eventPayload), nil
	case integrationevent.ImageActive:
		return imageActivePayload(eventPayload), nil
	case integrationevent.ImageDeleted:
		return imageDeletedPayload(eventPayload), nil
	default:
		return nil, fmt.Errorf("unknown integration event %T", eventPayload)
	}
}
