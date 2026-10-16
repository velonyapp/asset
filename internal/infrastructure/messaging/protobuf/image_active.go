package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageActivePayload(event integrationevent.ImageActive) *v1.ImageActivePayload {
	return &v1.ImageActivePayload{}
}
