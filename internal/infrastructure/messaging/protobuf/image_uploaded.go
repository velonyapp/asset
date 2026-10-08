package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/asset/internal/application/integrationevent"
)

func imageUploadedPayload(event integrationevent.ImageUploaded) *v1.ImageUploadedPayload {
	return &v1.ImageUploadedPayload{}
}
