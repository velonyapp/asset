package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/job/v1"
	"github.com/velonyapp/asset/internal/application/job"
)

func processImagePayload(job job.ProcessImage) *v1.ProcessImagePayload {
	return &v1.ProcessImagePayload{
		ImageId: job.ImageID(),
	}
}
