package integrationevent

import "time"

var _ IntegrationEvent = (*ImageProcessed)(nil)

type ImageProcessed struct {
	BaseIntegrationEvent
}

func NewImageProcessed(
	imageID string,
	tags []string,
	occurTime time.Time,
) *ImageProcessed {
	return &ImageProcessed{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),
	}
}

func (e *ImageProcessed) Type() string {
	return "asset.image.processed"
}

func (e *ImageProcessed) AggregateType() string {
	return "image"
}
