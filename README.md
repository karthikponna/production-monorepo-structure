# Production Monorepo Structure

A production-ready monorepo starter with a **Go backend** and shared **TypeScript packages**, managed by **Turborepo** and **Bun**.

The backend is a REST API with authentication, PostgreSQL, Redis, background jobs, transactional emails, structured logging, and New Relic observability already wired together. The TypeScript side holds the pieces shared between the backend and any frontend: Zod schemas, the API contract / OpenAPI spec, and React Email templates.

Use it as a starting point so a new project begins with the boring-but-important parts (config, logging, errors, migrations, health checks, graceful shutdown, tests) already in place.

## Tech stack

### Backend (Go)

| Area | Framework / library |
| --- | --- |
| HTTP routing & middleware | [Echo v4](https://echo.labstack.com/) |
| Logging | [zerolog](https://github.com/rs/zerolog) |
| Configuration | [koanf](https://github.com/knadh/koanf) (env provider) + [godotenv](https://github.com/joho/godotenv) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Database driver & pool | [pgx v5](https://github.com/jackc/pgx) (PostgreSQL) |
| Database migrations | [tern](https://github.com/jackc/tern) |
| Cache / key-value store | [go-redis v9](https://github.com/redis/go-redis) |
| Background jobs | [Asynq](https://github.com/hibiken/asynq) (Redis-backed queue) |
| Authentication | [Clerk Go SDK](https://github.com/clerk/clerk-sdk-go) |
| Transactional email | [Resend](https://github.com/resend/resend-go) |
| Observability / APM | [New Relic Go agent](https://github.com/newrelic/go-agent), with integrations for Echo, pgx, Redis, zerolog, and pkg/errors |
| Rate limiting | Echo rate limiter + [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate) |
| Error wrapping | [pkg/errors](https://github.com/pkg/errors) |
| IDs | [google/uuid](https://github.com/google/uuid) |
| Testing | [testify](https://github.com/stretchr/testify) + [Testcontainers](https://golang.testcontainers.org/) (real Postgres in Docker) |
| Task runner | [Task](https://taskfile.dev/) (`Taskfile.yml`) |
| Linting | [golangci-lint](https://golangci-lint.run/) (`.golangci.yml`) |

### Shared packages (TypeScript)

| Area | Framework / library |
| --- | --- |
| Monorepo build orchestration | [Turborepo](https://turbo.build/) |
| Package manager / runtime | [Bun](https://bun.sh/) (workspaces) |
| Language | [TypeScript](https://www.typescriptlang.org/) |
| Schemas & types | [Zod](https://zod.dev/) |
| API contract | [ts-rest](https://ts-rest.com/) |
| OpenAPI generation | `@ts-rest/open-api` + [`@anatine/zod-openapi`](https://github.com/anatine/zod-plugins) |
| Pattern matching | [ts-pattern](https://github.com/gvergnaud/ts-pattern) |
| Email templates | [React Email](https://react.email/) + React 19 |

### Infrastructure you need running

- **PostgreSQL**: primary database
- **Redis**: used by the app and by Asynq for the job queue

## Project structure

```text
.
├── backend/                 Go REST API
│   ├── cmd/
│   │   └── production-monorepo-structure/
│   │       └── main.go      Entry point: loads config, builds server, starts HTTP, graceful shutdown
│   ├── internal/
│   │   ├── config/          Loads and validates config from BOILERPLATE_* env vars (koanf + validator)
│   │   ├── database/        pgx connection pool, tracing, and the tern migrator
│   │   │   └── migrations/  SQL migration files (embedded into the binary)
│   │   ├── errs/            Standard HTTP error types (401, 403, 404, 400, 500) returned as JSON
│   │   ├── handler/         HTTP handlers (health check, OpenAPI docs UI)
│   │   ├── lib/
│   │   │   ├── email/       Resend client + HTML template rendering
│   │   │   ├── job/         Asynq job server, task definitions, and task handlers
│   │   │   └── utils/       Small shared helpers
│   │   ├── logger/          zerolog setup + New Relic log forwarding
│   │   ├── middleware/      Auth (Clerk), request ID, CORS, security headers, request logging,
│   │   │                    panic recovery, rate limiting, New Relic tracing, global error handler
│   │   ├── model/           Base model structs shared by entities
│   │   ├── repository/      Data-access layer (talks to the database)
│   │   ├── router/          Echo router, global middleware chain, and route registration
│   │   ├── server/          Server struct holding config, logger, DB, Redis, and job service
│   │   ├── service/         Business logic layer (e.g. auth service)
│   │   ├── sqlerr/          Converts Postgres errors into friendly HTTP errors
│   │   ├── testing/         Test helpers: Testcontainers Postgres, test server, transactions, assertions
│   │   └── validation/      Request validation helpers
│   ├── .env.sample          Example environment variables
│   ├── .golangci.yml        Linter configuration
│   └── Taskfile.yml         Common commands (run, migrations, tidy)
│
├── packages/                Shared TypeScript packages (Bun workspaces)
│   ├── zod/                 @boilerplate/zod: Zod schemas and types shared across apps
│   ├── openapi/             @boilerplate/openapi: ts-rest API contract and OpenAPI JSON generator
│   └── emails/              @boilerplate/emails: React Email templates, exported to HTML for the backend
│
├── templates/emails/        Exported HTML email templates
├── package.json             Root workspace config and Turborepo scripts
└── turbo.json               Turborepo task pipeline
```

### How a request flows through the backend

```text
HTTP request
  → router (Echo)
  → global middleware: rate limit → CORS → security headers → request ID
                       → New Relic tracing → context enrichment → request logger → recover
  → handler            (parses input, calls a service)
  → service            (business logic)
  → repository         (database access via pgx)
  → response, or an error turned into JSON by the global error handler
```

## Getting started

### Prerequisites

- [Go](https://go.dev/dl/) (version in `backend/go.mod`)
- [Bun](https://bun.sh/) and Node.js 22+
- [Task](https://taskfile.dev/installation/): `brew install go-task`
- PostgreSQL and Redis running locally (Docker is the easiest option)
- Docker, if you want to run the Testcontainers-based tests

### 1. Start PostgreSQL and Redis

```bash
docker run -d --name boilerplate-postgres \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=boilerplate \
  -p 5432:5432 postgres:16-alpine

docker run -d --name boilerplate-redis -p 6379:6379 redis:7-alpine
```

### 2. Install JavaScript dependencies

From the repo root:

```bash
bun install
```

### 3. Configure the backend

```bash
cd backend
cp .env.sample .env
```

Then edit `backend/.env`. At minimum, set:

- `BOILERPLATE_DATABASE.PASSWORD`: your Postgres password (`postgres` if you used the Docker command above)
- `BOILERPLATE_REDIS.ADDRESS`: use `localhost:6379` (host and port, without the `redis://` prefix)
- `BOILERPLATE_AUTH.SECRET_KEY`: your Clerk secret key
- `BOILERPLATE_INTEGRATION.RESEND_API_KEY`: your Resend API key

The New Relic license key is optional. If it is empty, New Relic is skipped.

The `.env` file is loaded automatically at startup. Every variable uses the `BOILERPLATE_` prefix, and the dots map to nested config fields (for example, `BOILERPLATE_DATABASE.HOST` becomes `database.host`).

### 4. Run database migrations

When `BOILERPLATE_PRIMARY.ENV` is anything other than `local`, migrations run automatically on startup. In `local`, run them yourself:

```bash
cd backend
BOILERPLATE_DB_DSN="postgres://postgres:postgres@localhost:5432/boilerplate?sslmode=disable" task migrations:up
```

To create a new migration:

```bash
task migrations:new name=create_users_table
```

### 5. Run the backend

```bash
cd backend
task run
```

The API starts on `http://localhost:8080`. Check it with:

```bash
curl http://localhost:8080/status
```

The response reports the health of the database and Redis.

### 6. Work on the shared packages

```bash
# Build the shared Zod schemas (other packages depend on them)
cd packages/zod && bun run build

# Regenerate the OpenAPI spec from the ts-rest contract
cd packages/openapi && bun run gen

# Preview email templates at http://localhost:3001
cd packages/emails && bun run dev

# Export email templates to HTML for the backend
cd packages/emails && bun run export
```

Or run every package's `dev` script together from the root:

```bash
bun run dev
```

## Useful commands

| Command | Where | What it does |
| --- | --- | --- |
| `task run` | `backend/` | Run the API server |
| `task migrations:new name=<name>` | `backend/` | Create a new SQL migration |
| `task migrations:up` | `backend/` | Apply all migrations (needs `BOILERPLATE_DB_DSN`) |
| `task tidy` | `backend/` | `go fmt`, `go mod tidy`, and `go mod verify` |
| `go test ./...` | `backend/` | Run backend tests (needs Docker for Testcontainers) |
| `golangci-lint run` | `backend/` | Lint the Go code |
| `bun run dev` | root | Run all package dev scripts through Turborepo |
| `bun run build` | root | Build all packages |
| `bun run typecheck` | root | Type-check all packages |

## Adding a new feature

1. Add a migration in `backend/internal/database/migrations/` with `task migrations:new`.
2. Add the model in `internal/model/`.
3. Add data access in `internal/repository/`.
4. Add business logic in `internal/service/`.
5. Add an HTTP handler in `internal/handler/` and register the route in `internal/router/`.
6. Add the matching Zod schema in `packages/zod` and the endpoint in the `packages/openapi` contract, then run `bun run gen` to update the OpenAPI spec.
