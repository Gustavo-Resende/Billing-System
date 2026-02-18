package configuracao

import (
	"database/sql"
	"fmt"

	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/shared"
)

type ConfiguracaoPostgres struct {
	db shared.DBTX
}

func NewConfiguracaoPostgres(db shared.DBTX) *ConfiguracaoPostgres {
	return &ConfiguracaoPostgres{db: db}
}

func (r *ConfiguracaoPostgres) Save(config *entity.Configuracao) error {
	_, err := r.db.Exec(`
		INSERT INTO configuracoes (id, usuario_id, dias_antes_lembrete, template_lembrete, template_cobranca, horario_inicio_envio, horario_fim_envio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (usuario_id) DO UPDATE SET
			dias_antes_lembrete = EXCLUDED.dias_antes_lembrete,
			template_lembrete = EXCLUDED.template_lembrete,
			template_cobranca = EXCLUDED.template_cobranca,
			horario_inicio_envio = EXCLUDED.horario_inicio_envio,
			horario_fim_envio = EXCLUDED.horario_fim_envio,
			updated_at = EXCLUDED.updated_at
	`,
		config.ID,
		config.UsuarioID,
		config.DiasAntesLembrete,
		config.TemplateLembrete,
		config.TemplateCobranca,
		config.HorarioInicioEnvio,
		config.HorarioFimEnvio,
		config.CreatedAt,
		config.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar configuracao: %w", err)
	}

	return nil
}

func (r *ConfiguracaoPostgres) FindByUsuarioID(usuarioID string) (*entity.Configuracao, error) {
	var c entity.Configuracao
	err := r.db.QueryRow(`
		SELECT id, usuario_id, dias_antes_lembrete, template_lembrete, template_cobranca, horario_inicio_envio, horario_fim_envio, created_at, updated_at
		FROM configuracoes
		WHERE usuario_id = $1
	`, usuarioID).Scan(
		&c.ID, &c.UsuarioID, &c.DiasAntesLembrete, &c.TemplateLembrete, &c.TemplateCobranca, &c.HorarioInicioEnvio, &c.HorarioFimEnvio, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar configuracao: %w", err)
	}

	return &c, nil
}
