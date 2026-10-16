package integrationevent

import "time"

var _ IntegrationEvent = (*ImageConfirmed)(nil)

type ImageConfirmed struct {
	BaseIntegrationEvent
}

func NewImageConfirmed(
	imageID string,
	tags []string,
	occurTime time.Time,
) ImageConfirmed {
	return ImageConfirmed{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),
	}
}

func (e ImageConfirmed) Type() string {
	return "asset.image.confirmed"
}

func (e ImageConfirmed) AggregateType() string {
	return "image"
}
