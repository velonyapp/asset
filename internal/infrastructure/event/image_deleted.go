package event

import (
	eventv1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageDeletedPayload(event integrationevent.ImageDeleted) *eventv1.ImageDeletedPayload {
	return &eventv1.ImageDeletedPayload{
		Tags: event.Tags(),
	}
}
