package fatura

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/shared"
)

type FaturaPostgres struct {
	db shared.DBTX
}

func NewFaturaPostgres(db shared.DBTX) *FaturaPostgres {
	return &FaturaPostgres{db: db}
}

func (r *FaturaPostgres) Save(fatura *entity.Fatura) error {
	_, err := r.db.Exec(`
		INSERT INTO faturas (id, cliente_id, numero, descricao, valor, data_vencimento, data_pagamento, status, lembrete_enviado, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		fatura.ID,
		fatura.ClienteID,
		fatura.Numero,
		fatura.Descricao,
		fatura.Valor,
		fatura.DataVencimento,
		fatura.DataPagamento,
		fatura.Status,
		fatura.LembreteEnviado,
		fatura.CreatedAt,
		fatura.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar fatura: %w", err)
	}

	return nil
}

func (r *FaturaPostgres) FindByID(id string) (*entity.Fatura, error) {
	var f entity.Fatura
	err := r.db.QueryRow(`
		SELECT id, cliente_id, numero, descricao, valor, data_vencimento, data_pagamento, status, lembrete_enviado, created_at, updated_at
		FROM faturas
		WHERE id = $1
	`, id).Scan(
		&f.ID,
		&f.ClienteID,
		&f.Numero,
		&f.Descricao,
		&f.Valor,
		&f.DataVencimento,
		&f.DataPagamento,
		&f.Status,
		&f.LembreteEnviado,
		&f.CreatedAt,
		&f.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fatura: %w", err)
	}

	return &f, nil
}

func (r *FaturaPostgres) FindByClienteID(clienteID string) ([]*entity.Fatura, error) {
	rows, err := r.db.Query(`
		SELECT id, cliente_id, numero, descricao, valor, data_vencimento, data_pagamento, status, lembrete_enviado, created_at, updated_at
		FROM faturas
		WHERE cliente_id = $1
	`, clienteID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faturas do cliente: %w", err)
	}
	defer rows.Close()

	var faturas []*entity.Fatura
	for rows.Next() {
		var f entity.Fatura
		if err := rows.Scan(
			&f.ID, &f.ClienteID, &f.Numero, &f.Descricao, &f.Valor, &f.DataVencimento, &f.DataPagamento, &f.Status, &f.LembreteEnviado, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao scanear fatura: %w", err)
		}
		faturas = append(faturas, &f)
	}

	return faturas, nil
}

func (r *FaturaPostgres) FindPendentes() ([]*entity.Fatura, error) {
	rows, err := r.db.Query(`
		SELECT id, cliente_id, numero, descricao, valor, data_vencimento, data_pagamento, status, lembrete_enviado, created_at, updated_at
		FROM faturas
		WHERE status = $1
	`, entity.StatusPendente)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faturas pendentes: %w", err)
	}
	defer rows.Close()

	// ... (copiar lógica de scan)
	return r.scanRows(rows)
}

func (r *FaturaPostgres) FindVencendoEm(dias int) ([]*entity.Fatura, error) {
	// A lógica de data pode ser complexa dependendo do banco.
	// PostgreSQL: NOW() + interval 'X days'
	// Mas aqui recebemos 'dias'.
	// Queremos faturas que vencem EXATAMENTE daqui a X dias? Ou ATÉ X dias?
	// Geralmente "Vencendo Em" para lembrete é: DataVencimento = Hoje + dias (ignorando hora)

	targetDate := time.Now().AddDate(0, 0, dias).Format("2006-01-02")

	rows, err := r.db.Query(`
		SELECT id, cliente_id, numero, descricao, valor, data_vencimento, data_pagamento, status, lembrete_enviado, created_at, updated_at
		FROM faturas
		WHERE status = $1 
		AND DATE(data_vencimento) = $2
	`, entity.StatusPendente, targetDate)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faturas vencendo: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *FaturaPostgres) Update(fatura *entity.Fatura) error {
	_, err := r.db.Exec(`
		UPDATE faturas
		SET status = $1, data_pagamento = $2, lembrete_enviado = $3, updated_at = $4
		WHERE id = $5
	`,
		fatura.Status,
		fatura.DataPagamento,
		fatura.LembreteEnviado,
		fatura.UpdatedAt,
		fatura.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar fatura: %w", err)
	}

	return nil
}

// Helper para evitar duplicação de scan
func (r *FaturaPostgres) scanRows(rows *sql.Rows) ([]*entity.Fatura, error) {
	var faturas []*entity.Fatura
	for rows.Next() {
		var f entity.Fatura
		if err := rows.Scan(
			&f.ID, &f.ClienteID, &f.Numero, &f.Descricao, &f.Valor, &f.DataVencimento, &f.DataPagamento, &f.Status, &f.LembreteEnviado, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao scanear fatura: %w", err)
		}
		faturas = append(faturas, &f)
	}
	return faturas, nil
}
