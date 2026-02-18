CREATE TABLE IF NOT EXISTS faturas (
    id UUID PRIMARY KEY,
    cliente_id UUID NOT NULL REFERENCES clientes(id),
    numero VARCHAR(50) NOT NULL UNIQUE,
    descricao TEXT,
    valor DECIMAL(10, 2) NOT NULL,
    data_vencimento TIMESTAMP NOT NULL,
    data_pagamento TIMESTAMP,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pendente', 'paga', 'vencida', 'cancelada')),
    lembrete_enviado BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_faturas_cliente_id ON faturas(cliente_id);
CREATE INDEX IF NOT EXISTS idx_faturas_status ON faturas(status);
CREATE INDEX IF NOT EXISTS idx_faturas_data_vencimento ON faturas(data_vencimento);
