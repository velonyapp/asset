package integrationevent

import "time"

var _ IntegrationEvent = (*ImageDeleted)(nil)

type ImageDeleted struct {
	BaseIntegrationEvent
}

func NewImageDeleted(
	imageID string,
	tags []string,
	occurTime time.Time,
) ImageDeleted {
	return ImageDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),
	}
}

func (e ImageDeleted) Type() string {
	return "asset.image.deleted"
}

func (e ImageDeleted) AggregateType() string {
	return "image"
}
