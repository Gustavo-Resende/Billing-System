package entity

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type StatusFatura string

const (
	StatusPendente  StatusFatura = "pendente"
	StatusPaga      StatusFatura = "paga"
	StatusVencida   StatusFatura = "vencida"
	StatusCancelada StatusFatura = "cancelada"
)

var (
	ErrValorInvalido        = errors.New("valor deve ser maior que zero")
	ErrVencimentoPassado    = errors.New("data de vencimento deve ser futura")
	ErrFaturaJaPaga         = errors.New("fatura ja esta paga")
	ErrFaturaJaCancelada    = errors.New("fatura ja esta cancelada")
	ErrCancelarFaturaPaga   = errors.New("nao e possivel cancelar uma fatura paga")
	ErrPagarFaturaCancelada = errors.New("nao e possivel pagar uma fatura cancelada")
)

type Fatura struct {
	BaseEntity
	ClienteID       string
	Numero          string
	Descricao       string
	Valor           float64
	DataVencimento  time.Time
	DataPagamento   *time.Time
	Status          StatusFatura
	LembreteEnviado bool
}

func NewFatura(clienteID string, valor float64, dataVencimento time.Time, descricao string) (*Fatura, error) {
	f := &Fatura{
		BaseEntity:      NewBase(),
		ClienteID:       clienteID,
		Numero:          GerarNumeroFatura(),
		Valor:           valor,
		DataVencimento:  dataVencimento,
		Descricao:       descricao,
		Status:          StatusPendente,
		LembreteEnviado: false,
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Fatura) Validate() error {
	// Valor tem que ser positivo
	if f.Valor <= 0 {
		return ErrValorInvalido
	}

	// Verifica se é uma NOVA fatura sendo criada agora
	isNewInvoice := time.Since(f.CreatedAt) < time.Second

	// Se é UMA NOVA fatura, ela OBRIGATORIAMENTE tem que vencer no futuro.
	if isNewInvoice && f.DataVencimento.Before(time.Now()) {
		return ErrVencimentoPassado
	}

	return nil
}

func (f *Fatura) MarcarComoPaga() error {
	if f.Status == StatusPaga {
		return ErrFaturaJaPaga
	}
	if f.Status == StatusCancelada {
		return ErrPagarFaturaCancelada
	}

	now := time.Now()
	f.Status = StatusPaga
	f.DataPagamento = &now
	f.Touch()
	return nil
}

func (f *Fatura) MarcarComoVencida() {
	if f.Status == StatusPendente && f.DataVencimento.Before(time.Now()) {
		f.Status = StatusVencida
		f.Touch()
	}
}

func (f *Fatura) Cancelar() error {
	if f.Status == StatusPaga {
		return ErrCancelarFaturaPaga
	}
	if f.Status == StatusCancelada {
		return ErrFaturaJaCancelada
	}

	f.Status = StatusCancelada
	f.Touch()
	return nil
}

func (f *Fatura) MarcarLembreteEnviado() {
	f.LembreteEnviado = true
	f.Touch()
}

func (f *Fatura) DiasAteVencimento() int {
	return int(time.Until(f.DataVencimento).Hours() / 24)
}

func (f *Fatura) DeveEnviarLembrete(diasAntes int) bool {
	// Só envia se estiver PENDENTE e ainda NÃO enviou
	if f.Status != StatusPendente || f.LembreteEnviado {
		return false
	}

	// A janela começa X dias antes do vencimento e vai até o momento do vencimento
	inicioDaJanela := f.DataVencimento.AddDate(0, 0, -diasAntes)
	agora := time.Now()

	// Precisa ser depois do início da janela e antes de vencer
	estaNaJanelaDeEnvio := agora.After(inicioDaJanela) && agora.Before(f.DataVencimento)

	return estaNaJanelaDeEnvio
}

// GerarNumeroFatura gera um ID legível: FAT-YYYYMMDD-Random
func GerarNumeroFatura() string {
	now := time.Now()
	randNum := rand.Intn(999999)
	return fmt.Sprintf("FAT-%s-%06d", now.Format("20060102"), randNum)
}
