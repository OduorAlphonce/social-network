# Local Development

## Prerequisites

- Go 1.22 or newer
- Node.js 20.19 or newer, or 22.12 or newer
- npm
- A C compiler and SQLite development libraries for `go-sqlite3`

## Backend

From the repository root, create the local environment file:

```bash
cp backend/.env.example backend/.env
```

Replace the placeholder paths in `backend/.env` with working values:

```ini
PORT=8080
DATABASE_PATH=./db.sqlite
APP_ENV=development
ALLOWED_ORIGIN=http://localhost:5173
MIGRATIONS_DIR=./internal/db/migrations
```

Run the API from `backend/` so these relative paths resolve correctly:

```bash
cd backend
go run ./cmd/server
```

The server listens on `http://localhost:8080`. Database migrations are applied
automatically at startup, and the SQLite database is created at
`backend/db.sqlite`.

## Frontend

In a second terminal:

```bash
cd frontend
npm install
npm run dev
```

Vite serves the frontend at `http://localhost:5173` by default.

The current frontend is primarily a UI scaffold. As API integration is added,
requests should target `http://localhost:8080/api` and include credentials for
authenticated routes so the session cookie is sent.

## Validation

Run backend checks:

```bash
cd backend
go test ./...
```

Run frontend checks:

```bash
cd frontend
npm run lint
npm run build
```

## Docker Status

The repository contains a backend `Dockerfile`, but it is not currently a
complete supported full-stack workflow: there is no frontend image or Compose
configuration, and the backend image's migration path still needs to be aligned
with `backend/internal/db/migrations/`. Use the local setup above for now.
