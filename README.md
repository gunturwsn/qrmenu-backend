# QRMenu Backend

Backend service for the QRMenu platform. It exposes public APIs for guests to browse tenant menus and place orders, plus admin APIs for managing categories, items, and orders. The service is written in Go (Fiber + GORM), uses PostgreSQL for persistence, Redis for caching, and ships a Swagger/OpenAPI spec.

## Features
- Public endpoints to retrieve menu data by tenant code and to create guest orders.
- Admin endpoints (cookie-authenticated) for CRUD operations on categories, items, options, and orders.
- Multi-tenant support with first-time setup flow for bootstrapping a tenant and its owner admin.
- Redis-backed caching of menu payloads with invalidation helpers.
- JWT-based admin session handling with Fiber middlewares (CORS, structured logging with request IDs, recover, etc.).
- Infrastructure-as-code via Docker Compose for local development, hot reload (Air), and migrations.

## Tech Stack
- **Language:** Go 1.23
- **Framework:** Fiber v2
- **Database:** PostgreSQL 16 (via GORM)
- **Cache:** Redis 7
- **Auth:** JWT
- **Migrations:** `golang-migrate` (Docker image)
- **Documentation:** `openapi/openapi.yaml` + Swagger UI handler
- **Dev tooling:** Air (live reload), Makefile helpers

## Repository Layout
```
cmd/api            # Application entry point
internal/config    # Env & config loader
internal/domain    # Domain models
internal/usecase   # Core business logic (auth, menu, setup, orders, etc.)
internal/repository# Data access with GORM
internal/handler   # HTTP handlers & adapters
internal/transport # Fiber route registration
internal/platform  # Cross-cutting infra (DB, cache, security, conversions)
internal/middleware# HTTP middleware & storage helpers
migrations         # SQL migration files (golang-migrate format)
openapi            # OpenAPI spec
```

The project follows a clean architecture approach:
- Handlers (HTTP adapters) depend only on use case interfaces and translate domain objects into DTOs.
- Use cases own business rules, returning domain entities without leaking HTTP concerns.
- Repositories encapsulate persistence using GORM and expose domain models back to the core.

## Prerequisites
- Go **1.23+**
- Docker & Docker Compose
- GNU Make
- Optional: `golangci-lint` for linting

## Configuration
All environment variables used by the service live in `.env.dev` (local) and `.env.prod.example` (reference). Important keys:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Fiber HTTP port inside the container | `8080` |
| `APP_ALLOWED_ORIGINS` | CORS allow list (comma separated) | `http://localhost:3000,...` |
| `DB_HOST` / `DB_NAME` / `DB_USER` / `DB_PASSWORD` | PostgreSQL connection | see `.env.dev` |
| `DB_URL` | Full DSN used by the migration container | `postgres://...` |
| `DB_MAX_OPEN_CONNS` / `DB_MAX_IDLE_CONNS` | DB connection pool sizes | `25` / `10` (dev) |
| `DB_CONN_MAX_LIFETIME_SEC` / `DB_CONN_MAX_IDLE_TIME_SEC` | Connection lifetime tuning (seconds) | `600` / `300` (dev) |
| `REDIS_ADDR` / `REDIS_DB` / `REDIS_TTL_SECONDS` | Redis connection + cache TTL | `redis:6379`, `0`, `300` |
| `JWT_SECRET` / `JWT_EXPIRES_MINUTES` | Admin JWT signing config | required |
| `SETUP_TOKEN` | Token for initial tenant setup flow | `my-super-secret-token` |

Copy `.env.dev.example` if you need a fresh dev environment file:

```bash
cp .env.dev.example .env.dev
```

## Development Workflow

### Start the stack with Air (hot reload)
```bash
make up-air
```
This builds the API image, launches PostgreSQL + Redis, and runs the API with Air for live reloads. Logs:
```bash
make logs-air
```

To stop the stack:
```bash
make down-air
```

### Start the full dev stack (all services)
```bash
make up-dev
```
Stops & removes everything:
```bash
make down-dev      # containers only
make clean-air     # containers + volumes (drops dev DB)
```

### Exec into containers
```bash
make exec-air      # shell inside api container
make redis-cli     # redis-cli inside redis service
make pg-cli        # psql inside postgres
```

### Database Migrations
We use the official `golang-migrate` Docker image. The Makefile ensures PostgreSQL is healthy before migrating.

```bash
make migrate-up        # apply all pending migrations
make migrate-down      # roll back the latest migration
make migrate-version   # print current schema version
make migrate-force v=20251016220000  # set dirty version
```

Create new migration files:
```bash
make migrate-create name=add_new_table
# creates migrations/2025..._add_new_table.up.sql and .down.sql
```

### Running tests & lint
```bash
make test             # go test ./...
make lint             # golangci-lint run ./...
```

### Direct Go run (outside Docker)
If you want to run the API directly:
```bash
go run cmd/api/main.go
```
Ensure you have PostgreSQL and Redis reachable at the addresses configured in your environment.

## API & Documentation
- The HTTP server boots from `cmd/api/main.go`.
- Routes are registered in `internal/transport/http/route.go`.
- Swagger UI is mounted via `handler.RegisterSwaggerUI`.
- The OpenAPI spec (`openapi/openapi.yaml`) mirrors the handler behaviour; update it whenever endpoints change.

Key public endpoints:
- `GET /api/v1/menu?tenant_code=CODE` – fetch menu (categories + items) by tenant code.
- `POST /api/v1/orders` – create guest order.

Admin endpoints (behind cookie-auth middleware) include:
- `/admin/categories` for category CRUD
- `/admin/items` for item management & stock toggle
- `/admin/orders` for order status updates

Setup endpoints:
- `GET /setup/status?tenant_code=CODE`
- `POST /setup/admin` to bootstrap a tenant’s first admin

## Caching & Invalidations
`MenuUC` caches menu payloads per tenant in Redis. Any admin-side changes should call `MenuUC.InvalidateTenantMenu(code)` to bust the cache. The current handlers invoke the use case via adapter structs; wire invalidation into future mutations as needed.

## Deployment Notes
- Production Compose file: `docker-compose.prod.yml` (single API + PostgreSQL service).
- `docker-compose.prod.yml` expects environment variables in `.env.prod`. Do not commit secrets.
- When building production images, the binary entry point (`./api`) runs with `--migrate` to ensure schema is up-to-date before serving.

## Troubleshooting
- **Redis DNS issue:** ensure the API service depends on Redis (`docker-compose.dev.yml`) so the hostname resolves on startup.
- **Migrations hang:** the Makefile waits for Postgres using `pg_isready`; if a command still hangs, add `-T` to disable TTY or run with `-verbose` for more logs.
- **Go test permission errors:** set a local GOCACHE inside the repo (e.g. `GOCACHE=$PWD/tmp/.gocache go test ./...`) when running under restrictive environments.

---

Happy hacking! Feel free to extend the use cases, add more automated tests, or integrate CI/CD around the provided Make targets and Compose setup.
