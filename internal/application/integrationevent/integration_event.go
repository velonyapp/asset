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
	OccurTime() time.Time
	Tags() []string
}

type BaseIntegrationEvent struct {
	id          string
	aggregateID string
	occurTime   time.Time
	tags        []string
}

func NewBaseIntegrationEvent(aggregateID string, occurTime time.Time, tags []string) BaseIntegrationEvent {
	return BaseIntegrationEvent{
		id:          uuid.Must(uuid.NewV7()).String(),
		aggregateID: aggregateID,
		occurTime:   occurTime,
		tags:        tags,
	}
}

func (e *BaseIntegrationEvent) ID() string           { return e.id }
func (e *BaseIntegrationEvent) AggregateID() string  { return e.aggregateID }
func (e *BaseIntegrationEvent) OccurTime() time.Time { return e.occurTime }
func (e *BaseIntegrationEvent) Tags() []string       { return append([]string(nil), e.tags...) }
