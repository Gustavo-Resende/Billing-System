.PHONY: up down logs ps run test clean

# Variáveis
DOCKER_COMPOSE_FILE=docker/docker-compose.yml
PROJECT_NAME=billing-system

# Comandos Docker
up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

ps:
	docker-compose -f $(DOCKER_COMPOSE_FILE) ps

clean:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v

# Comandos Go
run:
	go run cmd/api/main.go

test:
	go test ./... -v -p 1

tidy:
	go mod tidy

# Testes de Integração
test-db-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d postgres-test
	@echo "Aguardando banco de teste subir..."
	@timeout 5

test-db-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop postgres-test
	docker-compose -f $(DOCKER_COMPOSE_FILE) rm -f postgres-test
