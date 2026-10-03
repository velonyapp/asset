package integrationevent

import "time"

var _ IntegrationEvent = (*ImageObjectExistenceUpdated)(nil)

type ImageObjectExistenceUpdated struct {
	BaseIntegrationEvent

	objectExists bool
}

func NewImageObjectExistenceUpdated(
	imageID string,
	tags []string,
	occurTime time.Time,
	objectExists bool,
) *ImageObjectExistenceUpdated {
	return &ImageObjectExistenceUpdated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),

		objectExists: objectExists,
	}
}

func (e *ImageObjectExistenceUpdated) Type() string {
	return "asset.image.object-existence.updated"
}

func (e *ImageObjectExistenceUpdated) AggregateType() string {
	return "image"
}

func (e *ImageObjectExistenceUpdated) ObjectExists() bool {
	return e.objectExists
}
