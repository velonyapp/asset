package event

import (
	eventv1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageFinalizedPayload(event integrationevent.ImageFinalized) *eventv1.ImageFinalizedPayload {
	return &eventv1.ImageFinalizedPayload{}
}
