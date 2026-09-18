package event

import (
	eventv1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func imageCreatedPayload(event integrationevent.ImageCreated) *eventv1.ImageCreatedPayload {
	return &eventv1.ImageCreatedPayload{
		StorageKey: event.StorageKey,
		CreateTime: timestamppb.New(event.CreateTime),
	}
}
