package integrationevent

import "time"

type ImageDeleted struct {
	BaseIntegrationEvent

	DeleteTime time.Time
}

func NewImageDeleted(
	imageID string,
	deleteTime time.Time,
) ImageDeleted {
	return ImageDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID),

		DeleteTime: deleteTime,
	}
}

func (e ImageDeleted) Type() string {
	return "image.deleted"
}

func (e ImageDeleted) AggregateType() string {
	return "image"
}
