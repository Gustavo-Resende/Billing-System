package repository

import "github.com/teusf/billing-system/internal/domain/entity"

type MensagemRepository interface {
	Save(mensagem *entity.Mensagem) error
	FindByID(id string) (*entity.Mensagem, error)
	FindByStatus(status entity.StatusMensagem) ([]*entity.Mensagem, error)
	FindParaDLQ() ([]*entity.Mensagem, error)
	Update(mensagem *entity.Mensagem) error
}
