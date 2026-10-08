package integrationevent

import "time"

var _ IntegrationEvent = (*ImageUploaded)(nil)

type ImageUploaded struct {
	BaseIntegrationEvent
}

func NewImageUploaded(
	imageID string,
	tags []string,
	occurTime time.Time,
) ImageUploaded {
	return ImageUploaded{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),
	}
}

func (e ImageUploaded) Type() string {
	return "asset.image.uploaded"
}

func (e ImageUploaded) AggregateType() string {
	return "image"
}
