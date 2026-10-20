package protobuf

import (
	"fmt"

	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventEncoder struct{}

func NewEventEncoder() *EventEncoder {
	return &EventEncoder{}
}

func (e *EventEncoder) Encode(event integrationevent.IntegrationEvent) (*v1.Event, error) {
	payload, err := e.eventPayload(event)
	if err != nil {
		return nil, err
	}

	payloadAny, err := anypb.New(payload)
	if err != nil {
		return nil, err
	}

	return &v1.Event{
		Id:            event.ID(),
		Type:          event.Type(),
		AggregateId:   event.AggregateID(),
		AggregateType: event.AggregateType(),
		Tags:          event.Tags(),
		Payload:       payloadAny,
		OccurTime:     timestamppb.New(event.OccurTime()),
	}, nil
}

func (e *EventEncoder) eventPayload(event integrationevent.IntegrationEvent) (proto.Message, error) {
	switch event := event.(type) {
	case integrationevent.ImageCreated:
		return imageCreatedPayload(event), nil
	case integrationevent.ImageActive:
		return imageActivePayload(event), nil
	case integrationevent.ImageDeleted:
		return imageDeletedPayload(event), nil
	default:
		return nil, fmt.Errorf("unknown integration event %T", event)
	}
}
