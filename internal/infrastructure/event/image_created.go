package event

import (
	eventv1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageCreatedPayload(event integrationevent.ImageCreated) *eventv1.ImageCreatedPayload {
	return &eventv1.ImageCreatedPayload{
		StorageKey: event.StorageKey(),
	}
}
