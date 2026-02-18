package repository

import "github.com/teusf/billing-system/internal/domain/entity"

type ConfiguracaoRepository interface {
	Save(config *entity.Configuracao) error
	FindByUsuarioID(usuarioID string) (*entity.Configuracao, error)
	Update(config *entity.Configuracao) error
}
