package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBase(t *testing.T) {
	entity := NewBase()

	assert.NotEmpty(t, entity.ID)
	assert.NotZero(t, entity.CreatedAt)
	assert.NotZero(t, entity.UpdatedAt)
	assert.Equal(t, entity.CreatedAt, entity.UpdatedAt)
}

func TestBaseEntity_Touch(t *testing.T) {
	entity := NewBase()
	oldUpdatedAt := entity.UpdatedAt

	// Wait a bit to ensure time difference
	time.Sleep(time.Millisecond)

	entity.Touch()

	assert.NotEqual(t, oldUpdatedAt, entity.UpdatedAt)
	assert.True(t, entity.UpdatedAt.After(oldUpdatedAt))
}
