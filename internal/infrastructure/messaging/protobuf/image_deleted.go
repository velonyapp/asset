package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageDeletedPayload(event integrationevent.ImageDeleted) *v1.ImageDeletedPayload {
	return &v1.ImageDeletedPayload{}
}
