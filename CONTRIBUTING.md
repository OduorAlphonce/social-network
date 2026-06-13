# Contributing

Thanks for helping improve Social Network. This project contains a Go backend,
a React frontend, and design specifications that describe the intended
application.

## Before You Start

1. Follow the [local development setup](docs/setup/local-development.md).
2. Review the [API specification](docs/api/README.md) and
   [database documentation](docs/database/README.md) when your change affects
   a shared contract.
3. Keep changes focused and avoid committing generated dependencies, local
   databases, uploaded files, or `.env` files.

## Development Workflow

Create a branch for the change, make the smallest coherent implementation, and
update relevant tests and documentation alongside the code.

Backend checks:

```bash
cd backend
go test ./...
```

Frontend checks:

```bash
cd frontend
npm run format:check
npm run lint
npm run build
```

## API Changes

The OpenAPI document at
[`docs/api/openapi.json`](docs/api/openapi.json) is the source of truth for
planned request and response contracts. Update it when an API contract changes,
and keep handlers, frontend calls, and examples aligned with it.

The specification includes endpoints that are not implemented yet. Do not
present a specified endpoint as available until it is registered and handled by
the backend.

## Database Changes

Add schema changes as paired migration files under
`backend/internal/db/migrations/`:

```text
000004_short_description.up.sql
000004_short_description.down.sql
```

Migrations run automatically when the backend starts. Keep the conceptual
schema in [`docs/database/schema.dbml`](docs/database/schema.dbml) aligned with
the intended data model.

## Code Style

- Use `gofmt` for Go code.
- Use Prettier for frontend formatting. Run `npm run format` to apply formatting
  and `npm run format:check` to verify it.
- Follow the existing ESLint configuration for frontend code.
- Never commit credentials or real user data.

## Documentation

Keep the root README focused on visitors and the current product state. Put
setup instructions, implementation details, and reference material in `docs/`,
and contributor workflow in this file.
