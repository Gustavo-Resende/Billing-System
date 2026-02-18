package repository

import "github.com/teusf/billing-system/internal/domain/entity"

type ClienteRepository interface {
	Save(cliente *entity.Cliente) error
	FindByID(id string) (*entity.Cliente, error)
	FindByWhatsApp(whatsapp string) (*entity.Cliente, error)
	FindAll() ([]*entity.Cliente, error)
	Update(cliente *entity.Cliente) error
	Delete(id string) error
}
