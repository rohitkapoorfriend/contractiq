# ContractIQ

> Production-grade Contract Management System built with Go and Domain-Driven Design

[![CI](https://github.com/contractiq/contractiq/actions/workflows/ci.yml/badge.svg)](https://github.com/contractiq/contractiq/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/contractiq/contractiq)](https://goreportcard.com/report/github.com/contractiq/contractiq)

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
│  │  Commands    │  │   Queries    │  │  DTOs    │               │
│  │ Create,Sign  │  │ Get, List    │  │ Req/Res  │               │
│  │ Approve,...  │  │ w/ filters   │  │          │               │
│  └──────┬──────┘  └──────┬───────┘  └──────────┘               │
├─────────┼────────────────┼──────────────────────────────────────┤
│         ▼                ▼    Domain Layer (Pure Business Logic) │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  Contract Aggregate   │  Template  │  Party │Approval│       │
│  │  ┌──────┐ ┌────────┐ │           │        │        │       │
│  │  │Entity│ │  FSM   │ │  Factory  │  Value  │ Events │       │
│  │  │      │ │Draft──►│ │  Pattern  │Objects  │        │       │
│  │  │      │ │Review─►│ │           │ Money   │        │       │
│  │  │      │ │Approve►│ │           │ Clause  │        │       │
│  │  │      │ │Active─►│ │           │DateRange│        │       │
│  │  │      │ │Expired │ │           │         │        │       │
│  │  └──────┘ └────────┘ │           │         │        │       │
│  └──────────────────────────────────────────────────────┘       │
├─────────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                          │
│  ┌──────────┐ ┌─────────┐ ┌──────────┐ ┌────────┐              │
│  │PostgreSQL│ │Event Bus│ │   JWT    │ │ Config │              │
│  │  Repos   │ │In-Memory│ │  Auth    │ │ Viper  │              │
│  │  + UoW   │ │         │ │+ bcrypt  │ │        │              │
│  └──────────┘ └─────────┘ └──────────┘ └────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

## DDD Patterns Implemented

| Pattern | Implementation |
|---------|---------------|
| **Aggregate Root** | `Contract` with unexported fields, state transitions via methods |
| **Value Objects** | `Money` (cents + currency), `Clause`, `DateRange`, `Status` |
| **Status FSM** | Draft → PendingReview → Approved → Active → Expired/Terminated |
| **Repository Pattern** | Interfaces in domain, PostgreSQL implementations in infra |
| **Domain Events** | ContractCreated, Submitted, Approved, Signed, Terminated |
| **Factory Pattern** | Contract creation with validation, template-based creation |
| **Specification Pattern** | Typed `Filter` criteria for contract queries |
| **CQRS-lite** | Separate command/query handlers for contracts |
| **Unit of Work** | Transaction wrapper providing scoped repositories |
| **Optimistic Concurrency** | Version field on aggregates, checked on update |

## Tech Stack

- **Go 1.22+** with standard project layout
- **Chi** - Lightweight HTTP router
- **PostgreSQL** + **sqlx** - Database with compile-time safe queries
- **golang-migrate** - Database migrations
- **JWT** (golang-jwt) + **bcrypt** - Authentication
- **Zap** - Structured logging
- **Viper** - Configuration management
- **Docker** - Multi-stage build + docker-compose
- **GitHub Actions** - CI pipeline (lint, test, build)

## Project Structure

```
contractiq/
├── cmd/api/main.go              # Entrypoint with graceful shutdown
├── internal/
│   ├── domain/                  # Pure business logic (ZERO infra imports)
│   │   ├── contract/            # Core aggregate: entity, FSM, value objects
│   │   ├── template/            # Contract templates aggregate
│   │   ├── party/               # External parties aggregate
│   │   ├── approval/            # Approval workflow aggregate
│   │   └── event/               # Domain event interfaces
│   ├── application/             # Use cases (CQRS-lite)
│   │   ├── contract/command/    # Write operations
│   │   ├── contract/query/      # Read operations
│   │   ├── contract/dto/        # Request/Response types
│   │   ├── template/            # Template CRUD service
│   │   ├── party/               # Party CRUD service
│   │   └── unitofwork/          # Transaction interface
│   ├── infrastructure/          # External concerns
│   │   ├── persistence/postgres/# Repository + UoW implementations
│   │   ├── eventbus/            # In-memory event dispatcher
│   │   ├── auth/                # JWT, bcrypt, user service
│   │   └── config/              # Viper configuration
│   └── interfaces/http/         # REST API layer
│       ├── handler/             # Request handlers
│       ├── middleware/           # Auth, logging, CORS, rate limit, recovery
│       ├── response/            # JSON response helpers
│       └── validation/          # Request validation
├── pkg/                         # Shared utilities
│   ├── apperror/                # Error types with HTTP mapping
│   ├── identifier/              # UUID helpers
│   └── clock/                   # Time abstraction for testing
├── migrations/                  # 6 migration pairs (up/down)
├── .github/workflows/ci.yml    # CI pipeline
├── Makefile                     # Build, test, lint, Docker commands
├── Dockerfile                   # Multi-stage build
└── docker-compose.yml           # Full local stack
```

## Getting Started

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Make (optional)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/contractiq/contractiq.git
cd contractiq

# Start infrastructure (PostgreSQL)
make docker-up

# Copy environment config
cp .env.example .env

# Run database migrations
make migrate-up

# Start the API server
make dev
```

### Without Docker

```bash
# Set up PostgreSQL manually and update .env
cp .env.example .env

# Download dependencies
go mod download

# Run migrations
make migrate-up

# Build and run
make run
```

## API Reference

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Login and receive JWT |

### Contracts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/contracts` | Create a new contract |
| GET | `/api/v1/contracts` | List contracts (filtered, paginated) |
| GET | `/api/v1/contracts/{id}` | Get contract by ID |
| PUT | `/api/v1/contracts/{id}` | Update a draft contract |
| POST | `/api/v1/contracts/{id}/submit` | Submit for review |
| POST | `/api/v1/contracts/{id}/approve` | Approve contract |
| POST | `/api/v1/contracts/{id}/sign` | Sign and activate |
| POST | `/api/v1/contracts/{id}/terminate` | Terminate contract |

### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/templates` | Create template |
| GET | `/api/v1/templates` | List templates |
| GET | `/api/v1/templates/{id}` | Get template |
| PUT | `/api/v1/templates/{id}` | Update template |
| DELETE | `/api/v1/templates/{id}` | Delete template |

### Parties

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/parties` | Create party |
| GET | `/api/v1/parties` | List parties |
| GET | `/api/v1/parties/{id}` | Get party |
| PUT | `/api/v1/parties/{id}` | Update party |
| DELETE | `/api/v1/parties/{id}` | Delete party |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |

### Example: Create and Sign a Contract

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","name":"John Doe","password":"secret123"}'

# Create contract (use token from register response)
curl -X POST http://localhost:8080/api/v1/contracts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "SaaS Agreement",
    "description": "Annual SaaS subscription",
    "value": {"amount_cents": 1200000, "currency": "USD"},
    "start_date": "2025-01-01T00:00:00Z",
    "end_date": "2025-12-31T00:00:00Z"
  }'

# Submit → Approve → Sign
curl -X POST http://localhost:8080/api/v1/contracts/{id}/submit \
  -H "Authorization: Bearer <token>"

curl -X POST http://localhost:8080/api/v1/contracts/{id}/approve \
  -H "Authorization: Bearer <token>"

curl -X POST http://localhost:8080/api/v1/contracts/{id}/sign \
  -H "Authorization: Bearer <token>"
```

## Development

```bash
# Run tests
make test

# Run tests with coverage
make test-cover

# Lint
make lint

# Vet
make vet

# Build binary
make build

# Clean build artifacts
make clean
```

## License

MIT
