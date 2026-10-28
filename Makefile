include .env
export

COMPOSE := docker compose

.PHONY: app-run db-up db-down db-reset db-proxy-up db-proxy-down \
        migrate-create migrate-up migrate-down migrate-action

app-run:
	@go run ./cmd/taskflow --config=./config/config.yml

db-up:
	@$(COMPOSE) up -d db

db-down:
	@$(COMPOSE) stop db

db-reset:
	@read -p "Delete the environment's volumes and all DB data? Data will be reset. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		$(COMPOSE) down -v; \
		echo "Environment stopped, volumes removed."; \
	else \
		echo "Reset cancelled."; \
	fi

db-proxy-up:
	@$(COMPOSE) up -d db-proxy

db-proxy-down:
	@$(COMPOSE) stop db-proxy

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "seq parameter is required. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	mkdir -p migrations; \
	docker compose run --rm --user $(shell id -u):$(shell id -g) migrate \
    	create \
    	-ext sql \
    	-dir /migrations \
    	-seq "$(seq)"

migrate-up:
	@$(MAKE) migrate-action action=up

migrate-down:
	@$(MAKE) migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "action parameter is required. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE} \
	"$(action)"
