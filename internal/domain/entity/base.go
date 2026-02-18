package entity

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBase() BaseEntity {
	now := time.Now()
	return BaseEntity{
		ID:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *BaseEntity) Touch() {
	b.UpdatedAt = time.Now()
}
