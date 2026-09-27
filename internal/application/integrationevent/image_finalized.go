package integrationevent

import "time"

var _ IntegrationEvent = (*ImageFinalized)(nil)

type ImageFinalized struct {
	BaseIntegrationEvent
}

func NewImageFinalized(
	imageID string,
	occurTime time.Time,
) *ImageFinalized {
	return &ImageFinalized{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, occurTime),
	}
}

func (e *ImageFinalized) Type() string {
	return "image.finalized"
}

func (e *ImageFinalized) AggregateType() string {
	return "image"
}
