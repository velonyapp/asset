package integrationevent

import "time"

var _ IntegrationEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseIntegrationEvent

	storageKey string
}

func NewImageCreated(
	imageID string,
	storageKey string,
	occurTime time.Time,
) *ImageCreated {
	return &ImageCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, occurTime),

		storageKey: storageKey,
	}
}

func (e *ImageCreated) Type() string {
	return "image.created"
}

func (e *ImageCreated) AggregateType() string {
	return "image"
}

func (e *ImageCreated) StorageKey() string {
	return e.storageKey
}
