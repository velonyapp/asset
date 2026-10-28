package integrationevent

import (
	"maps"
	"time"
)

var _ IntegrationEvent = (*ImageDeleted)(nil)

type ImageDeleted struct {
	BaseIntegrationEvent

	tags     []string
	metadata map[string]string
}

func NewImageDeleted(
	aggregateID string,
	occurTime time.Time,
	tags []string,
	metadata map[string]string,
) ImageDeleted {
	return ImageDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		tags:     append([]string(nil), tags...),
		metadata: maps.Clone(metadata),
	}
}

func (e ImageDeleted) Type() string {
	return "asset.image.deleted"
}

func (e ImageDeleted) AggregateType() string {
	return "image"
}

func (e ImageDeleted) Tags() []string {
	return append([]string(nil), e.tags...)
}

func (e ImageDeleted) Metadata() map[string]string {
	return maps.Clone(e.metadata)
}
