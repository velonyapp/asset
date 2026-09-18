package event

import (
	"fmt"

	eventv1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(event integrationevent.IntegrationEvent) (*eventv1.Event, error) {
	if event == nil {
		return nil, fmt.Errorf("integration event is nil")
	}

	payload, err := e.payload(event)
	if err != nil {
		return nil, err
	}

	payloadAny, err := anypb.New(payload)
	if err != nil {
		return nil, err
	}

	return &eventv1.Event{
		Id:            event.ID(),
		Type:          event.Type(),
		AggregateId:   event.AggregateID(),
		AggregateType: event.AggregateType(),
		OccurTime:     timestamppb.New(event.OccurTime()),
		Payload:       payloadAny,
	}, nil
}

func (e *Encoder) payload(event integrationevent.IntegrationEvent) (proto.Message, error) {
	switch event := event.(type) {
	case integrationevent.ImageCreated:
		return imageCreatedPayload(event), nil

	case *integrationevent.ImageCreated:
		if event == nil {
			return nil, fmt.Errorf("ImageCreated event is nil")
		}

		return imageCreatedPayload(*event), nil

	case integrationevent.ImageDeleted:
		return imageDeletedPayload(event), nil

	case *integrationevent.ImageDeleted:
		if event == nil {
			return nil, fmt.Errorf("ImageDeleted event is nil")
		}

		return imageDeletedPayload(*event), nil

	default:
		return nil, fmt.Errorf("unknown integration event %T", event)
	}
}
