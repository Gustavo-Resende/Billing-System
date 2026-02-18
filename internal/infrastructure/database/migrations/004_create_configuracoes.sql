CREATE TABLE IF NOT EXISTS configuracoes (
    id UUID PRIMARY KEY,
    usuario_id VARCHAR(100) NOT NULL UNIQUE, -- Pode ser UUID ou string externa
    dias_antes_lembrete INTEGER NOT NULL DEFAULT 3,
    template_lembrete TEXT,
    template_cobranca TEXT,
    whatsapp_financeiro VARCHAR(20),
    envio_automatico_ativo BOOLEAN NOT NULL DEFAULT TRUE,
    horario_inicio_envio VARCHAR(5) NOT NULL DEFAULT '08:00',
    horario_fim_envio VARCHAR(5) NOT NULL DEFAULT '18:00',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
