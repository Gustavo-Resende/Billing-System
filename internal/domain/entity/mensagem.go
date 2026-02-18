package entity

import (
	"errors"
	"time"
)

type StatusMensagem string
type TipoMensagem string

const (
	StatusMensagemPendente StatusMensagem = "pendente"
	StatusMensagemEnviada  StatusMensagem = "enviada"
	StatusMensagemFalha    StatusMensagem = "falha"

	TipoMensagemLembrete    TipoMensagem = "lembrete"
	TipoMensagemConfirmacao TipoMensagem = "confirmacao"
	TipoMensagemCobranca    TipoMensagem = "cobranca"
)

var (
	ErrWhatsAppVazio = errors.New("whatsapp nao pode ser vazio")
	ErrConteudoVazio = errors.New("conteudo nao pode ser vazio")
)

type Mensagem struct {
	BaseEntity
	FaturaID        string
	ClienteID       string
	WhatsApp        string
	Tipo            TipoMensagem
	Conteudo        string
	Status          StatusMensagem
	TentativasEnvio int
	ErroMensagem    string
	EnviadoEm       *time.Time
}

func NewMensagem(faturaID, clienteID, whatsapp, conteudo string, tipo TipoMensagem) (*Mensagem, error) {
	m := &Mensagem{
		BaseEntity:      NewBase(),
		FaturaID:        faturaID,
		ClienteID:       clienteID,
		WhatsApp:        whatsapp,
		Tipo:            tipo,
		Conteudo:        conteudo,
		Status:          StatusMensagemPendente,
		TentativasEnvio: 0,
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Mensagem) Validate() error {
	if m.WhatsApp == "" {
		return ErrWhatsAppVazio
	}
	if m.Conteudo == "" {
		return ErrConteudoVazio
	}
	return nil
}

func (m *Mensagem) MarcarComoEnviada() {
	now := time.Now()
	m.Status = StatusMensagemEnviada
	m.EnviadoEm = &now
	m.TentativasEnvio++ // Conta como uma tentativa bem sucedida
	m.Touch()
}

func (m *Mensagem) MarcarComoFalha(erro string) {
	m.Status = StatusMensagemFalha // Pode ser temporário se houver retry
	m.ErroMensagem = erro
	m.TentativasEnvio++
	m.Touch()
}

func (m *Mensagem) PodeRetentar() bool {
	return m.TentativasEnvio < 5
}

func (m *Mensagem) DeveIrParaDLQ() bool {
	return m.TentativasEnvio >= 5 && m.Status != StatusMensagemEnviada
}
