package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageReconciledPayload(event integrationevent.ImageReconciled) *v1.ImageReconciledPayload {
	return &v1.ImageReconciledPayload{
		ObjectExists: event.ObjectExists(),
	}
}
