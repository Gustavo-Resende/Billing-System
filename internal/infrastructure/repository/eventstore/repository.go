package eventstore

import (
	"encoding/json"
	"fmt"

	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/shared"
)

type EventStorePostgres struct {
	DB shared.DBTX
}

func NewEventStorePostgres(db shared.DBTX) *EventStorePostgres {
	return &EventStorePostgres{DB: db}
}

func (r *EventStorePostgres) Save(event *entity.Event) error {
	query := `
		INSERT INTO events (id, event_type, aggregate_id, aggregate_type, event_data, metadata, timestamp, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// Garante que Metadata seja um JSON válido se for nil ou vazio
	metadata := event.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage("{}")
	}

	_, err := r.DB.Exec(
		query,
		event.ID,
		event.EventType,
		event.AggregateID,
		event.AggregateType,
		event.EventData,
		metadata,
		event.Timestamp,
		event.Version,
	)

	if err != nil {
		return fmt.Errorf("falha ao salvar evento: %w", err)
	}

	return nil
}
