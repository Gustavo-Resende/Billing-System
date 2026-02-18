package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguracao(t *testing.T) {
	t.Run("should create with defaults", func(t *testing.T) {
		c, err := NewConfiguracao("user-1")
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, 3, c.DiasAntesLembrete)
		assert.True(t, c.EnvioAutomaticoAtivo)
		assert.Equal(t, "08:00", c.HorarioInicioEnvio)
		assert.Equal(t, "18:00", c.HorarioFimEnvio)
	})

	t.Run("should validate fields", func(t *testing.T) {
		_, err := NewConfiguracao("")
		assert.Equal(t, ErrUsuarioIDObrigatorio, err)

		c, _ := NewConfiguracao("user-1")
		c.DiasAntesLembrete = 31
		assert.Equal(t, ErrDiasInvalidos, c.Validate())

		c.DiasAntesLembrete = 3
		c.HorarioInicioEnvio = "25:00"
		assert.Equal(t, ErrFormatoHoraInvalido, c.Validate())
	})
}

func TestConfiguracao_HorarioEnvio(t *testing.T) {
	c, _ := NewConfiguracao("user-1")
	c.HorarioInicioEnvio = "09:00"
	c.HorarioFimEnvio = "17:00"

	layout := "15:04"

	// 10:00 -> Sim
	t1, _ := time.Parse(layout, "10:00")
	assert.True(t, c.EstaDentroHorarioEnvio(t1))

	// 08:00 -> Não
	t2, _ := time.Parse(layout, "08:00")
	assert.False(t, c.EstaDentroHorarioEnvio(t2))

	// 17:00 -> Sim (inclusivo)
	t3, _ := time.Parse(layout, "17:00")
	assert.True(t, c.EstaDentroHorarioEnvio(t3))

	// Teste virada de dia
	c.HorarioInicioEnvio = "22:00"
	c.HorarioFimEnvio = "05:00"

	t4, _ := time.Parse(layout, "23:00")
	assert.True(t, c.EstaDentroHorarioEnvio(t4))

	t5, _ := time.Parse(layout, "04:00")
	assert.True(t, c.EstaDentroHorarioEnvio(t5))

	t6, _ := time.Parse(layout, "12:00")
	assert.False(t, c.EstaDentroHorarioEnvio(t6))
}
