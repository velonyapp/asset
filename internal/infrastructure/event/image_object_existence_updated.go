package event

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageObjectExistenceUpdatedPayload(event integrationevent.ImageObjectExistenceUpdated) *v1.ImageObjectExistenceUpdatedPayload {
	return &v1.ImageObjectExistenceUpdatedPayload{
		ObjectExists: event.ObjectExists(),
	}
}
