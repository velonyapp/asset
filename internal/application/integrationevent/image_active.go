package integrationevent

import "time"

var _ IntegrationEvent = (*ImageActive)(nil)

type ImageActive struct {
	BaseIntegrationEvent
}

func NewImageActive(
	imageID string,
	tags []string,
	occurTime time.Time,
) ImageActive {
	return ImageActive{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),
	}
}

func (e ImageActive) Type() string {
	return "asset.image.active"
}

func (e ImageActive) AggregateType() string {
	return "image"
}
