include .env
export

COMPOSE_DEV=infra/dev/docker-compose.yml

# ── Docker ────────────────────────────────────────────────
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
	go build ./...
