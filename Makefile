include .env
export

COMPOSE := docker compose

.PHONY: db-up db-down db-reset proxy-up proxy-down \
        migrate-create migrate-up migrate-down migrate-action

db-up:
	@$(COMPOSE) up -d db

db-down:
	@$(COMPOSE) down db

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
	@$(COMPOSE) down db-proxy

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "seq parameter is required. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	mkdir -p migrations; \
	docker compose run --rm --user $(shell id -u):$(shell id -g) migrate \
    	create \
    	-ext sql \
    	-path /migrations \
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
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable \
	"$(action)"
