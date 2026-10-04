package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageProcessedPayload(event integrationevent.ImageProcessed) *v1.ImageProcessedPayload {
	return &v1.ImageProcessedPayload{}
}
