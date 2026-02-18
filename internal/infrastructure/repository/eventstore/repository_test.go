package eventstore

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEventStore(t *testing.T) {
	db, err := testutils.SetupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Salvar Evento", func(t *testing.T) {
		tx, err := db.Begin()
		assert.NoError(t, err)
		defer tx.Rollback()

		repoWithTx := NewEventStorePostgres(tx)

		eventID := uuid.New().String()
		aggregateID := uuid.New().String()
		eventType := "PedidoCriado"
		aggregateType := "Pedido"
		payload := json.RawMessage(`{"valor": 100}`)
		metadata := json.RawMessage(`{"user": "admin"}`)

		event := &entity.Event{
			ID:            eventID,
			EventType:     eventType,
			AggregateID:   aggregateID,
			AggregateType: aggregateType,
			EventData:     payload,
			Metadata:      metadata,
			Timestamp:     time.Now(),
			Version:       1,
		}

		err = repoWithTx.Save(event)
		assert.NoError(t, err)

		// Verifica inserção
		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM events WHERE id = $1", eventID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		// Valida dados persistidos
		var savedPayload []byte
		err = tx.QueryRow("SELECT event_data FROM events WHERE id = $1", eventID).Scan(&savedPayload)
		assert.NoError(t, err)
		assert.JSONEq(t, string(payload), string(savedPayload))
	})
}
