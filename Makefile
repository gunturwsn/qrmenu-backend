SHELL := /bin/bash

# -----------------------------
# Global vars
# -----------------------------
PROJECT ?= qrmenu-dev
ENV ?= .env.dev
DB_URL ?= postgres://qrmenu:qrmenu@postgres:5432/qrmenu_dev?sslmode=disable

DC_DEV := docker compose -p $(PROJECT) -f docker-compose.dev.yml --env-file $(ENV)
DC_PROD := docker compose -p $(PROJECT)-prod -f docker-compose.prod.yml --env-file $(ENV)

# Gunakan variabel PG* yang konsisten (fallback ke default yang kita pakai di compose)
PGHOST ?= postgres
PGUSER ?= qrmenu
PGDB   ?= qrmenu_dev

# Helper: tunggu Postgres ready dengan kredensial yang benar
define WAIT_FOR_PG
	@printf "⏳ Waiting for postgres to be ready"; \
	until $(DC_DEV) exec postgres pg_isready -h $(PGHOST) -U $(PGUSER) -d $(PGDB) >/dev/null 2>&1; do \
		printf "."; \
		sleep 1; \
	done; \
	printf " done\n"
endef

# -----------------------------
# DEV (compose.dev)
# -----------------------------
.PHONY: up-dev down-dev logs-dev ps-dev config-dev restart-api
up-dev:
	$(DC_DEV) up -d --build

down-dev:
	$(DC_DEV) down

logs-dev:
	$(DC_DEV) logs -f api

ps-dev:
	$(DC_DEV) ps

config-dev:
	$(DC_DEV) config

restart-api:
	$(DC_DEV) restart api

# -----------------------------
# DEV with Air (hot reload)
# (sama saja dengan up-dev jika service api sudah memakai air di Dockerfile.dev)
# -----------------------------
.PHONY: up-air down-air logs-air exec-air rebuild-air clean-air
up-air:
	$(DC_DEV) up -d --build api
	@echo "✅ Running QRMenu backend with Air (hot reload enabled)"

down-air:
	$(DC_DEV) down

logs-air:
	$(DC_DEV) logs -f api

exec-air:
	$(DC_DEV) exec api sh

rebuild-air:
	$(DC_DEV) build --no-cache api

# Hapus container + network + volume (⚠️ data DB dev hilang)
clean-air:
	$(DC_DEV) down -v

# -----------------------------
# PROD (compose.prod)
# -----------------------------
.PHONY: up-prod down-prod logs-prod
up-prod:
	$(DC_PROD) up -d --build

down-prod:
	$(DC_PROD) down

logs-prod:
	$(DC_PROD) logs -f api

# -----------------------------
# QUALITY
# -----------------------------
.PHONY: lint test
lint:
	golangci-lint run ./...

test:
	go test ./... -v

# -----------------------------
# MIGRATIONS (A) via app (go run --migrate)
# -----------------------------
.PHONY: migrate-app
migrate-app:
	# jalankan main.go dengan flag --migrate dalam container api
	$(DC_DEV) up -d postgres
	$(call WAIT_FOR_PG)
	$(DC_DEV) run --rm --no-deps --entrypoint "" api \
		sh -lc 'go run cmd/api/main.go --migrate'

# -----------------------------
# MIGRATIONS (B) via golang-migrate image
# -----------------------------
.PHONY: migrate-create migrate-up migrate-down migrate-version migrate-force
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=add_table_xxx"; exit 1; fi; \
	ts=$$(date +%Y%m%d%H%M%S); \
	touch migrations/$$ts_$(name).up.sql; \
	touch migrations/$$ts_$(name).down.sql; \
	echo "Created migrations/$$ts_$(name).up.sql and .down.sql"

migrate-up:
	$(DC_DEV) up -d postgres
	$(call WAIT_FOR_PG)
	$(DC_DEV) run --rm --no-deps --entrypoint /migrate migrate \
		-path=/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	$(DC_DEV) up -d postgres
	$(call WAIT_FOR_PG)
	$(DC_DEV) run --rm --no-deps --entrypoint /migrate migrate \
		-path=/migrations -database "$(DB_URL)" -verbose down 1

migrate-version:
	$(DC_DEV) up -d postgres
	$(call WAIT_FOR_PG)
	$(DC_DEV) run --rm --no-deps --entrypoint /migrate migrate \
		-path=/migrations -database "$(DB_URL)" -verbose version

migrate-force:
	@if [ -z "$(v)" ]; then echo "Usage: make migrate-force v=<version>"; exit 1; fi
	$(DC_DEV) up -d postgres
	$(call WAIT_FOR_PG)
	$(DC_DEV) run --rm --no-deps --entrypoint /migrate migrate \
		-path=/migrations -database "$(DB_URL)" -verbose force $(v)


# -----------------------------
# UTILITIES
# -----------------------------
.PHONY: redis-cli pg-cli
redis-cli:
	$(DC_DEV) exec redis redis-cli

pg-cli:
	$(DC_DEV) exec postgres psql -U qrmenu -d qrmenu_dev
