CREATE TABLE IF NOT EXISTS mensagens (
    id UUID PRIMARY KEY,
    fatura_id UUID NOT NULL REFERENCES faturas(id),
    cliente_id UUID NOT NULL REFERENCES clientes(id),
    whatsapp VARCHAR(20) NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    conteudo TEXT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pendente', 'enviada', 'falha')),
    tentativas_envio INTEGER NOT NULL DEFAULT 0,
    erro_mensagem TEXT,
    enviado_em TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mensagens_fatura_id ON mensagens(fatura_id);
CREATE INDEX IF NOT EXISTS idx_mensagens_status ON mensagens(status);
