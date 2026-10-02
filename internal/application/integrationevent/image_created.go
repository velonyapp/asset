package integrationevent

import "time"

var _ IntegrationEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseIntegrationEvent

	objectKey string
}

func NewImageCreated(
	imageID string,
	tags []string,
	objectKey string,
	occurTime time.Time,
) *ImageCreated {
	return &ImageCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, occurTime, tags),

		objectKey: objectKey,
	}
}

func (e *ImageCreated) Type() string {
	return "asset.image.created"
}

func (e *ImageCreated) AggregateType() string {
	return "image"
}

func (e *ImageCreated) ObjectKey() string {
	return e.objectKey
}
