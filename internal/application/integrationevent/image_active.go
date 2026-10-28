package integrationevent

import (
	"maps"
	"time"
)

var _ IntegrationEvent = (*ImageActive)(nil)

type ImageActive struct {
	BaseIntegrationEvent

	tags     []string
	metadata map[string]string
}

func NewImageActive(
	aggregateID string,
	occurTime time.Time,
	tags []string,
	metadata map[string]string,
) ImageActive {
	return ImageActive{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		tags:     append([]string(nil), tags...),
		metadata: maps.Clone(metadata),
	}
}

func (e ImageActive) Type() string {
	return "asset.image.active"
}

func (e ImageActive) AggregateType() string {
	return "image"
}

func (e ImageActive) Tags() []string {
	return append([]string(nil), e.tags...)
}

func (e ImageActive) Metadata() map[string]string {
	return maps.Clone(e.metadata)
}
