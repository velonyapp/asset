package event

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageCreatedPayload(event integrationevent.ImageCreated) *v1.ImageCreatedPayload {
	return &v1.ImageCreatedPayload{
		ObjectKey: event.ObjectKey(),
	}
}
