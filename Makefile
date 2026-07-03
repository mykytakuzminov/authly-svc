include .env
export

.PHONY: dev-env-up dev-env-down dev-env-rebuild dev-env-clear dev-env-logs
.PHONY: check
.PHONY: generate-proto

# ── Docker ────────────────────────────────────────────────
COMPOSE_DEV := infra/dev/docker-compose.yml

dev-env-up:
	docker compose -f $(COMPOSE_DEV) up -d

dev-env-down:
	docker compose -f $(COMPOSE_DEV) down

dev-env-rebuild:
	docker compose -f $(COMPOSE_DEV) up --build -d

dev-env-clear:
	docker compose -f $(COMPOSE_DEV) down -v

dev-env-logs:
	docker compose -f $(COMPOSE_DEV) logs -f

# ── Code quality ──────────────────────────────────────────
check:
	go mod tidy
	gofmt -w .
	go vet ./...
	golangci-lint run
	go build -o /dev/null ./...

# ── Proto ─────────────────────────────────────────────────
PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

# ── Migrations ────────────────────────────────────────────
MIGRATIONS_DIR := services/auth_service/internal/infra/migrations
DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
