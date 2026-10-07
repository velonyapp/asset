package integrationevent

import "time"

var _ IntegrationEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseIntegrationEvent

	objectKey string
	state     string
}

func NewImageCreated(
	imageID string,
	tags []string,
	occurTime time.Time,
	objectKey string,
	state string,
) ImageCreated {
	return ImageCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),

		objectKey: objectKey,
		state:     state,
	}
}

func (e ImageCreated) Type() string {
	return "asset.image.created"
}

func (e ImageCreated) AggregateType() string {
	return "image"
}

func (e ImageCreated) ObjectKey() string {
	return e.objectKey
}

func (e ImageCreated) State() string {
	return e.state
}
