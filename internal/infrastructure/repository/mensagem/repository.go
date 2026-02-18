package mensagem

import (
	"database/sql"
	"fmt"

	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/shared"
)

type MensagemPostgres struct {
	db shared.DBTX
}

func NewMensagemPostgres(db shared.DBTX) *MensagemPostgres {
	return &MensagemPostgres{db: db}
}

func (r *MensagemPostgres) Save(msg *entity.Mensagem) error {
	_, err := r.db.Exec(`
		INSERT INTO mensagens (id, fatura_id, cliente_id, whatsapp, tipo, conteudo, status, tentativas_envio, erro_mensagem, enviado_em, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`,
		msg.ID,
		msg.FaturaID,
		msg.ClienteID,
		msg.WhatsApp,
		msg.Tipo,
		msg.Conteudo,
		msg.Status,
		msg.TentativasEnvio,
		msg.ErroMensagem,
		msg.EnviadoEm,
		msg.CreatedAt,
		msg.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar mensagem: %w", err)
	}

	return nil
}

func (r *MensagemPostgres) FindByID(id string) (*entity.Mensagem, error) {
	var m entity.Mensagem
	err := r.db.QueryRow(`
		SELECT id, fatura_id, cliente_id, whatsapp, tipo, conteudo, status, tentativas_envio, erro_mensagem, enviado_em, created_at, updated_at
		FROM mensagens
		WHERE id = $1
	`, id).Scan(
		&m.ID, &m.FaturaID, &m.ClienteID, &m.WhatsApp, &m.Tipo, &m.Conteudo, &m.Status, &m.TentativasEnvio, &m.ErroMensagem, &m.EnviadoEm, &m.CreatedAt, &m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mensagem: %w", err)
	}

	return &m, nil
}

func (r *MensagemPostgres) FindByStatus(status entity.StatusMensagem) ([]*entity.Mensagem, error) {
	rows, err := r.db.Query(`
		SELECT id, fatura_id, cliente_id, whatsapp, tipo, conteudo, status, tentativas_envio, erro_mensagem, enviado_em, created_at, updated_at
		FROM mensagens
		WHERE status = $1
	`, status)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mensagens por status: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *MensagemPostgres) FindParaDLQ() ([]*entity.Mensagem, error) {
	// Mensagens que falharam e excederam tentativas E NÃO SÃO 'enviada'.
	// Mas na verdade, 'FindParaDLQ' seria para processar talvez?
	// Ou somente para listar as que morreram?
	// Vamos assumir que buscamos as que estao com status FALHA e tentativas >= 5
	rows, err := r.db.Query(`
		SELECT id, fatura_id, cliente_id, whatsapp, tipo, conteudo, status, tentativas_envio, erro_mensagem, enviado_em, created_at, updated_at
		FROM mensagens
		WHERE status = $1 AND tentativas_envio >= 5
	`, entity.StatusMensagemFalha)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mensagens para DLQ: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *MensagemPostgres) Update(msg *entity.Mensagem) error {
	_, err := r.db.Exec(`
		UPDATE mensagens
		SET status = $1, tentativas_envio = $2, erro_mensagem = $3, enviado_em = $4, updated_at = $5
		WHERE id = $6
	`,
		msg.Status,
		msg.TentativasEnvio,
		msg.ErroMensagem,
		msg.EnviadoEm,
		msg.UpdatedAt,
		msg.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar mensagem: %w", err)
	}

	return nil
}

func (r *MensagemPostgres) scanRows(rows *sql.Rows) ([]*entity.Mensagem, error) {
	var msgs []*entity.Mensagem
	for rows.Next() {
		var m entity.Mensagem
		if err := rows.Scan(
			&m.ID, &m.FaturaID, &m.ClienteID, &m.WhatsApp, &m.Tipo, &m.Conteudo, &m.Status, &m.TentativasEnvio, &m.ErroMensagem, &m.EnviadoEm, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao scanear mensagem: %w", err)
		}
		msgs = append(msgs, &m)
	}
	return msgs, nil
}
