package repository

import "github.com/teusf/billing-system/internal/domain/entity"

type FaturaRepository interface {
	Save(fatura *entity.Fatura) error
	FindByID(id string) (*entity.Fatura, error)
	FindByClienteID(clienteID string) ([]*entity.Fatura, error)
	FindPendentes() ([]*entity.Fatura, error)
	FindVencendoEm(dias int) ([]*entity.Fatura, error)
	Update(fatura *entity.Fatura) error
}
