# ContractIQ

> Production-grade Contract Lifecycle Management API built with Go, Domain-Driven Design, and clean architecture principles.

[![CI](https://github.com/contractiq/contractiq/actions/workflows/ci.yml/badge.svg)](https://github.com/contractiq/contractiq/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/contractiq/contractiq)](https://goreportcard.com/report/github.com/contractiq/contractiq)
[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

---

## What Is ContractIQ?

ContractIQ is a **production-ready REST API** for managing the full contract lifecycle — from drafting and review, through approval and signing, to expiry or termination. It is architected with **Domain-Driven Design (DDD)**, a **CQRS-lite** pattern, and a strict **Clean Architecture** layer separation.

Built to demonstrate senior Go engineering practices: aggregate design, domain event publishing, optimistic concurrency, structured logging, graceful shutdown, and full Docker-based deployment.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        HTTP Interface                           │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐ ┌──────┐ ┌──────────┐  │
│  │ Handlers │ │Middleware│ │Response │ │Router│ │Validation│  │
│  └────┬─────┘ └──────────┘ └─────────┘ └──────┘ └──────────┘  │
├───────┼─────────────────────────────────────────────────────────┤
│       ▼         Application Layer (CQRS-lite)                   │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────┐               │
│  │  Commands   │  │   Queries    │  │  DTOs    │               │
│  │ Create,Sign │  │ Get, List    │  │ Req/Res  │               │
│  │ Approve,... │  │ w/ filters   │  │          │               │
│  └──────┬──────┘  └──────┬───────┘  └──────────┘               │
├─────────┼────────────────┼──────────────────────────────────────┤
│         ▼                ▼    Domain Layer (Pure Business Logic) │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  Contract Aggregate  │  Template  │  Party │Approval │       │
│  │  ┌──────┐ ┌───────┐  │            │        │         │       │
│  │  │Entity│ │  FSM  │  │  Factory   │ Value  │  Events │       │
│  │  │      │ │Draft─►│  │  Pattern   │Objects │         │       │
│  │  │      │ │Review►│  │            │ Money  │         │       │
│  │  │      │ │Approv►│  │            │ Clause │         │       │
│  │  │      │ │Active►│  │            │DateRng │         │       │
│  │  │      │ │Expired│  │            │        │         │       │
│  │  └──────┘ └───────┘  │            │        │         │       │
│  └──────────────────────────────────────────────────────┘       │
├─────────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                          │
│  ┌──────────┐ ┌─────────┐ ┌──────────┐ ┌────────┐              │
│  │PostgreSQL│ │EventBus │ │   JWT    │ │ Config │              │
│  │Repos+UoW │ │In-Memory│ │+ bcrypt  │ │ Viper  │              │
│  └──────────┘ └─────────┘ └──────────┘ └────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

---

## DDD Patterns Implemented

| Pattern | Implementation |
|---|---|
| **Aggregate Root** | `Contract` with unexported fields; all mutations via domain methods |
| **Value Objects** | `Money` (cents + ISO currency), `Clause`, `DateRange`, `Status` |
| **Status FSM** | Draft → PendingReview → Approved → Active → Expired / Terminated |
| **Repository Pattern** | Interfaces in domain layer; PostgreSQL implementations in infrastructure |
| **Domain Events** | `ContractCreated`, `Submitted`, `Approved`, `Signed`, `Terminated` |
| **Factory Pattern** | Validated construction + template-based contract creation |
| **Specification / Filter** | Typed `Filter` struct for composable contract queries |
| **CQRS-lite** | Separate command / query handlers for contracts |
| **Unit of Work** | Transaction wrapper that provides scoped, consistent repositories |
| **Optimistic Concurrency** | `version` field on aggregates; enforced via `WHERE id = $1 AND version = $2` |
| **Reconstitute Pattern** | Hydrate aggregates from DB without re-running validation logic |

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP Router | [Chi v5](https://github.com/go-chi/chi) |
| Database | PostgreSQL 16 + [sqlx](https://github.com/jmoiron/sqlx) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Auth | [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) + bcrypt (cost 12) |
| Logging | [Uber Zap](https://github.com/uber-go/zap) |
| Config | [Viper](https://github.com/spf13/viper) + `.env` |
| Validation | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Containerization | Docker (multi-stage) + Docker Compose |
| CI | GitHub Actions |

---

## Project Structure

```
contractiq/
├── cmd/api/main.go                  # Entrypoint — wiring, graceful shutdown
├── internal/
│   ├── domain/                      # Pure business logic — zero infra imports
│   │   ├── contract/                # Core aggregate, FSM, value objects, events
│   │   ├── template/                # Reusable contract template aggregate
│   │   ├── party/                   # External party aggregate
│   │   ├── approval/                # Approval workflow aggregate
│   │   └── event/                   # Domain event interfaces & base type
│   ├── application/                 # Use cases orchestrating domain + infra
│   │   ├── contract/
│   │   │   ├── command/             # Write: Create, Update, Submit, Approve, Sign, Terminate
│   │   │   ├── query/               # Read: GetByID, List (paginated + filtered)
│   │   │   └── dto/                 # Request / Response types
│   │   ├── template/                # Template CRUD service
│   │   ├── party/                   # Party CRUD service
│   │   └── unitofwork/              # UoW interface + Repositories struct
│   ├── infrastructure/
│   │   ├── persistence/postgres/    # Repo + UoW implementations
│   │   ├── eventbus/                # In-memory event dispatcher
│   │   ├── auth/                    # JWT service, bcrypt helpers, UserService
│   │   └── config/                  # Viper-based config loading with defaults
│   └── interfaces/http/
│       ├── handler/                 # HTTP request handlers (contract, template, party, auth, health)
│       ├── middleware/              # Auth JWT, request ID, logging, recovery, rate limit
│       ├── response/                # Standardized JSON response helpers
│       └── validation/              # JSON decode + struct validation
├── pkg/
│   ├── apperror/                    # Typed application errors with HTTP status mapping
│   ├── identifier/                  # UUID generation + validation
│   └── clock/                       # Time abstraction (real + mock for tests)
├── migrations/                      # 6 up/down SQL migration pairs
├── .github/workflows/ci.yml         # Lint → Test → Build CI pipeline
├── Makefile                         # build, run, test, lint, migrate, docker targets
├── Dockerfile                       # Multi-stage build (builder + alpine runtime)
└── docker-compose.yml               # API + PostgreSQL with health checks
```

---

## Getting Started

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- `make` (optional but recommended)

### Quickstart with Docker

```bash
git clone https://github.com/contractiq/contractiq.git
cd contractiq

# Start the full stack (API + PostgreSQL)
make docker-up

# The API is now available at http://localhost:8080
```

### Local Development (without Docker)

```bash
# 1. Start PostgreSQL (or point .env at an existing instance)
cp .env.example .env
# Edit .env with your DB credentials

# 2. Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 3. Run migrations
make migrate-up

# 4. Start the server with live reload
make dev
```

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `APP_ENV` | `development` | `development` or `production` |
| `SERVER_PORT` | `8080` | HTTP listen port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | — | PostgreSQL user |
| `DB_PASSWORD` | — | PostgreSQL password |
| `DB_NAME` | — | PostgreSQL database name |
| `DB_SSL_MODE` | `disable` | `disable`, `require`, `verify-full` |
| `JWT_SECRET` | — | **Required.** HS256 signing secret |
| `JWT_EXPIRY` | `24h` | Token expiry duration |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:3000` | Comma-separated allowed origins |

---

## API Reference

All protected endpoints require `Authorization: Bearer <token>`.

### Authentication

| Method | Endpoint | Body | Description |
|---|---|---|---|
| `POST` | `/api/v1/auth/register` | `{email, name, password}` | Create account, returns JWT |
| `POST` | `/api/v1/auth/login` | `{email, password}` | Authenticate, returns JWT |

### Contracts

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/contracts` | Create a new draft contract |
| `GET` | `/api/v1/contracts` | List contracts (filter by status, party, search; paginated) |
| `GET` | `/api/v1/contracts/{id}` | Get contract by ID |
| `PUT` | `/api/v1/contracts/{id}` | Update a draft contract |
| `POST` | `/api/v1/contracts/{id}/submit` | Submit for review |
| `POST` | `/api/v1/contracts/{id}/approve` | Approve (moves to Approved) |
| `POST` | `/api/v1/contracts/{id}/sign` | Sign and activate |
| `POST` | `/api/v1/contracts/{id}/terminate` | Terminate with reason |

**List query parameters:** `?status=DRAFT&party_id=<uuid>&search=keyword&page=1&page_size=20`

### Templates

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/templates` | Create template |
| `GET` | `/api/v1/templates` | List templates (`?active_only=true`) |
| `GET` | `/api/v1/templates/{id}` | Get template |
| `PUT` | `/api/v1/templates/{id}` | Update template |
| `DELETE` | `/api/v1/templates/{id}` | Delete template |

### Parties

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/parties` | Create party (ORGANIZATION or INDIVIDUAL) |
| `GET` | `/api/v1/parties` | List parties owned by current user |
| `GET` | `/api/v1/parties/{id}` | Get party |
| `PUT` | `/api/v1/parties/{id}` | Update party |
| `DELETE` | `/api/v1/parties/{id}` | Delete party |

### Health

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/health` | Returns API + DB status |

---

## End-to-End Example

```bash
BASE=http://localhost:8080/api/v1

# 1. Register
TOKEN=$(curl -sX POST $BASE/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","name":"Alice","password":"password123"}' \
  | jq -r '.token')

AUTH="Authorization: Bearer $TOKEN"

# 2. Create a contract
CONTRACT_ID=$(curl -sX POST $BASE/contracts \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{
    "title": "Enterprise SaaS Agreement",
    "description": "Annual subscription — 50 seats",
    "value": {"amount_cents": 2400000, "currency": "USD"},
    "start_date": "2025-01-01T00:00:00Z",
    "end_date": "2025-12-31T00:00:00Z"
  }' | jq -r '.id')

# 3. Lifecycle: Draft → PendingReview → Approved → Active
curl -sX POST $BASE/contracts/$CONTRACT_ID/submit  -H "$AUTH" | jq .status
curl -sX POST $BASE/contracts/$CONTRACT_ID/approve -H "$AUTH" | jq .status
curl -sX POST $BASE/contracts/$CONTRACT_ID/sign    -H "$AUTH" | jq .status
# → "ACTIVE"

# 4. Terminate with reason
curl -sX POST $BASE/contracts/$CONTRACT_ID/terminate \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"reason": "Mutual agreement — scope reduced"}' | jq .status
```

---

## Contract Status FSM

```
         ┌─────────────────────────────────────────────────────┐
         │                                                     │
  CREATE │                                                     │
    ▼    │                                                     ▼
  DRAFT ──► PENDING_REVIEW ──► APPROVED ──► ACTIVE ──► EXPIRED
                  │                            │
                  │                            └──► TERMINATED
                  ▼
               DRAFT  (reject back)
```

Invalid transitions return `409 Conflict`.

---

## Development

```bash
make test          # Run all tests with race detector
make test-cover    # Tests + HTML coverage report
make lint          # golangci-lint
make vet           # go vet
make build         # Compile binary to ./bin/contractiq
make clean         # Remove build artifacts

# Database
make migrate-up    # Apply all migrations
make migrate-down  # Roll back one migration
make migrate-create name=add_audit_log  # Create new migration pair
```

---

## Key Design Decisions

**Why `amount_cents int64` instead of `float64`?**  
Floating-point arithmetic is unsafe for money. Storing cents as an integer eliminates rounding errors entirely.

**Why optimistic concurrency instead of row-level locks?**  
Pessimistic locking degrades under read-heavy load. A `version` column incremented on every write catches conflicts without holding DB connections open.

**Why an in-memory event bus?**  
The current publisher is intentionally simple and swappable. Replace `InMemoryPublisher` with a Kafka/NATS adapter behind the same `event.Publisher` interface — no domain code changes required.

**Why CQRS-lite (not full CQRS)?**  
Contracts warrant separated read/write models due to filtering complexity. Templates and Parties are simple CRUD and don't benefit from the overhead of full event sourcing.

---

## Roadmap

- [ ] Refresh token endpoint
- [ ] Role-based access control (owner / reviewer / admin)
- [ ] Contract PDF export
- [ ] Webhook delivery for domain events
- [ ] Redis-backed distributed rate limiter
- [ ] OpenAPI 3.0 spec (auto-generated)
- [ ] Integration test suite with testcontainers

---

## About This Project

This project was built to demonstrate production-level Go API engineering for contract management platforms. It is actively maintained and serves as a portfolio piece for **US remote / contract engineering roles**.

If you're looking for a Go engineer with deep experience in DDD, clean architecture, and production API design — [let's connect](mailto:your@email.com).

---

## License

MIT — see [LICENSE](LICENSE) for details.