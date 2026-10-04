package integrationevent

import (
	"time"

	"github.com/google/uuid"
)

type IntegrationEvent interface {
	ID() string
	Type() string
	AggregateID() string
	AggregateType() string
	Tags() []string
	OccurTime() time.Time
}

type BaseIntegrationEvent struct {
	id          string
	aggregateID string
	occurTime   time.Time
	tags        []string
}

func NewBaseIntegrationEvent(aggregateID string, tags []string, occurTime time.Time) BaseIntegrationEvent {
	return BaseIntegrationEvent{
		id:          uuid.Must(uuid.NewV7()).String(),
		aggregateID: aggregateID,
		tags:        tags,
		occurTime:   occurTime,
	}
}

func (e BaseIntegrationEvent) ID() string           { return e.id }
func (e BaseIntegrationEvent) AggregateID() string  { return e.aggregateID }
func (e BaseIntegrationEvent) Tags() []string       { return append([]string(nil), e.tags...) }
func (e BaseIntegrationEvent) OccurTime() time.Time { return e.occurTime }
