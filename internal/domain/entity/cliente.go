package entity

import (
	"errors"
	"regexp"
)

var (
	ErrNomeCurto        = errors.New("nome deve ter pelo menos 3 digitos")
	ErrWhatsAppInvalido = errors.New("whatsapp deve conter apenas numeros e ter entre 10 e 15 digitos")
	ErrEmailInvalido    = errors.New("email invalido")
)

type Cliente struct {
	BaseEntity
	Nome     string
	WhatsApp string
	Email    string
	Ativo    bool
}

func NewCliente(nome, whatsapp, email string) (*Cliente, error) {
	c := &Cliente{
		BaseEntity: NewBase(),
		Nome:       nome,
		WhatsApp:   whatsapp,
		Email:      email,
		Ativo:      true,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Cliente) Validate() error {
	if len(c.Nome) < 3 {
		return ErrNomeCurto
	}

	// Validate WhatsApp (only numbers, 10-15 digits)
	match, _ := regexp.MatchString(`^\d{10,15}$`, c.WhatsApp)
	if !match {
		return ErrWhatsAppInvalido
	}

	// Validate Email (optional)
	if c.Email != "" {
		match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, c.Email)
		if !match {
			return ErrEmailInvalido
		}
	}

	return nil
}

func (c *Cliente) Ativar() {
	c.Ativo = true
	c.Touch()
}

func (c *Cliente) Desativar() {
	c.Ativo = false
	c.Touch()
}
