package job

import "time"

var _ Job = (*ProcessImage)(nil)

type ProcessImage struct {
	BaseJob

	imageID string
}

func NewProcessImage(
	imageID string,
	enqueueTime time.Time,
) ProcessImage {
	return ProcessImage{
		BaseJob: NewBaseJob(enqueueTime),

		imageID: imageID,
	}
}

func (e ProcessImage) Type() string {
	return "asset.image.process"
}

func (e ProcessImage) ImageID() string {
	return e.imageID
}
