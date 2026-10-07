package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageUpdatedPayload(event integrationevent.ImageUpdated) *v1.ImageUpdatedPayload {
	return &v1.ImageUpdatedPayload{
		State: event.State(),
	}
}
