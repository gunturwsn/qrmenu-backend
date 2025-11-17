SHELL := /bin/bash

# -----------------------------
# Global vars
# -----------------------------
PROJECT ?= qrmenu-dev
ENV ?= .env.dev
DB_URL ?= postgres://qrmenu:qrmenu@postgres:5432/qrmenu_dev?sslmode=disable

# --- Container runtime defaults ---
COLIMA_PROFILE ?= default
COLIMA_RUNTIME ?= docker
COLIMA_CPU ?= 4
COLIMA_MEM ?= 6
COLIMA_VM ?= vz
COLIMA_MOUNT ?= virtiofs
# Docker CLI will ALWAYS use this context (safe even if your global context changes)
DOCKER_CONTEXT ?= colima
DOCKER_COMPOSE_BIN ?= docker-compose

# Runtime prefix matches the desired Docker context (ex: colima, docker, default)
CTX_PREFIX := $(DOCKER_CONTEXT)
CTX_COMMAND_TARGETS := start stop status restart use auto-enable auto-disable
.PHONY: $(foreach suffix,$(CTX_COMMAND_TARGETS),$(CTX_PREFIX)-$(suffix))

ifeq ($(DOCKER_CONTEXT),colima)
CTX_START_HELP := Start Colima VM (runtime=$(COLIMA_RUNTIME), cpu=$(COLIMA_CPU), mem=$(COLIMA_MEM)G)
CTX_STOP_HELP := Stop Colima VM
CTX_STATUS_HELP := Show Colima status and Docker contexts
CTX_RESTART_HELP := Restart Colima VM
CTX_USE_HELP := Ensure Docker context uses '$(DOCKER_CONTEXT)' (and create if missing)
CTX_AUTO_ENABLE_HELP := Autostart Colima on login
CTX_AUTO_DISABLE_HELP := Disable Colima autostart

$(CTX_PREFIX)-start:
	@echo "🚀 Starting Colima (profile=$(COLIMA_PROFILE), runtime=$(COLIMA_RUNTIME))..."
	colima start --profile $(COLIMA_PROFILE) \
		--runtime $(COLIMA_RUNTIME) \
		--vm-type $(COLIMA_VM) \
		--cpu $(COLIMA_CPU) \
		--memory $(COLIMA_MEM) \
		--mount-type $(COLIMA_MOUNT)
	@$(MAKE) -s $(CTX_PREFIX)-use

$(CTX_PREFIX)-stop:
	@echo "🛑 Stopping Colima..."
	colima stop --profile $(COLIMA_PROFILE)

$(CTX_PREFIX)-status:
	colima status --profile $(COLIMA_PROFILE) || true
	@echo
	@echo "Docker contexts:"
	@docker context ls || true

$(CTX_PREFIX)-use:
	@# Ensure a docker context named $(DOCKER_CONTEXT) exists and points to colima socket
	@if ! docker context inspect $(DOCKER_CONTEXT) >/dev/null 2>&1; then \
		echo "ℹ️ Creating docker context '$(DOCKER_CONTEXT)' for Colima..."; \
		docker context create $(DOCKER_CONTEXT) --docker "host=unix://$$HOME/.colima/$(COLIMA_PROFILE)/docker.sock"; \
	fi
	@docker context use $(DOCKER_CONTEXT)
	@echo "✅ Docker context now: $$(docker context show)"

$(CTX_PREFIX)-auto-enable:
	colima start --profile $(COLIMA_PROFILE) --auto
	@echo "✅ Colima autostart enabled."

$(CTX_PREFIX)-auto-disable:
	colima stop --profile $(COLIMA_PROFILE)
	colima delete --profile $(COLIMA_PROFILE) --keep
	@launchctl remove io.colima.$(COLIMA_PROFILE) || true
	@echo "✅ Colima autostart disabled (profile kept). To fully delete: 'colima delete --profile $(COLIMA_PROFILE)'."

else
CTX_START_HELP := No-op start hook for Docker context '$(DOCKER_CONTEXT)' (nothing to start)
CTX_STOP_HELP := No-op stop hook for Docker context '$(DOCKER_CONTEXT)'
CTX_STATUS_HELP := Show Docker contexts
CTX_RESTART_HELP := No-op restart hook for Docker context '$(DOCKER_CONTEXT)'
CTX_USE_HELP := Switch Docker context to '$(DOCKER_CONTEXT)'
CTX_AUTO_ENABLE_HELP := Autostart management not required for '$(DOCKER_CONTEXT)'
CTX_AUTO_DISABLE_HELP := Autostart management not required for '$(DOCKER_CONTEXT)'

$(CTX_PREFIX)-start:
	@echo "ℹ️ Nothing to start for Docker context '$(DOCKER_CONTEXT)'."

$(CTX_PREFIX)-stop:
	@echo "ℹ️ Nothing to stop for Docker context '$(DOCKER_CONTEXT)'."

$(CTX_PREFIX)-status:
	@echo "Docker contexts:"
	@docker context ls || true

$(CTX_PREFIX)-use:
	@if docker context inspect $(DOCKER_CONTEXT) >/dev/null 2>&1; then \
		docker context use $(DOCKER_CONTEXT); \
	else \
		echo "⚠️ Docker context '$(DOCKER_CONTEXT)' not found, keeping current context '$$(docker context show)'."; \
	fi
	@echo "✅ Docker context now: $$(docker context show)"

$(CTX_PREFIX)-auto-enable:
	@echo "ℹ️ Autostart is not managed for Docker context '$(DOCKER_CONTEXT)'."

$(CTX_PREFIX)-auto-disable:
	@echo "ℹ️ Autostart is not managed for Docker context '$(DOCKER_CONTEXT)'."

endif

$(CTX_PREFIX)-restart: $(CTX_PREFIX)-stop $(CTX_PREFIX)-start

DC_DEV  := $(DOCKER_COMPOSE_BIN) -p $(PROJECT) -f docker-compose.dev.yml --env-file $(ENV)
DC_PROD := $(DOCKER_COMPOSE_BIN) -p $(PROJECT)-prod -f docker-compose.prod.yml --env-file $(ENV)

# Default connection settings shared across targets
PGHOST ?= postgres
PGUSER ?= qrmenu
PGDB   ?= qrmenu_dev

# Helper: wait for Postgres to be ready with the expected credentials
define WAIT_FOR_PG
	@printf "⏳ Waiting for postgres to be ready"; \
	until $(DC_DEV) exec postgres pg_isready -h $(PGHOST) -U $(PGUSER) -d $(PGDB) >/dev/null 2>&1; do \
		printf "."; \
		sleep 1; \
	done; \
	printf " done\n"
endef

# -----------------------------
# HELP
# -----------------------------
.PHONY: help
help:
	@echo "Targets:"
	@echo "  $(CTX_PREFIX)-start        $(CTX_START_HELP)"
	@echo "  $(CTX_PREFIX)-stop         $(CTX_STOP_HELP)"
	@echo "  $(CTX_PREFIX)-status       $(CTX_STATUS_HELP)"
	@echo "  $(CTX_PREFIX)-restart      $(CTX_RESTART_HELP)"
	@echo "  $(CTX_PREFIX)-use          $(CTX_USE_HELP)"
	@echo "  $(CTX_PREFIX)-auto-enable  $(CTX_AUTO_ENABLE_HELP)"
	@echo "  $(CTX_PREFIX)-auto-disable $(CTX_AUTO_DISABLE_HELP)"
	@echo "  up-dev / down-dev / logs-dev / ps-dev / config-dev / restart-api"
	@echo "  up-air / down-air / logs-air / exec-air / rebuild-air / clean-air"
	@echo "  up-debug-deps / down-debug-deps / logs-debug-deps"
	@echo "  up-prod / down-prod / logs-prod"
	@echo "  migrate-* / lint / test / redis-cli / pg-cli"

# -----------------------------
# DOCKER CONTEXT MANAGEMENT
# -----------------------------
# Targets defined dynamically above based on DOCKER_CONTEXT.

stop-all:
	$(MAKE) down-dev
	$(MAKE) $(CTX_PREFIX)-stop


# -----------------------------
# DEV (compose.dev)
# -----------------------------
.PHONY: up-dev down-dev logs-dev ps-dev config-dev restart-api
up-dev: $(CTX_PREFIX)-use
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

# ===============================================
# LOCAL DEBUG DEPENDENCIES (Postgres + Redis)
# For local debugging via VS Code (F5). API runs locally, not inside container.
# ===============================================
.PHONY: up-debug-deps down-debug-deps logs-debug-deps
up-debug-deps: $(CTX_PREFIX)-use
	@echo "🟢 Starting local debug dependencies (Postgres + Redis)..."
	$(DC_DEV) up -d postgres redis
	@echo "✅ Dependencies ready! You can now start debugging via VS Code (F5)."
	@echo "💡 API will listen on port $$(grep APP_PORT $(ENV) | cut -d '=' -f2)."

down-debug-deps:
	@echo "🛑 Stopping local debug dependencies..."
	$(DC_DEV) down postgres redis || true
	@echo "✅ Dependencies stopped."

logs-debug-deps:
	@echo "📜 Showing logs for Postgres & Redis..."
	$(DC_DEV) logs -f postgres redis

# -----------------------------
# DEV with Air (hot reload)
# (equivalent to up-dev when the api service already uses Air in Dockerfile.dev)
# -----------------------------
.PHONY: up-air down-air logs-air exec-air rebuild-air clean-air
up-air: $(CTX_PREFIX)-use
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

# Remove containers + network + volumes (⚠️ drops dev DB data)
clean-air:
	$(DC_DEV) down -v

# -----------------------------
# PROD (compose.prod)
# -----------------------------
.PHONY: up-prod down-prod logs-prod
up-prod: $(CTX_PREFIX)-use
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
migrate-app: $(CTX_PREFIX)-use
	# run main.go with the --migrate flag inside the api container
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

migrate-up: $(CTX_PREFIX)-use
	$(DC_DEV) up -d postgres
	$(call WAIT_FOR_PG)
	$(DC_DEV) run --rm --no-deps --entrypoint /migrate migrate \
		-path=/migrations -database "$(DB_URL)" -verbose up

migrate-down: $(CTX_PREFIX)-use
	$(DC_DEV) up -d postgres
	$(call WAIT_FOR_PG)
	$(DC_DEV) run --rm --no-deps --entrypoint /migrate migrate \
		-path=/migrations -database "$(DB_URL)" -verbose down 1

migrate-version: $(CTX_PREFIX)-use
	$(DC_DEV) up -d postgres
	$(call WAIT_FOR_PG)
	$(DC_DEV) run --rm --no-deps --entrypoint /migrate migrate \
		-path=/migrations -database "$(DB_URL)" -verbose version

migrate-force: $(CTX_PREFIX)-use
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
