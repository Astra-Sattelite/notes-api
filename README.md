# Notes API

A simple REST API built with Go, Chi, PostgreSQL, and Redis.

## Tech Stack
- Go + Chi
- PostgreSQL 17
- Redis 7
- Docker Compose
- Swagger

## Requirements
- Go
- Docker + Docker Compose
- Make (optional)

## Setup

**1. Configure environment**

Create `.env` in the project root:

```env
APP_ENV=dev
PORT=8080

DB_CREDS=postgres://app:app@localhost:5432/notes?sslmode=disable
RDB_CREDS=localhost:6379

SWAGGER_HOST=localhost:8080
```

**2. Start PostgreSQL and Redis**

```bash
docker compose up -d
```

**3. Install dependencies**

```bash
go mod download
```

**4. Run the API**

```bash
go run .
```

API: `http://localhost:8080`

Swagger: `http://localhost:8080/swagger/index.html`

## API Routes

| Method |    Endpoint    |    Description    |
|--------|----------------|-------------------|
| GET    | `/ping`        | Check PostgreSQL  |
| GET    | `/redis-test`  | Check Redis       |
| POST   | `/notes`       | Create a note     |

**Example request:**

```json
{
  "title": "First note",
  "body": "Hello from Notes API"
}
```

## Development Commands

| Command             | Description             |
|---------------------|-------------------------|
| `make run`          | Run the API             |
| `make build`        | Build the project       |
| `make test`         | Run tests               |
| `make check`        | Run checks(tests)       |
| `make fmt`          | Format code             |
| `make tidy`         | Update Go dependencies  |
| `make swagger`      | Regenerate Swagger docs |
| `make docker-up`    | Start PSQL and Redis    |
| `make docker-down`  | Stop services           |
| `make docker-logs`  | View Docker logs        |
| `make docker-clean` | Stop services           |

## Database

Database initialization runs automatically from `db/init.sql` when PostgreSQL creates a new data directory.

To reset the database (deletes all data):

```bash
docker compose down -v
docker compose up -d
```

## Project Structure

```text
.
├── db/          # SQL initialization
├── docs/        # Swagger documentation
├── health/      # PostgreSQL health checks
├── rcache/      # Redis routes
├── main.go      # Entry point
├── compose.yaml
├── Makefile
├── go.mod
└── .env
```