package integrationevent

import "time"

type ImageCreated struct {
	BaseIntegrationEvent

	StorageKey string
	CreateTime time.Time
}

func NewImageCreated(
	imageID string,
	storageKey string,
	createTime time.Time,
) ImageCreated {
	return ImageCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID),

		StorageKey: storageKey,
		CreateTime: createTime,
	}
}

func (e ImageCreated) Type() string {
	return "image.created"
}

func (e ImageCreated) AggregateType() string {
	return "image"
}
