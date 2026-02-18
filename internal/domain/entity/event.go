package entity

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID            string          `json:"id"`
	EventType     string          `json:"event_type"`
	AggregateID   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	EventData     json.RawMessage `json:"event_data"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Timestamp     time.Time       `json:"timestamp"`
	Version       int             `json:"version"`
}

// NewEvent cria uma nova instância de evento
func NewEvent(eventType, aggregateID, aggregateType string, data, metadata json.RawMessage, version int) *Event {
	return &Event{
		EventType:     eventType,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventData:     data,
		Metadata:      metadata,
		Timestamp:     time.Now(),
		Version:       version,
	}
}
