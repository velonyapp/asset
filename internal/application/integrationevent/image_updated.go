package integrationevent

import "time"

var _ IntegrationEvent = (*ImageUpdated)(nil)

type ImageUpdated struct {
	BaseIntegrationEvent

	state string
}

func NewImageUpdated(
	imageID string,
	tags []string,
	occurTime time.Time,
	state string,
) ImageUpdated {
	return ImageUpdated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),

		state: state,
	}
}

func (e ImageUpdated) Type() string {
	return "asset.image.updated"
}

func (e ImageUpdated) AggregateType() string {
	return "image"
}

func (e ImageUpdated) State() string {
	return e.state
}
