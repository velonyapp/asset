package event

import (
	eventv1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func imageDeletedPayload(event integrationevent.ImageDeleted) *eventv1.ImageDeletedPayload {
	return &eventv1.ImageDeletedPayload{
		DeleteTime: timestamppb.New(event.DeleteTime),
	}
}
