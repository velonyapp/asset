package integrationevent

import (
	"maps"
	"time"
)

var _ IntegrationEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseIntegrationEvent

	tags      []string
	metadata  map[string]string
	objectKey string
}

func NewImageCreated(
	aggregateID string,
	occurTime time.Time,
	tags []string,
	metadata map[string]string,
	objectKey string,
) ImageCreated {
	return ImageCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		tags:      append([]string(nil), tags...),
		metadata:  maps.Clone(metadata),
		objectKey: objectKey,
	}
}

func (e ImageCreated) Type() string {
	return "asset.image.created"
}

func (e ImageCreated) AggregateType() string {
	return "image"
}

func (e ImageCreated) Tags() []string {
	return append([]string(nil), e.tags...)
}

func (e ImageCreated) Metadata() map[string]string {
	return maps.Clone(e.metadata)
}

func (e ImageCreated) ObjectKey() string {
	return e.objectKey
}
