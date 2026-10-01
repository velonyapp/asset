package event

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageFinalizedPayload(event integrationevent.ImageFinalized) *v1.ImageFinalizedPayload {
	return &v1.ImageFinalizedPayload{
		Tags: event.Tags(),
	}
}
