package integrationevent

import "time"

var _ IntegrationEvent = (*ImageDeleted)(nil)

type ImageDeleted struct {
	BaseIntegrationEvent

	tags []string
}

func NewImageDeleted(
	imageID string,
	tags []string,
	occurTime time.Time,
) *ImageDeleted {
	return &ImageDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, occurTime),

		tags: tags,
	}
}

func (e *ImageDeleted) Type() string {
	return "asset.image.deleted"
}

func (e *ImageDeleted) AggregateType() string {
	return "image"
}

func (i *ImageDeleted) Tags() []string {
	return append([]string(nil), i.tags...)
}
