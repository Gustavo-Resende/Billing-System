package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMensagem(t *testing.T) {
	t.Run("should create valid mensagem", func(t *testing.T) {
		m, err := NewMensagem("fat-1", "cli-1", "5511999998888", "Olá", TipoMensagemLembrete)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.Equal(t, StatusMensagemPendente, m.Status)
		assert.Equal(t, 0, m.TentativasEnvio)
	})

	t.Run("should validate required fields", func(t *testing.T) {
		_, err := NewMensagem("fat-1", "cli-1", "", "Olá", TipoMensagemLembrete)
		assert.Equal(t, ErrWhatsAppVazio, err)

		_, err = NewMensagem("fat-1", "cli-1", "5511...", "", TipoMensagemLembrete)
		assert.Equal(t, ErrConteudoVazio, err)
	})
}

func TestMensagem_Lifecycle(t *testing.T) {
	m, _ := NewMensagem("fat-1", "cli-1", "5511999998888", "Olá", TipoMensagemLembrete)

	t.Run("should mark as sent", func(t *testing.T) {
		m.MarcarComoEnviada()
		assert.Equal(t, StatusMensagemEnviada, m.Status)
		assert.NotNil(t, m.EnviadoEm)
		assert.Equal(t, 1, m.TentativasEnvio)
	})
}

func TestMensagem_RetryLogic(t *testing.T) {
	m, _ := NewMensagem("fat-1", "cli-1", "5511", "M", TipoMensagemLembrete)

	// Simula 4 falhas
	for i := 0; i < 4; i++ {
		assert.True(t, m.PodeRetentar())
		assert.False(t, m.DeveIrParaDLQ())
		m.MarcarComoFalha("timeout")
	}

	assert.Equal(t, 4, m.TentativasEnvio)

	// 5ª falha
	assert.True(t, m.PodeRetentar()) // Ainda pode tentar a 5ª vez
	m.MarcarComoFalha("timeout final")

	// Agora esgotou
	assert.Equal(t, 5, m.TentativasEnvio)
	assert.False(t, m.PodeRetentar())
	assert.True(t, m.DeveIrParaDLQ())
	assert.Equal(t, StatusMensagemFalha, m.Status)
	assert.Equal(t, "timeout final", m.ErroMensagem)
}
