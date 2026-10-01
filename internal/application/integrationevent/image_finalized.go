package integrationevent

import "time"

var _ IntegrationEvent = (*ImageFinalized)(nil)

type ImageFinalized struct {
	BaseIntegrationEvent

	tags []string
}

func NewImageFinalized(
	imageID string,
	tags []string,
	occurTime time.Time,
) *ImageFinalized {
	return &ImageFinalized{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, occurTime),

		tags: tags,
	}
}

func (e *ImageFinalized) Type() string {
	return "asset.image.finalized"
}

func (e *ImageFinalized) AggregateType() string {
	return "image"
}

func (i *ImageFinalized) Tags() []string {
	return append([]string(nil), i.tags...)
}
