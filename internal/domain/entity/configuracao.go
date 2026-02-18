package entity

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrUsuarioIDObrigatorio = errors.New("usuario id e obrigatorio")
	ErrDiasInvalidos        = errors.New("dias antes do lembrete deve estar entre 0 e 30")
	ErrFormatoHoraInvalido  = errors.New("formato de hora invalido, use HH:MM")
)

type Configuracao struct {
	BaseEntity
	UsuarioID            string
	DiasAntesLembrete    int
	TemplateLembrete     string
	TemplateCobranca     string
	WhatsAppFinanceiro   string
	EnvioAutomaticoAtivo bool
	HorarioInicioEnvio   string // HH:MM
	HorarioFimEnvio      string // HH:MM
}

func NewConfiguracao(usuarioID string) (*Configuracao, error) {
	c := &Configuracao{
		BaseEntity:           NewBase(),
		UsuarioID:            usuarioID,
		DiasAntesLembrete:    3, // Default
		EnvioAutomaticoAtivo: true,
		HorarioInicioEnvio:   "08:00",
		HorarioFimEnvio:      "18:00",
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Configuracao) Validate() error {
	if c.UsuarioID == "" {
		return ErrUsuarioIDObrigatorio
	}

	if c.DiasAntesLembrete < 0 || c.DiasAntesLembrete > 30 {
		return ErrDiasInvalidos
	}

	timeRegex := regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
	if !timeRegex.MatchString(c.HorarioInicioEnvio) || !timeRegex.MatchString(c.HorarioFimEnvio) {
		return ErrFormatoHoraInvalido
	}

	return nil
}

func (c *Configuracao) EstaDentroHorarioEnvio(agora time.Time) bool {
	// Helper para converter "HH:MM" em minutos desde meia-noite
	minutosDoDia := func(horario string) int {
		t, _ := time.Parse("15:04", horario) // Ignora erro pois validamos no Validate()
		return t.Hour()*60 + t.Minute()
	}

	inicio := minutosDoDia(c.HorarioInicioEnvio)
	fim := minutosDoDia(c.HorarioFimEnvio)
	atual := agora.Hour()*60 + agora.Minute()

	// Caso 1: Janela no mesmo dia (ex: 08:00 as 18:00)
	if inicio <= fim {
		return atual >= inicio && atual <= fim
	}

	// Caso 2: Janela cruza a meia-noite (ex: 22:00 as 06:00)
	return atual >= inicio || atual <= fim
}
