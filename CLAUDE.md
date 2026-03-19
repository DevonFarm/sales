# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Devon Farm Sales is a Go web application for managing horse farms. Users authenticate via Stytch (passwordless magic links), manage farms, and perform CRUD operations on horses. Built with Fiber, CockroachDB, and server-rendered HTML templates.

## Commands

Uses [go-task](https://taskfile.dev/) as the task runner. Environment variables are loaded from `.env` (see `.env.example`).

```bash
# Setup
task setup                    # Install all dependencies (Homebrew, CockroachDB, go-migrate, DB cert)

# Run
go run .                      # Start server on :4242 (requires COCKROACH_DSN)

# Tests
task test-unit                # Unit tests (./utils/... ./horse/... ./user/... ./farm/...)
task test-handlers            # HTTP handler tests (./tests/handlers/...)
task test-all                 # All Go tests (unit + handlers)
task test-e2e                 # Playwright E2E tests (cd tests/e2e && npm test)

# Run a single Go test
go test -v -run TestFuncName ./package/...

# Database
task db-migrate               # Run migrations (requires MIGRATIONS_DSN)
task db-create-migration -- name  # Create new migration
task sql                      # SQL shell (requires COCKROACH_DSN)
task clean-db:dbname          # Drop and recreate a test database
```

## Architecture

**Entry point**: `main.go` — embeds `templates/` and `assets/` via `//go:embed`, creates a server, registers routes, and listens on `:4242`.

**Modular by entity** — each domain package (`horse/`, `farm/`, `user/`) owns its model, database queries, HTTP handlers, and route registration via `RegisterRoutes()`.

**Key packages**:

- `server/` — Fiber app initialization, middleware (logger), template engine setup, static asset serving
- `auth/` — Stytch integration: magic link login/callback/logout, session middleware (`RequireAuth()`)
- `database/` — pgx connection wrapper; migrations in `database/migrations/` (uses go-migrate)
- `templates/` — HTML templates with `templates/layouts/main.html` as the shared layout
- `utils/` — Date/age helpers, image handling, HTML path utilities

**Database**: CockroachDB (PostgreSQL-compatible) with three tables: `horses`, `farms`, `users`. Horses belong to farms via `farm_id` FK.

**Auth flow**: Stytch passwordless magic links. `auth.StytchAuth` registers its own routes and provides `RequireAuth()` middleware for protected routes.

## Testing

Three test layers:

1. **Unit tests** — colocated `_test.go` files in each package. Pure logic tests (age calculation, validation, etc.)
2. **Handler tests** — `tests/handlers/` — test HTTP endpoints using mock database implementations and Fiber's `httptest`
3. **E2E tests** — `tests/e2e/` — Playwright (TypeScript), run with `npm test` from that directory

Test utilities live in `tests/util/` (database helpers, auth mocking, HTTP assertion helpers).
