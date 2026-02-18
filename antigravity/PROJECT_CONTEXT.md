# CONTEXTO DO PROJETO - Sistema de Faturamento com WhatsApp

## OBJETIVO
Desenvolver um sistema de cobrança automatizado que:
- Permite cadastrar clientes e faturas via interface web
- Envia lembretes automáticos via WhatsApp antes do vencimento
- Utiliza arquitetura orientada a eventos (EDA - Event-Driven Architecture)
- Serve como projeto de portfólio para demonstrar conhecimento em sistemas distribuídos
- Frontend será desenvolvido em projeto SEPARADO posteriormente

## STACK TECNOLÓGICA OBRIGATÓRIA

### Backend
- **Linguagem:** Go 1.22+
- **Router:** Chi (go-chi/chi/v5)
- **Banco de Dados:** PostgreSQL 15+ com SQL PURO (❌ SEM GORM ou qualquer ORM)
- **Message Broker:** RabbitMQ
- **Driver PostgreSQL:** lib/pq
- **Driver RabbitMQ:** amqp091-go
- **Agendamento:** gocron
- **Logging:** zap (Uber)
- **WhatsApp API:** Evolution API (container separado)
- **Documentação API:** Swagger (swaggo/swag)
- **UUID:** google/uuid

### Frontend (PROJETO FUTURO SEPARADO)
- Next.js 14+ com App Router
- TypeScript
- Tailwind CSS
- ShadcnUI
- TanStack Query

### Infraestrutura
- Docker & Docker Compose
- PostgreSQL (container)
- RabbitMQ (container)
- Evolution API (container separado)
- Deploy: VPS na Hostinger

## RESTRIÇÕES TÉCNICAS RÍGIDAS

### ❌ PROIBIDO USAR:
- GORM ou qualquer ORM (usar SQL puro com database/sql)
- Redis (pelo menos na primeira versão MVP)
- Kafka (usar RabbitMQ)
- Node.js no backend
- Gin ou Fiber (usar Chi Router)
- Prisma, TypeORM ou similares
- GraphQL (usar REST API)

### ✅ OBRIGATÓRIO:
- SQL escrito manualmente em arquivos .sql
- Migrations em SQL puro
- Repository pattern com queries SQL explícitas
- Clean Architecture (domain, usecase, infrastructure, interface)
- Validações de negócio nas entidades
- Event Sourcing com Event Store
- UUID em todas as entidades (google/uuid)
- BaseEntity para auditoria automática
- Swagger para documentação
- Chi Router para API REST

## ARQUITETURA

### Estrutura de Pastas
```
billing-system/
├── cmd/                              # Executáveis
│   ├── api/                         # API REST
│   │   └── main.go
│   ├── consumer-persistence/        # Consumer que salva no DB
│   │   └── main.go
│   ├── consumer-scheduler/          # Consumer que agenda
│   │   └── main.go
│   └── consumer-notification/       # Consumer que envia WhatsApp
│       └── main.go
│
├── internal/
│   ├── domain/                      # Camada de domínio (regras de negócio)
│   │   ├── entity/                  # Entidades
│   │   │   ├── base.go              # BaseEntity com UUID e auditoria
│   │   │   ├── cliente.go
│   │   │   ├── fatura.go
│   │   │   ├── mensagem.go
│   │   │   └── configuracao.go
│   │   ├── event/                   # Eventos de domínio
│   │   │   ├── fatura_criada.go
│   │   │   ├── fatura_atualizada.go
│   │   │   ├── cliente_criado.go
│   │   │   └── enviar_lembrete.go
│   │   └── repository/              # Interfaces dos repositórios
│   │       ├── cliente_repository.go
│   │       ├── fatura_repository.go
│   │       ├── mensagem_repository.go
│   │       ├── configuracao_repository.go
│   │       └── event_store.go
│   │
│   ├── usecase/                     # Casos de uso (lógica de aplicação)
│   │   ├── cliente/
│   │   │   ├── criar_cliente.go
│   │   │   └── listar_clientes.go
│   │   ├── fatura/
│   │   │   ├── criar_fatura.go
│   │   │   ├── listar_faturas.go
│   │   │   └── atualizar_fatura.go
│   │   └── notification/
│   │       └── enviar_whatsapp.go
│   │
│   ├── infrastructure/              # Implementações concretas
│   │   ├── database/
│   │   │   ├── postgres.go          # Conexão
│   │   │   ├── migrations/          # SQL files
│   │   │   │   ├── 001_create_clientes.sql
│   │   │   │   ├── 002_create_faturas.sql
│   │   │   │   ├── 003_create_mensagens.sql
│   │   │   │   ├── 004_create_configuracoes.sql
│   │   │   │   └── 005_create_events.sql
│   │   │   └── repository/          # Implementações SQL
│   │   │       ├── cliente_postgres.go
│   │   │       ├── fatura_postgres.go
│   │   │       ├── mensagem_postgres.go
│   │   │       ├── configuracao_postgres.go
│   │   │       └── event_store_postgres.go
│   │   │
│   │   ├── messaging/               # RabbitMQ
│   │   │   ├── rabbitmq.go          # Conexão
│   │   │   ├── publisher.go         # Producer
│   │   │   ├── consumer.go          # Consumer base
│   │   │   └── config.go            # Exchanges/Queues
│   │   │
│   │   ├── whatsapp/               # Evolution API client
│   │   │   ├── evolution_client.go
│   │   │   └── template.go
│   │   │
│   │   └── scheduler/              # gocron
│   │       └── cron.go
│   │
│   └── interface/                   # Adapters
│       ├── http/                    # REST API
│       │   ├── server.go
│       │   ├── middleware.go
│       │   ├── handler/
│       │   │   ├── cliente_handler.go
│       │   │   ├── fatura_handler.go
│       │   │   └── dashboard_handler.go
│       │   └── dto/
│       │       ├── cliente_dto.go
│       │       └── fatura_dto.go
│       │
│       └── consumer/                # Event handlers
│           ├── persistence_handler.go
│           ├── scheduler_handler.go
│           └── notification_handler.go
│
├── config/                          # Arquivos de configuração
│   ├── config.yaml
│   └── config.dev.yaml
│
├── scripts/                         # Scripts auxiliares
│   ├── setup-rabbitmq.sh
│   └── run-migrations.sh
│
├── docker/                          # Dockerfiles
│   ├── Dockerfile.api
│   ├── Dockerfile.consumer
│   └── docker-compose.yml
│
├── docs/                            # Documentação Swagger
│   └── swagger.json (gerado)
│
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Padrões de Projeto Obrigatórios

1. **Repository Pattern**
   - Interface na camada domain
   - Implementação com SQL puro na infrastructure

2. **Factory Pattern**
   - Funções `New*()` para criar entidades válidas
   - Ex: `NewCliente()`, `NewFatura()`

3. **Dependency Injection**
   - Passar dependências via construtor
   - Não usar variáveis globais

4. **Event-Driven Architecture**
   - Producer publica eventos no RabbitMQ
   - Múltiplos consumers escutam eventos
   - Event Store registra tudo

5. **Clean Architecture**
   - Domain não depende de nada externo
   - Use cases orquestram lógica
   - Infrastructure implementa detalhes
   - Interface adapta entrada/saída

6. **BaseEntity Pattern**
   - Todos entidades embedam BaseEntity
   - UUID automático
   - CreatedAt/UpdatedAt automático

## FLUXO DE EVENTOS

```
1. API REST recebe POST /api/faturas
2. Use case valida e cria entidade Fatura
3. Producer publica evento "FaturaCriada" no RabbitMQ
4. Exchange roteia para múltiplas queues
5. Consumers processam em paralelo:
   - Persistence: salva no PostgreSQL + Event Store
   - Scheduler: agenda lembrete para 3 dias antes
6. No dia agendado: Scheduler publica "EnviarLembrete"
7. Notification Consumer:
   - Renderiza template de mensagem
   - Envia WhatsApp via Evolution API
   - Registra "MensagemEnviada" no Event Store
8. Retry automático em caso de falha (max 5 tentativas)
9. Após 5 falhas: mensagem vai para DLQ
```

## MODELAGEM DE DADOS

### BaseEntity (Todas entidades herdam)

```go
// internal/domain/entity/base.go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type BaseEntity struct {
    ID        string    // UUID v4
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewBase() BaseEntity {
    now := time.Now()
    return BaseEntity{
        ID:        uuid.New().String(),
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func (b *BaseEntity) Touch() {
    b.UpdatedAt = time.Now()
}

func (b *BaseEntity) GetID() string {
    return b.ID
}
```

### Entidades Principais

**Cliente:**
```go
type Cliente struct {
    BaseEntity           // ID, CreatedAt, UpdatedAt (UUID automático)
    Nome       string    // min 3 chars
    WhatsApp   string    // apenas números, 10-15 dígitos, unique
    Email      string    // opcional, validação regex
    Ativo      bool      // default true
}
```

**Fatura:**
```go
type Fatura struct {
    BaseEntity                  // ID, CreatedAt, UpdatedAt (UUID automático)
    ClienteID       string      // UUID, FK para clientes
    Numero          string      // único, formato: FAT-YYYYMMDD-HHMMSS
    Descricao       string      // opcional
    Valor           float64     // > 0
    DataVencimento  time.Time   // não pode ser passado
    DataPagamento   *time.Time  // nullable, quando foi pago
    Status          StatusFatura // pendente, paga, vencida, cancelada
    LembreteEnviado bool        // default false
}
```

**Mensagem:**
```go
type Mensagem struct {
    BaseEntity                     // ID, CreatedAt, UpdatedAt (UUID automático)
    FaturaID        string         // UUID, FK
    ClienteID       string         // UUID, FK
    WhatsApp        string         // número do destinatário
    Tipo            TipoMensagem   // lembrete, confirmacao, cobranca
    Conteudo        string         // texto renderizado, max 4096 chars
    Status          StatusMensagem // pendente, enviada, falha
    TentativasEnvio int            // contador para DLQ
    ErroMensagem    string         // nullable, guarda erro quando falha
    EnviadoEm       *time.Time     // nullable, quando Evolution confirmou
}
```

**Configuracao:**
```go
type Configuracao struct {
    BaseEntity                   // ID, CreatedAt, UpdatedAt (UUID automático)
    UsuarioID            string  // UUID, unique
    DiasAntesLembrete    int     // default 3, range 0-30
    TemplateLembrete     string  // template Go text/template
    TemplateCobranca     string  // template Go text/template
    WhatsAppFinanceiro   string  // número do financeiro
    EnvioAutomaticoAtivo bool    // default true
    HorarioInicioEnvio   string  // formato HH:MM, ex: "08:00"
    HorarioFimEnvio      string  // formato HH:MM, ex: "18:00"
}
```

**Event Store:**
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,        -- "FaturaCriada", "MensagemEnviada"
    aggregate_id VARCHAR(100) NOT NULL,      -- ID da entidade relacionada
    aggregate_type VARCHAR(50) NOT NULL,     -- "Fatura", "Cliente", "Mensagem"
    event_data JSONB NOT NULL,               -- Dados do evento
    metadata JSONB,                          -- Metadados extras
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INTEGER NOT NULL DEFAULT 1
);
```

### Explicações dos Campos

**Numero (Fatura):**
- Identificador amigável da fatura
- Ex: FAT-20240315-143022
- Para cliente saber qual fatura está consultando
- Único por fatura

**DataVencimento vs DataPagamento:**
- **DataVencimento:** Prazo para pagar (definido na criação)
- **DataPagamento:** Quando realmente foi pago (NULL se não pagou)
- Permite calcular atraso: `DataPagamento - DataVencimento`

**Status (Mensagem):**
- **FASE 1 (MVP):** Apenas `pendente`, `enviada`, `falha`
- Campos `EntregueEm` e `LidoEm` ficam para fase futura
- WhatsApp nem sempre confirma entrega/leitura

**ErroMensagem:**
- Guarda texto do erro quando falha
- Ex: "timeout após 30s", "número inválido"
- Essencial para debug e analytics

**Configuracao vs Mensagem:**
- **Configuracao:** Template/molde de COMO será a mensagem
- **Mensagem:** Conteúdo renderizado que FOI enviado (log)
- São propósitos totalmente diferentes

**Event Store:**
- Histórico completo e imutável de tudo que aconteceu
- Permite auditoria, replay e debug
- NUNCA deletar eventos (append-only)

## REGRAS DE NEGÓCIO CRÍTICAS

### Validações
- Todas validações DEVEM estar em métodos `Validate()` das entidades
- Factory methods `New*()` SEMPRE chamam Validate()
- Nunca criar entidades inválidas

### WhatsApp
- Apenas números (sem +, espaços, parênteses)
- Formato: 5511999998888
- Validação: regex `^\d{10,15}$`
- Único por cliente

### Fatura
- Data vencimento não pode ser passado
- Valor deve ser > 0
- Não pode pagar fatura cancelada
- Não pode cancelar fatura paga
- Número único por fatura

### Mensagem
- Máximo 5 tentativas de envio
- Após 5 falhas → DLQ
- Conteúdo máximo 4096 caracteres
- Sempre registrar em Event Store

### Lembrete
- Enviar 3 dias antes do vencimento (configurável)
- Apenas durante horário comercial (configurável)
- Não enviar se flag `lembrete_enviado = true`
- Apenas faturas com status "pendente"

### Event Store
- NUNCA deletar eventos
- Apenas INSERT (append-only)
- Registrar TODOS eventos importantes:
  - FaturaCriada
  - FaturaAtualizada
  - FaturaPaga
  - ClienteCriado
  - MensagemEnviada
  - MensagemFalhou

## EVOLUTION API - INTEGRAÇÃO

### Arquitetura
```
┌─────────────────────────────┐
│  SEU PROJETO (billing-system)│
│                              │
│  NotificationConsumer ─────┐ │
└────────────────────────────┼─┘
                             │
                             │ HTTP POST
                             ▼
┌──────────────────────────────────┐
│  EVOLUTION API (container)       │
│                                  │
│  ┌────────────────────────────┐ │
│  │  WhatsApp Web Client       │ │
│  │  (escaneia QR Code)        │ │
│  └────────────────────────────┘ │
└──────────────────────────────────┘
```

### Setup no docker-compose.yml
```yaml
services:
  evolution-api:
    image: atendai/evolution-api:latest
    ports:
      - "8081:8081"
    environment:
      - DATABASE_ENABLED=false
      - AUTHENTICATION_API_KEY=sua-chave-secreta-aqui
    volumes:
      - evolution_instances:/evolution/instances
```

### Cliente Go
```go
// infrastructure/whatsapp/evolution_client.go
type EvolutionClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func (c *EvolutionClient) SendMessage(to, message string) error {
    url := fmt.Sprintf("%s/message/sendText/instance1", c.baseURL)
    
    payload := map[string]interface{}{
        "number": to,
        "text":   message,
    }
    
    // HTTP POST com retry e timeout
}
```

**Você NÃO desenvolve o Evolution API!**
- Evolution API é projeto separado (já pronto)
- Você apenas CONSOME a API dele via HTTP
- Roda em container separado
- Configuração via variáveis de ambiente

## TASKS DE DESENVOLVIMENTO

### 🎯 FASE 0: Setup Inicial (2-3 dias)

**TASK 0.1: Inicializar Projeto**
- [ ] Criar repositório Git
- [ ] `go mod init github.com/seu-usuario/billing-system`
- [ ] Instalar dependências básicas
- [ ] Criar estrutura completa de pastas
- [ ] `.gitignore`, `README.md`, `.env.example`

**TASK 0.2: Docker Environment**
- [ ] `docker-compose.yml` com PostgreSQL, RabbitMQ, Evolution API
- [ ] Variáveis de ambiente
- [ ] Script `make docker-up` para subir ambiente
- [ ] Testar conexões

**TASK 0.3: Configuração Básica**
- [ ] Struct de Config para ler `.env`
- [ ] Logger com Zap
- [ ] Health check endpoint básico

---

### 🎯 FASE 1: Domain Layer (3-4 dias)

**TASK 1.1: BaseEntity**
- [ ] Criar `internal/domain/entity/base.go`
- [ ] Struct BaseEntity com UUID
- [ ] Função `NewBase()`
- [ ] Método `Touch()`
- [ ] Testes unitários

**TASK 1.2: Entidade Cliente**
- [ ] Criar `internal/domain/entity/cliente.go`
- [ ] Struct Cliente (embeda BaseEntity)
- [ ] Método `Validate()` com todas validações
- [ ] Factory `NewCliente()`
- [ ] Métodos `Ativar()` e `Desativar()`
- [ ] Testes unitários completos

**TASK 1.3: Entidade Fatura**
- [ ] Criar `internal/domain/entity/fatura.go`
- [ ] Enums `StatusFatura`
- [ ] Struct Fatura (embeda BaseEntity)
- [ ] Método `Validate()`
- [ ] Factory `NewFatura()`
- [ ] Métodos de domínio:
  - `MarcarComoPaga()`
  - `MarcarComoVencida()`
  - `Cancelar()`
  - `MarcarLembreteEnviado()`
  - `EstaVencida()`
  - `DiasAteVencimento()`
  - `DeveEnviarLembrete()`
- [ ] Função `GerarNumeroFatura()`
- [ ] Testes unitários

**TASK 1.4: Entidade Mensagem**
- [ ] Criar `internal/domain/entity/mensagem.go`
- [ ] Enums `TipoMensagem` e `StatusMensagem`
- [ ] Struct Mensagem (embeda BaseEntity)
- [ ] Método `Validate()`
- [ ] Factory `NewMensagem()`
- [ ] Métodos:
  - `MarcarComoEnviada()`
  - `MarcarComoFalha(erro)`
  - `PodeRetentar()`
  - `DeveIrParaDLQ()`
- [ ] Testes unitários

**TASK 1.5: Entidade Configuracao**
- [ ] Criar `internal/domain/entity/configuracao.go`
- [ ] Struct Configuracao (embeda BaseEntity)
- [ ] Método `Validate()`
- [ ] Factory `NewConfiguracao()` com defaults
- [ ] Método `EstaDentroHorarioEnvio()`
- [ ] Funções de templates padrão
- [ ] Testes unitários

**TASK 1.6: Interfaces de Repository**
- [ ] `internal/domain/repository/cliente_repository.go`
- [ ] `internal/domain/repository/fatura_repository.go`
- [ ] `internal/domain/repository/mensagem_repository.go`
- [ ] `internal/domain/repository/configuracao_repository.go`
- [ ] `internal/domain/repository/event_store.go`
- [ ] Definir TODAS as assinaturas de métodos

---

### 🎯 FASE 2: Database Layer (4-5 dias)

**TASK 2.1: Migrations SQL**
- [ ] `001_create_clientes.sql` (com índices)
- [ ] `002_create_faturas.sql` (com enums e índices)
- [ ] `003_create_mensagens.sql` (com enums e índices)
- [ ] `004_create_configuracoes.sql` (com índices)
- [ ] `005_create_events.sql` (com índices GIN)
- [ ] Script para rodar migrations

**TASK 2.2: Conexão PostgreSQL**
- [ ] `internal/infrastructure/database/postgres.go`
- [ ] Função `NewPostgresConnection()`
- [ ] Pool de conexões configurado
- [ ] Ping e health check
- [ ] Função `RunMigrations()`

**TASK 2.3: ClienteRepository (SQL Puro)**
- [ ] `internal/infrastructure/database/repository/cliente_postgres.go`
- [ ] Implementar métodos:
  - `Save(cliente)` - INSERT
  - `FindByID(id)` - SELECT
  - `FindByWhatsApp(whatsapp)` - SELECT
  - `FindAll()` - SELECT
  - `Update(cliente)` - UPDATE
  - `Delete(id)` - DELETE (soft delete)
- [ ] Usar prepared statements
- [ ] Tratar `sql.ErrNoRows`
- [ ] Testes de integração (testcontainers)

**TASK 2.4: FaturaRepository (SQL Puro)**
- [ ] `internal/infrastructure/database/repository/fatura_postgres.go`
- [ ] Implementar métodos:
  - `Save(fatura)` - INSERT
  - `FindByID(id)` - SELECT com JOIN cliente
  - `FindByClienteID(clienteID)` - SELECT
  - `FindAll()` - SELECT com paginação
  - `FindPendentes()` - SELECT WHERE status
  - `FindVencendoEm(dias)` - SELECT WHERE data
  - `FindParaEnviarLembrete()` - SELECT complexo
  - `Update(fatura)` - UPDATE
- [ ] Queries otimizadas
- [ ] Testes de integração

**TASK 2.5: MensagemRepository (SQL Puro)**
- [ ] `internal/infrastructure/database/repository/mensagem_postgres.go`
- [ ] Implementar métodos:
  - `Save(mensagem)` - INSERT
  - `FindByID(id)` - SELECT
  - `FindByFaturaID(faturaID)` - SELECT
  - `FindByStatus(status)` - SELECT
  - `FindParaDLQ()` - SELECT WHERE tentativas >= 5
  - `Update(mensagem)` - UPDATE
- [ ] Testes de integração

**TASK 2.6: ConfiguracaoRepository (SQL Puro)**
- [ ] `internal/infrastructure/database/repository/configuracao_postgres.go`
- [ ] Implementar métodos:
  - `Save(config)` - INSERT
  - `FindByUsuarioID(usuarioID)` - SELECT
  - `Update(config)` - UPDATE
- [ ] Testes de integração

**TASK 2.7: EventStore (SQL Puro)**
- [ ] `internal/infrastructure/database/repository/event_store_postgres.go`
- [ ] Implementar métodos:
  - `Append(event)` - INSERT
  - `FindByAggregateID(id)` - SELECT
  - `FindByEventType(type)` - SELECT
  - `FindAll(limit, offset)` - SELECT
- [ ] JSONB queries
- [ ] Testes de integração

---

### 🎯 FASE 3: RabbitMQ Layer (3-4 dias)

**TASK 3.1: RabbitMQ Connection**
- [ ] `internal/infrastructure/messaging/rabbitmq.go`
- [ ] Função `NewRabbitMQConnection()`
- [ ] Criar exchanges na inicialização
- [ ] Criar queues e bindings
- [ ] Configurar DLQ

**TASK 3.2: Publisher**
- [ ] `internal/infrastructure/messaging/publisher.go`
- [ ] Struct Publisher
- [ ] Método `Publish(eventType, data)`
- [ ] Serialização JSON
- [ ] Publisher confirms
- [ ] Retry em caso de falha
- [ ] Testes de integração

**TASK 3.3: Consumer Base**
- [ ] `internal/infrastructure/messaging/consumer.go`
- [ ] Struct Consumer base
- [ ] Método `Consume(queue, handler)`
- [ ] ACK/NACK automático
- [ ] Prefetch configuration
- [ ] Graceful shutdown
- [ ] Testes de integração

**TASK 3.4: Eventos de Domínio**
- [ ] `internal/domain/event/fatura_criada.go`
- [ ] `internal/domain/event/cliente_criado.go`
- [ ] `internal/domain/event/enviar_lembrete.go`
- [ ] `internal/domain/event/mensagem_enviada.go`
- [ ] Todos com método `ToJSON()`
- [ ] Testes unitários

---

### 🎯 FASE 4: API REST (5-6 dias)

**TASK 4.1: Setup Chi Router**
- [ ] `internal/interface/http/server.go`
- [ ] Configurar Chi router
- [ ] Middleware (CORS, Logger, Recovery)
- [ ] Estrutura de rotas
- [ ] Health check endpoint

**TASK 4.2: Setup Swagger**
- [ ] Instalar swaggo/swag
- [ ] Anotações básicas no main
- [ ] Endpoint `/swagger/*` para docs
- [ ] Script `make swagger` para gerar docs

**TASK 4.3: DTOs**
- [ ] `internal/interface/http/dto/cliente_dto.go`
  - CreateClienteRequest
  - UpdateClienteRequest
  - ClienteResponse
- [ ] `internal/interface/http/dto/fatura_dto.go`
  - CreateFaturaRequest
  - UpdateFaturaRequest
  - FaturaResponse
- [ ] Validação com tags

**TASK 4.4: Use Cases - Cliente**
- [ ] `internal/usecase/cliente/criar_cliente.go`
- [ ] `internal/usecase/cliente/listar_clientes.go`
- [ ] `internal/usecase/cliente/buscar_cliente.go`
- [ ] `internal/usecase/cliente/atualizar_cliente.go`
- [ ] Injetar dependências (repo, publisher)
- [ ] Publicar eventos quando necessário

**TASK 4.5: Use Cases - Fatura**
- [ ] `internal/usecase/fatura/criar_fatura.go`
- [ ] `internal/usecase/fatura/listar_faturas.go`
- [ ] `internal/usecase/fatura/buscar_fatura.go`
- [ ] `internal/usecase/fatura/atualizar_fatura.go`
- [ ] `internal/usecase/fatura/marcar_paga.go`
- [ ] `internal/usecase/fatura/cancelar.go`
- [ ] Publicar eventos

**TASK 4.6: Cliente Handler**
- [ ] `internal/interface/http/handler/cliente_handler.go`
- [ ] POST `/api/clientes` (criar)
- [ ] GET `/api/clientes` (listar)
- [ ] GET `/api/clientes/:id` (buscar)
- [ ] PUT `/api/clientes/:id` (atualizar)
- [ ] DELETE `/api/clientes/:id` (desativar)
- [ ] Anotações Swagger
- [ ] Testes de integração

**TASK 4.7: Fatura Handler**
- [ ] `internal/interface/http/handler/fatura_handler.go`
- [ ] POST `/api/faturas` (criar)
- [ ] GET `/api/faturas` (listar com filtros)
- [ ] GET `/api/faturas/:id` (buscar)
- [ ] PUT `/api/faturas/:id` (atualizar)
- [ ] PATCH `/api/faturas/:id/pagar` (marcar paga)
- [ ] PATCH `/api/faturas/:id/cancelar` (cancelar)
- [ ] Anotações Swagger
- [ ] Testes de integração

**TASK 4.8: Dashboard Handler**
- [ ] `internal/interface/http/handler/dashboard_handler.go`
- [ ] GET `/api/dashboard/metricas`
  - Total de faturas
  - Faturas pendentes
  - Faturas pagas
  - Faturas vencidas
  - Valor total em aberto
- [ ] GET `/api/dashboard/mensagens`
  - Total enviadas
  - Taxa de sucesso
  - Em DLQ

**TASK 4.9: Main da API**
- [ ] `cmd/api/main.go`
- [ ] Inicializar dependências
- [ ] Conectar PostgreSQL
- [ ] Conectar RabbitMQ
- [ ] Iniciar servidor HTTP
- [ ] Graceful shutdown

---

### 🎯 FASE 5: Consumers (4-5 dias)

**TASK 5.1: Persistence Consumer**
- [ ] `cmd/consumer-persistence/main.go`
- [ ] `internal/interface/consumer/persistence_handler.go`
- [ ] Consumir eventos:
  - FaturaCriada → salvar no PostgreSQL
  - ClienteCriado → salvar no PostgreSQL
  - FaturaAtualizada → atualizar no PostgreSQL
- [ ] Salvar TODOS eventos no Event Store
- [ ] ACK/NACK correto
- [ ] Logging estruturado
- [ ] Testes de integração

**TASK 5.2: Scheduler Consumer**
- [ ] `cmd/consumer-scheduler/main.go`
- [ ] `internal/interface/consumer/scheduler_handler.go`
- [ ] Setup gocron
- [ ] Consumir FaturaCriada:
  - Calcular data de lembrete (3 dias antes)
  - Agendar job com gocron
- [ ] No momento agendado:
  - Verificar se deve enviar
  - Publicar EnviarLembrete
- [ ] Job diário para marcar vencidas
- [ ] Testes

**TASK 5.3: Notification Consumer**
- [ ] `cmd/consumer-notification/main.go`
- [ ] `internal/interface/consumer/notification_handler.go`
- [ ] Consumir EnviarLembrete:
  - Buscar configuração do usuário
  - Verificar horário comercial
  - Renderizar template
  - Enviar via Evolution API
  - Salvar mensagem no banco
  - Publicar MensagemEnviada ou MensagemFalhou
- [ ] Retry automático (max 5)
- [ ] Se 5 falhas → NACK (vai pra DLQ)
- [ ] Logging estruturado
- [ ] Testes

---

### 🎯 FASE 6: Evolution API Integration (2-3 dias)

**TASK 6.1: Evolution Client**
- [ ] `internal/infrastructure/whatsapp/evolution_client.go`
- [ ] Struct EvolutionClient
- [ ] Método `SendMessage(to, text)`
- [ ] Método `GetStatus()`
- [ ] HTTP client com:
  - Timeout (30s)
  - Retry (3 tentativas)
  - Exponential backoff
- [ ] Tratamento de erros específicos
- [ ] Testes (mock HTTP)

**TASK 6.2: Template Rendering**
- [ ] `internal/infrastructure/whatsapp/template.go`
- [ ] Função `RenderTemplate(tmpl, data)`
- [ ] Usar `text/template` do Go
- [ ] Suportar placeholders:
  - {{.NomeCliente}}
  - {{.Valor}}
  - {{.DataVencimento}}
  - {{.DiasRestantes}}
  - {{.NomeEmpresa}}
- [ ] Validação de templates
- [ ] Testes unitários

**TASK 6.3: Configurar Evolution no Docker**
- [ ] Adicionar serviço no docker-compose.yml
- [ ] Variáveis de ambiente
- [ ] Volume para instâncias
- [ ] Documentar como escanear QR Code
- [ ] Script para verificar se Evolution está online

---

### 🎯 FASE 7: Scheduler Jobs (2 dias)

**TASK 7.1: Setup gocron**
- [ ] `internal/infrastructure/scheduler/cron.go`
- [ ] Configurar scheduler
- [ ] Timezone
- [ ] Singleton instance

**TASK 7.2: Job de Lembretes**
- [ ] Job que roda a cada 1 hora
- [ ] Buscar faturas com `DeveEnviarLembrete() = true`
- [ ] Para cada fatura:
  - Publicar evento EnviarLembrete
  - Marcar `lembrete_enviado = true`
- [ ] Logging
- [ ] Tratamento de erros

**TASK 7.3: Job de Vencimento**
- [ ] Job que roda diariamente (00:00)
- [ ] Buscar faturas pendentes com vencimento passado
- [ ] Marcar como vencidas
- [ ] Publicar evento FaturaVencida
- [ ] Logging

---

### 🎯 FASE 8: Testes & Qualidade (3-4 dias)

**TASK 8.1: Testes Unitários - Domain**
- [ ] Cobertura > 80% nas entidades
- [ ] Testar todas validações
- [ ] Testar métodos de domínio
- [ ] Table-driven tests

**TASK 8.2: Testes Integração - Repositories**
- [ ] Usar testcontainers (PostgreSQL)
- [ ] Testar CRUD completo
- [ ] Testar queries complexas
- [ ] Testar transações

**TASK 8.3: Testes Integração - RabbitMQ**
- [ ] Usar testcontainers (RabbitMQ)
- [ ] Testar publish/consume
- [ ] Testar ACK/NACK
- [ ] Testar DLQ

**TASK 8.4: Testes E2E**
- [ ] Setup ambiente completo
- [ ] Criar fatura → verificar salvou
- [ ] Criar fatura → verificar lembrete agendado
- [ ] Simular envio WhatsApp
- [ ] Verificar Event Store

**TASK 8.5: Code Quality**
- [ ] golangci-lint
- [ ] go fmt
- [ ] go vet
- [ ] gosec (security)
- [ ] Makefile com comandos

---

### 🎯 FASE 9: DevOps & Deploy (3-4 dias)

**TASK 9.1: Dockerização**
- [ ] Multi-stage Dockerfile otimizado
- [ ] .dockerignore
- [ ] Imagens pequenas (alpine)
- [ ] docker-compose para produção

**TASK 9.2: CI/CD**
- [ ] GitHub Actions workflow
- [ ] Lint + Test + Build
- [ ] Build de imagens Docker
- [ ] Deploy automático VPS Hostinger

**TASK 9.3: Monitoramento Básico**
- [ ] Prometheus metrics endpoint
- [ ] Métricas de negócio:
  - Faturas criadas (counter)
  - Mensagens enviadas (counter)
  - Erros de envio (counter)
  - Latência API (histogram)
- [ ] Health checks robustos

**TASK 9.4: Documentação Final**
- [ ] README.md completo
- [ ] Diagrama de arquitetura
- [ ] Como rodar localmente
- [ ] Como fazer deploy
- [ ] Troubleshooting
- [ ] Documentação Swagger completa

---

## CONVENÇÕES DE CÓDIGO

### Nomenclatura
- **Arquivos:** snake_case (`cliente_repository.go`)
- **Funções/Métodos:** PascalCase (`NewCliente`, `FindByID`)
- **Variáveis:** camelCase (`clienteID`, `faturaRepo`)
- **Constantes:** PascalCase ou UPPER_CASE
- **Interfaces:** Suffixo "Repository", "Service", "Client"

### Organização
- 1 struct/interface por arquivo quando possível
- Agrupar funções relacionadas
- Comentários godoc em funções públicas
- Erros sempre retornados (nunca panic)
- Usar context.Context em operações I/O

### Erros
- Sempre usar `fmt.Errorf` com `%w` para wrapping
- Mensagens descritivas
- Validar entrada no início das funções
- Não ignorar erros (evitar `_`)

### SQL
- Usar prepared statements SEMPRE
- Nomear queries complexas como constantes
- Comentar queries não-óbvias
- Evitar N+1 queries (usar JOINs)
- Usar transações quando necessário

### Logging
- Usar Zap structured logging
- Níveis: Debug, Info, Warn, Error
- Incluir contexto relevante
- Não logar dados sensíveis (senhas, tokens)

## EXEMPLO DE IMPLEMENTAÇÃO ESPERADA

### ✅ CORRETO - Repository com SQL puro

```go
package repository

import (
    "database/sql"
    "fmt"
    "github.com/seu-usuario/billing-system/internal/domain/entity"
)

type PostgresClienteRepository struct {
    db *sql.DB
}

func NewPostgresClienteRepository(db *sql.DB) *PostgresClienteRepository {
    return &PostgresClienteRepository{db: db}
}

func (r *PostgresClienteRepository) Save(cliente *entity.Cliente) error {
    query := `
        INSERT INTO clientes (id, nome, whatsapp, email, ativo, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    _, err := r.db.Exec(
        query,
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

func (r *PostgresClienteRepository) FindByID(id string) (*entity.Cliente, error) {
    query := `
        SELECT id, nome, whatsapp, email, ativo, created_at, updated_at
        FROM clientes
        WHERE id = $1
    `
    
    cliente := &entity.Cliente{}
    err := r.db.QueryRow(query, id).Scan(
        &cliente.ID,
        &cliente.Nome,
        &cliente.WhatsApp,
        &cliente.Email,
        &cliente.Ativo,
        &cliente.CreatedAt,
        &cliente.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("cliente não encontrado: %s", id)
    }
    
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
    }
    
    return cliente, nil
}
```

### ❌ ERRADO - Usando ORM

```go
// NÃO FAZER ISSO!
type ClienteRepository struct {
    db *gorm.DB
}

func (r *ClienteRepository) Save(cliente *entity.Cliente) error {
    return r.db.Create(cliente).Error
}
```

## MELHORES PRÁTICAS COM CURSOR/WINDSURF

### Como Trabalhar de Forma Produtiva

**1. TASKS Pequenas e Incrementais**
```
❌ "Implemente todo o sistema"
✅ "Crie apenas a struct BaseEntity com UUID e testes"
✅ "Implemente ClienteRepository.Save() com SQL puro"
✅ "Crie migration 001_create_clientes.sql"
```

**2. Revise CADA Resposta**
- Leia todo código gerado
- Entenda o que foi feito
- Teste manualmente
- Não aceite cegamente

**3. Peça Explicações**
```
"Por que você usou prepared statement aqui?"
"Explique essa query SQL linha por linha"
"O que esse defer faz?"
"Por que context.Context nessa função?"
```

**4. Iterate Incrementalmente**
```
Ciclo 1: Criar struct
Ciclo 2: Adicionar validações
Ciclo 3: Adicionar testes
Ciclo 4: Refatorar se necessário
```

**5. Comandos Específicos**
```
"Adicione logs estruturados nessa função"
"Adicione tratamento de erro completo aqui"
"Refatore para melhorar legibilidade"
"Adicione comentários godoc"
"Crie testes para essa função"
```

**6. Peça Revisão do SEU Código**
```
"Revise esse código e sugira melhorias"
"Esse SQL está otimizado?"
"Tem algum code smell aqui?"
"Como posso melhorar essa função?"
```

**7. Sempre Teste**
```
Após CADA funcionalidade:
"Crie testes unitários para essa função"
"Adicione teste de caso de erro"
"Teste com dados inválidos"
```

**8. Documente Enquanto Desenvolve**
```
"Adicione comentários explicando essa lógica"
"Crie exemplo de uso dessa função"
"Documente esse endpoint no Swagger"
```

## PRIORIZAÇÃO DE DESENVOLVIMENTO

### 🔴 PRIORIDADE MÁXIMA (Fazer Primeiro)
1. BaseEntity + Entidades de domínio
2. Migrations SQL
3. Repositories com SQL puro
4. API REST básica (CRUD)
5. RabbitMQ Publisher
6. Persistence Consumer

### 🟡 PRIORIDADE ALTA (Fazer Depois)
7. Scheduler Consumer
8. Notification Consumer
9. Evolution API integration
10. gocron jobs
11. Swagger completo

### 🟢 PRIORIDADE MÉDIA (Pode Deixar por Último)
12. Analytics/Dashboard avançado
13. Testes de integração completos
14. CI/CD
15. Monitoramento
16. Otimizações de performance

### ⚪ PRIORIDADE BAIXA (Futuro)
17. Frontend (projeto separado)
18. Webhooks Evolution API
19. Relatórios avançados
20. Multi-tenancy

## PERGUNTAS PARA ESCLARECER

Antes de começar cada TASK, pergunte:

1. **Entendi o objetivo?** O que essa task deve entregar?
2. **Tenho as dependências?** O que precisa estar pronto antes?
3. **Como vou testar?** Qual o critério de sucesso?
4. **Quanto tempo?** Estimativa realista.

## RESUMO EXECUTIVO

Sistema de faturamento com notificações WhatsApp automatizadas usando:
- **Go** com SQL puro (sem ORM)
- **Clean Architecture** bem definida
- **RabbitMQ** para mensageria
- **Event Sourcing** para auditoria
- **Evolution API** para WhatsApp
- **Chi Router** para API REST
- **Swagger** para documentação

Cada componente deve ser desenvolvido **incrementalmente**, **testado** e **documentado** antes de prosseguir para o próximo.

O projeto é educacional, focado em aprendizado de arquitetura distribuída, mensageria e boas práticas Go.

---

## COMEÇAR AGORA

**Primeira TASK sugerida:**
```
TASK 0.1: Inicializar Projeto
- Criar repo Git
- go mod init
- Estrutura de pastas completa
- .gitignore, README.md, .env.example
```

**Está pronto para começar?** 🚀
