package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFatura(t *testing.T) {
	t.Run("should create valid fatura", func(t *testing.T) {
		vencimento := time.Now().AddDate(0, 0, 5) // 5 dias no futuro
		f, err := NewFatura("cust-123", 100.50, vencimento, "Consultoria")

		assert.NoError(t, err)
		assert.NotNil(t, f)
		assert.Equal(t, StatusPendente, f.Status)
		assert.NotEmpty(t, f.Numero)
		assert.False(t, f.LembreteEnviado)
		assert.Nil(t, f.DataPagamento)
	})

	t.Run("should validate invalid value", func(t *testing.T) {
		vencimento := time.Now().AddDate(0, 0, 5)
		f, err := NewFatura("cust-123", 0, vencimento, "Invalid")
		assert.Error(t, err)
		assert.Nil(t, f)
		assert.Equal(t, ErrValorInvalido, err)
	})

	t.Run("should validate past due date", func(t *testing.T) {
		vencimento := time.Now().AddDate(0, 0, -1)
		f, err := NewFatura("cust-123", 100, vencimento, "Late")
		assert.Error(t, err)
		assert.Nil(t, f)
		assert.Equal(t, ErrVencimentoPassado, err)
	})
}

func TestFatura_StateTransitions(t *testing.T) {
	vencimento := time.Now().AddDate(0, 0, 5)
	f, _ := NewFatura("cust-123", 100, vencimento, "Test")

	t.Run("should mark as paid", func(t *testing.T) {
		err := f.MarcarComoPaga()
		assert.NoError(t, err)
		assert.Equal(t, StatusPaga, f.Status)
		assert.NotNil(t, f.DataPagamento)
	})

	t.Run("should fail to pay already paid", func(t *testing.T) {
		err := f.MarcarComoPaga()
		assert.Equal(t, ErrFaturaJaPaga, err)
	})

	t.Run("should fail to cancel paid", func(t *testing.T) {
		err := f.Cancelar()
		assert.Equal(t, ErrCancelarFaturaPaga, err)
	})
}

func TestFatura_Cancelamento(t *testing.T) {
	vencimento := time.Now().AddDate(0, 0, 5)
	f, _ := NewFatura("cust-123", 100, vencimento, "Test")

	err := f.Cancelar()
	assert.NoError(t, err)
	assert.Equal(t, StatusCancelada, f.Status)

	err = f.MarcarComoPaga()
	assert.Equal(t, ErrPagarFaturaCancelada, err)
}

func TestFatura_Lembrete(t *testing.T) {
	// Fatura vence em 2 dias (dentro da janela de 3 dias)
	vencimento := time.Now().AddDate(0, 0, 2)
	f, _ := NewFatura("cust-123", 100, vencimento, "Test")

	// Configurado para avisar 3 dias antes
	// Vence em 2 dias. 2 <= 3. Deve enviar.

	shouldSend := f.DeveEnviarLembrete(3)
	assert.True(t, shouldSend)

	f.MarcarLembreteEnviado()
	assert.True(t, f.LembreteEnviado)
	assert.False(t, f.DeveEnviarLembrete(3))
}

func TestGerarNumeroFatura(t *testing.T) {
	num := GerarNumeroFatura()
	assert.Contains(t, num, "FAT-2") // Check year start
	assert.Len(t, num, 19)           // FAT-YYYYMMDD-XXXXXX (4+8+1+6 = 19 length) formula: FAT + - + 8 chars date + - + 6 digits = 3+1+8+1+6 = 19
	// FAT-20240315-123456
}
