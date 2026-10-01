package integrationevent

import "time"

var _ IntegrationEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseIntegrationEvent

	tags       []string
	storageKey string
}

func NewImageCreated(
	imageID string,
	tags []string,
	storageKey string,
	occurTime time.Time,
) *ImageCreated {
	return &ImageCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, occurTime),

		tags:       tags,
		storageKey: storageKey,
	}
}

func (e *ImageCreated) Type() string {
	return "asset.image.created"
}

func (e *ImageCreated) AggregateType() string {
	return "image"
}

func (i *ImageCreated) Tags() []string {
	return append([]string(nil), i.tags...)
}

func (e *ImageCreated) StorageKey() string {
	return e.storageKey
}
