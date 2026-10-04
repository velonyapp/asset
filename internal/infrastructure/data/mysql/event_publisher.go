package mysql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/velonyapp/asset/internal/application/integrationevent"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/infrastructure/messaging/protobuf"

	"google.golang.org/protobuf/encoding/protojson"
)

var _ port.EventPublisher = (*eventPublisher)(nil)

type eventPublisher struct {
	db      *sql.DB
	encoder *protobuf.Encoder
}

func NewEventPublisher(
	db *sql.DB,
	encoder *protobuf.Encoder,
) port.EventPublisher {
	return &eventPublisher{
		db:      db,
		encoder: encoder,
	}
}

func (pub *eventPublisher) Publish(ctx context.Context, integrationEvent integrationevent.IntegrationEvent) error {
	event, err := pub.encoder.Encode(integrationEvent)
	if err != nil {
		return err
	}

	tags, err := json.Marshal(event.Tags)
	if err != nil {
		return err
	}
	payload, err := protojson.Marshal(event.Payload)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO outbox_events (
			id,
			type,
			aggregate_id,
			aggregate_type,
			tags,
			payload,
			occur_time
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = executor(ctx, pub.db).ExecContext(ctx, query,
		event.Id,
		event.Type,
		event.AggregateId,
		event.AggregateType,
		string(tags),
		string(payload),
		event.OccurTime.AsTime(),
	)

	return err
}
