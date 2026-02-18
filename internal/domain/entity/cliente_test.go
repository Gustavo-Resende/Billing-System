package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCliente(t *testing.T) {
	t.Run("should create valid cliente", func(t *testing.T) {
		c, err := NewCliente("John Doe", "5511999998888", "john@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, "John Doe", c.Nome)
		assert.Equal(t, "5511999998888", c.WhatsApp)
		assert.True(t, c.Ativo)
		assert.NotEmpty(t, c.ID)
	})

	t.Run("should validate name length", func(t *testing.T) {
		c, err := NewCliente("Jo", "5511999998888", "john@example.com")
		assert.Error(t, err)
		assert.Nil(t, c)
		assert.Equal(t, ErrNomeCurto, err)
	})

	t.Run("should validate whatsapp format", func(t *testing.T) {
		invalidPhones := []string{
			"123",              // too short
			"1234567890123456", // too long
			"551199999abcd",    // letters
			"+5511999998888",   // special chars
		}

		for _, phone := range invalidPhones {
			c, err := NewCliente("John Doe", phone, "john@example.com")
			assert.Error(t, err)
			assert.Nil(t, c)
			assert.Equal(t, ErrWhatsAppInvalido, err)
		}
	})

	t.Run("should validate email format if provided", func(t *testing.T) {
		c, err := NewCliente("John Doe", "5511999998888", "invalid-email")
		assert.Error(t, err)
		assert.Nil(t, c)
		assert.Equal(t, ErrEmailInvalido, err)
	})

	t.Run("should accept empty email", func(t *testing.T) {
		c, err := NewCliente("John Doe", "5511999998888", "")
		assert.NoError(t, err)
		assert.NotNil(t, c)
	})
}

func TestCliente_AtivarDesativar(t *testing.T) {
	c, _ := NewCliente("John Doe", "5511999998888", "john@example.com")

	oldUpdate := c.UpdatedAt
	time.Sleep(time.Millisecond)

	c.Desativar()
	assert.False(t, c.Ativo)
	assert.True(t, c.UpdatedAt.After(oldUpdate))

	oldUpdate = c.UpdatedAt
	time.Sleep(time.Millisecond)

	c.Ativar()
	assert.True(t, c.Ativo)
	assert.True(t, c.UpdatedAt.After(oldUpdate))
}
