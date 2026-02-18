package cliente

import (
	"database/sql"
	"fmt"

	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/shared"
)

type ClientePostgres struct {
	db shared.DBTX
}

func NewClientePostgres(db shared.DBTX) *ClientePostgres {
	return &ClientePostgres{db: db}
}

func (r *ClientePostgres) Save(cliente *entity.Cliente) error {
	_, err := r.db.Exec(`
		INSERT INTO clientes (id, nome, whatsapp, email, ativo, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		cliente.ID,
		cliente.Nome,
		cliente.WhatsApp,
		cliente.Email,
		cliente.Ativo,
		cliente.CreatedAt,
		cliente.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar cliente: %w", err)
	}

	return nil
}

func (r *ClientePostgres) FindByID(id string) (*entity.Cliente, error) {
	var c entity.Cliente
	err := r.db.QueryRow(`
		SELECT id, nome, whatsapp, email, ativo, created_at, updated_at
		FROM clientes
		WHERE id = $1
	`, id).Scan(
		&c.ID,
		&c.Nome,
		&c.WhatsApp,
		&c.Email,
		&c.Ativo,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Retorna nil se não encontrar, sem erro
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	return &c, nil
}

func (r *ClientePostgres) FindByWhatsApp(whatsapp string) (*entity.Cliente, error) {
	var c entity.Cliente
	err := r.db.QueryRow(`
		SELECT id, nome, whatsapp, email, ativo, created_at, updated_at
		FROM clientes
		WHERE whatsapp = $1
	`, whatsapp).Scan(
		&c.ID,
		&c.Nome,
		&c.WhatsApp,
		&c.Email,
		&c.Ativo,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cliente por whats: %w", err)
	}

	return &c, nil
}

func (r *ClientePostgres) FindAll() ([]*entity.Cliente, error) {
	rows, err := r.db.Query(`
		SELECT id, nome, whatsapp, email, ativo, created_at, updated_at
		FROM clientes
	`)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar clientes: %w", err)
	}
	defer rows.Close()

	var clientes []*entity.Cliente
	for rows.Next() {
		var c entity.Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.WhatsApp, &c.Email, &c.Ativo, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("erro ao scanear cliente: %w", err)
		}
		clientes = append(clientes, &c)
	}

	return clientes, nil
}

func (r *ClientePostgres) Update(cliente *entity.Cliente) error {
	_, err := r.db.Exec(`
		UPDATE clientes
		SET nome = $1, whatsapp = $2, email = $3, ativo = $4, updated_at = $5
		WHERE id = $6
	`,
		cliente.Nome,
		cliente.WhatsApp,
		cliente.Email,
		cliente.Ativo,
		cliente.UpdatedAt,
		cliente.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar cliente: %w", err)
	}

	return nil
}

func (r *ClientePostgres) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM clientes WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("erro ao deletar cliente: %w", err)
	}

	return nil
}
