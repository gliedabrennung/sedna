# Messenger Core API

A production-ready **Go REST API** backend for a messaging application.  
Demonstrates routing, JWT auth, RBAC, PostgreSQL, migrations, goroutines, graceful shutdown, and OpenAPI docs.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22 |
| HTTP Framework | [Hertz](https://github.com/cloudwego/hertz) |
| Database | PostgreSQL 16 |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Auth | JWT (HS256) |
| Docs | Swagger / OpenAPI (swaggo) |
| Containers | Docker + Docker Compose |

---

## Architecture

```
messenger-core/
├── cmd/
│   ├── server/       # Main application entry point
│   └── seeder/       # Demo data seeder
├── docs/             # Auto-generated Swagger files
├── internal/
│   ├── config/       # Config loading (env vars)
│   ├── controller/http/
│   │   ├── middleware/   # JWT, RoleAuth, RateLimiter, Logger
│   │   ├── auth_handler.go
│   │   ├── message_handler.go
│   │   ├── contact_handler.go
│   │   └── router.go
│   ├── entity/       # Domain models (User, Message, Contact)
│   ├── messenger/    # WebSocket hub & client
│   ├── repository/postgres/  # Data access layer
│   ├── usecase/      # Business logic (+ unit tests)
│   └── pkg/          # Shared helpers (JWT, API responses)
└── migrations/       # Versioned SQL migrations
```

### Entities & Relationships

- **User** — registered accounts with roles (`user`, `admin`)
- **Message** — one-to-many messages between two users
- **Contact** — many-to-many friendship with status (`pending` / `accepted`)

---

## Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- Or: Go 1.22+ and PostgreSQL 16+

### Run with Docker (recommended)

```bash
# 1. Clone the repo
git clone https://github.com/bekzat707/messenger-core
cd messenger-core

# 2. Start server + database
docker compose up --build

# Server is now available at http://localhost:8080
```

### Run Locally

```bash
# 1. Copy and fill in env vars
cp .env.example .env

# 2. Start PostgreSQL (or use existing instance)
# Make sure your DSN in .env points to the right host

# 3. Run the server
go run ./cmd/server

# 4. (Optional) Seed demo data
go run ./cmd/seeder
```

---

## Environment Variables

| Variable | Example | Description |
|----------|---------|-------------|
| `DSN` | `postgres://user:pass@localhost:5432/db?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | `your-secret-key` | Secret for signing JWTs |
| `JWT_TTL` | `24h` | Token expiry duration |

---

## API Reference

> 📖 **Interactive Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/auth/register` | ❌ | Register a new user |
| `POST` | `/auth/login` | ❌ | Login and receive a JWT token |

### Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/users` | ✅ JWT | List all users |

### Messages

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/messages?contact_id=<id>` | ✅ JWT | Get chat history with a user |
| `GET` | `/api/messages/unread` | ✅ JWT | Get unread message counts per sender |
| `POST` | `/api/messages/upload` | ✅ JWT | Upload an audio message (multipart) |

### Contacts

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/contacts` | ✅ JWT | Send a contact request |
| `PATCH` | `/api/contacts/:id` | ✅ JWT | Accept a contact request |
| `GET` | `/api/contacts` | ✅ JWT | List your accepted contacts |
| `DELETE` | `/api/contacts/:id` | ✅ JWT | Remove a contact |

### Admin (role: admin only)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `DELETE` | `/api/admin/users/:id` | ✅ JWT + admin role | Delete a user |

### Real-time

| Protocol | Endpoint | Auth | Description |
|----------|----------|------|-------------|
| WebSocket | `/ws` | ✅ JWT query param | Real-time messaging |

---

## Demo Flow (Postman / cURL)

### Step 1 — Register two users

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"alice1234"}'

curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"bob12345"}'
```

### Step 2 — Login and get JWT tokens

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"alice1234"}' | jq -r .token)

echo "Alice's token: $TOKEN"
```

### Step 3 — Access protected endpoint

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/users
```

### Step 4 — Send a contact request

```bash
# Get Bob's user ID first
BOB_ID=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/users | jq '.users[] | select(.username=="bob") | .id')

curl -X POST http://localhost:8080/api/contacts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"contact_id\": $BOB_ID}"
```

### Step 5 — Login as Bob and accept the request

```bash
BOB_TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"bob12345"}' | jq -r .token)

ALICE_ID=$(curl -s -H "Authorization: Bearer $BOB_TOKEN" \
  http://localhost:8080/api/users | jq '.users[] | select(.username=="alice") | .id')

curl -X PATCH http://localhost:8080/api/contacts/$ALICE_ID \
  -H "Authorization: Bearer $BOB_TOKEN"
```

### Step 6 — Check chat history

```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/messages?contact_id=$BOB_ID"
```

### Step 7 — Or seed everything automatically

```bash
go run ./cmd/seeder
```

---

## Running Tests

```bash
# Run all unit tests
go test ./internal/usecase/... -v

# Run with coverage report
go test ./internal/usecase/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

**Test coverage includes:**
- `TestAuthUseCase_Register` — happy path + duplicate username
- `TestAuthUseCase_Login` — success, wrong password, user not found
- `TestAuthUseCase_GetAllUsers` — returns all registered users
- `TestAuthUseCase_DeleteUser` — removes user from store
- `TestMessageUseCase_SendMessage` — creates message with correct status
- `TestMessageUseCase_GetChatHistory` — bidirectional conversation retrieval
- `TestMessageUseCase_GetChatHistory_Pagination` — limit/offset
- `TestMessageUseCase_UpdateStatus` — status changes (sent → read)
- `TestMessageUseCase_GetUnreadCounts` — counts per sender
- `TestContactUseCase_SendAndAcceptRequest` — request lifecycle
- `TestContactUseCase_GetUserContacts` — lists accepted contacts
- `TestContactUseCase_RemoveContact` — deletes contact relationship
- `TestContactUseCase_DuplicateRequest` — prevents duplicates

---

## Key Features

| Feature | Implementation |
|---------|---------------|
| **JWT Authentication** | `middleware/auth.go` — validates `Authorization: Bearer <token>` |
| **Role-Based Access (RBAC)** | `middleware/role.go` — `admin` group protected by role check |
| **Database Migrations** | `golang-migrate` — versioned SQL files in `/migrations` |
| **Graceful Shutdown** | Hertz `h.Spin()` catches SIGINT/SIGTERM; `defer dbpool.Close()` |
| **Background Worker** | `middleware/ratelimiter.go` — goroutine cleans stale IP entries every minute |
| **Rate Limiting** | Token bucket per IP (20 req/s), returns `429 Too Many Requests` |
| **Request Logging** | Structured log per request: method, path, status, latency, IP |
| **OpenAPI Docs** | `/swagger/index.html` — full interactive API documentation |
| **WebSocket** | Real-time messaging + status updates (sent/delivered/read) |

---

## Database Schema

```sql
-- Users with role support
CREATE TABLE users (
    id         BIGINT PRIMARY KEY,
    username   VARCHAR(24) UNIQUE NOT NULL,
    password   TEXT NOT NULL,
    role       VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Messages with delivery status
CREATE TABLE messages (
    id         SERIAL PRIMARY KEY,
    from_id    BIGINT REFERENCES users(id) ON DELETE CASCADE,
    to_id      BIGINT REFERENCES users(id) ON DELETE CASCADE,
    type       VARCHAR(10) NOT NULL DEFAULT 'text',
    content    TEXT NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'sent',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Many-to-many contacts
CREATE TABLE contacts (
    user_id    BIGINT REFERENCES users(id) ON DELETE CASCADE,
    contact_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status     VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, contact_id)
);
```
