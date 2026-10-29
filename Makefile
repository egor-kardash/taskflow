-include .env
export

COMPOSE := docker compose
BIN := ./bin/taskflow
SSLMODE := $(or $(POSTGRES_SSLMODE),disable)

.PHONY: check-env app-build app-run dev dev-down db-up db-down db-reset db-proxy-up db-proxy-down \
        migrate-create migrate-up migrate-down migrate-action

check-env:
	@test -f .env || { echo "error: .env is missing. Copy .env.example to .env and fill it in." >&2; exit 1; }
	@test -n "$(POSTGRES_USER)" || { echo "error: POSTGRES_USER is empty in .env" >&2; exit 1; }
	@test -n "$(POSTGRES_PASSWORD)" || { echo "error: POSTGRES_PASSWORD is empty in .env" >&2; exit 1; }
	@test -n "$(POSTGRES_DB)" || { echo "error: POSTGRES_DB is empty in .env" >&2; exit 1; }

app-build:
	@go build -o $(BIN) ./cmd/taskflow

app-run: check-env app-build
	@$(BIN) --config=./config/config.yml

dev: db-up db-proxy-up migrate-up app-run

dev-down:
	@$(COMPOSE) down

db-up: check-env
	@$(COMPOSE) up -d db

db-down:
	@$(COMPOSE) stop db

db-reset:
	@printf "Delete the environment's volumes and all DB data? [y/N]: "; \
	read ans; \
	if [ "$$ans" = "y" ]; then \
		$(COMPOSE) down -v; \
		echo "Environment stopped, volumes removed."; \
	else \
		echo "Reset cancelled."; \
	fi

db-proxy-up: check-env
	@$(COMPOSE) up -d db-proxy

db-proxy-down:
	@$(COMPOSE) stop db-proxy

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "seq parameter is required. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	mkdir -p migrations; \
	$(COMPOSE) run --rm --no-deps --user $(shell id -u):$(shell id -g) migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@$(MAKE) migrate-action action=up

migrate-down:
	@$(MAKE) migrate-action action=down steps=$(or $(steps),1)

migrate-action: check-env
	@if [ -z "$(action)" ]; then \
		echo "action parameter is required. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	$(COMPOSE) run --rm migrate \
		-path /migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@db:5432/$(POSTGRES_DB)?sslmode=$(SSLMODE)" \
		"$(action)" $(if $(steps),$(steps))
