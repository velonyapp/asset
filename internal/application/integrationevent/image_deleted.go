package integrationevent

import "time"

var _ IntegrationEvent = (*ImageDeleted)(nil)

type ImageDeleted struct {
	BaseIntegrationEvent
}

func NewImageDeleted(
	imageID string,
	occurTime time.Time,
) *ImageDeleted {
	return &ImageDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, occurTime),
	}
}

func (e *ImageDeleted) Type() string {
	return "image.deleted"
}

func (e *ImageDeleted) AggregateType() string {
	return "image"
}
