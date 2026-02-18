# Sistema de Faturamento com WhatsApp

Sistema de cobrança automatizado que permite cadastrar clientes e faturas, enviando lembretes automáticos via WhatsApp antes do vencimento.

## 🚀 Objetivo
Projeto para demonstrar conhecimento em:
- **Go** (Golang) com Clean Architecture
- **RabbitMQ** (Event-Driven Architecture)
- **PostgreSQL** com SQL Puro
- **Docker** & **Docker Compose**
- Integração com **Evolution API** (WhatsApp)

## 🛠️ Stack Tecnológica
- **Linguagem:** Go 1.22+
- **Router:** Chi
- **Banco:** PostgreSQL 15+
- **Mensageria:** RabbitMQ
- **Agendamento:** gocron
- **Logs:** Zap
- **UUID:** google/uuid

## 📂 Estrutura do Projeto
```
billing-system/
├── cmd/               # Executáveis (API, Consumers)
├── internal/          # Código da aplicação (Domain, Usecase, Infra)
├── config/            # Configurações
├── scripts/           # Scripts auxiliares
├── docker/            # Dockerfiles e Compose
└── docs/              # Documentação
```

## 🚦 Como Rodar
### Pré-requisitos
- Go 1.22+
- Docker & Docker Compose
- Make (opcional)

### Passo a Passo
1. Clone o repositório
2. Copie o arquivo de exemplo de ambiente:
   ```bash
   cp .env.example .env
   ```
3. Suba o ambiente Docker:
   ```bash
   docker-compose up -d
   ```
4. Rode as migrações (em breve)
5. Inicie a API (em breve)

## 📝 Documentação
A documentação da API estará disponível em `/swagger/index.html` após iniciar a aplicação.
