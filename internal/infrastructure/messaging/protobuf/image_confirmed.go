package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageConfirmedPayload(event integrationevent.ImageConfirmed) *v1.ImageConfirmedPayload {
	return &v1.ImageConfirmedPayload{}
}
