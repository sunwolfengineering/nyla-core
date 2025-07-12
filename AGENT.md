# Nyla Development Guide for AI Agents

## Commands
- Build: `make nyla-api` (API) / `make nyla-ui` (UI)
- Test: `go test ./...` (no Makefile target yet - tests not implemented)
- Migrations: `make migrate` / `make migrate-reset` / `make migrate-status`
- Database: `make seed` (reset DB and import sample data)
- Dependencies: `go mod download` / `go mod tidy`

## Architecture
- **Stack**: Go 1.24+ backend, SQLite database, HTMX+Hyperscript frontend
- **Structure**: `cmd/` (binaries), `pkg/` (packages: db, geo, handlers, hash), `migrations/` (Goose SQL), `data/` (SQLite DB), `js-collector/` (tracker), `specs/` (documentation)
- **Services**: API server (`cmd/api`) and UI server (`cmd/ui`) as separate binaries
- **Database**: SQLite with Goose migrations, uses direnv for env vars

## Code Style & Conventions
- **Domain**: Always use `getnyla.app` (root), subdomains: `app.`, `dashboard.`, `api.`, `cdn.`
- **Go Style**: Standard Go formatting (gofmt), use testify for tests, table-driven test patterns
- **Privacy First**: GDPR compliant, no cookies, IP anonymization, <5KB JS tracker
- **Commits**: Conventional commits format: `type(scope): description`
  - **Types**: `feat` (features), `fix` (bug fixes), `docs` (documentation), `refactor` (code refactoring), `test` (tests), `chore` (maintenance)
  - **Scopes**: `db` (database/migrations), `api` (REST API), `ui` (user interface), `js` (JavaScript tracker), `deploy` (deployment), `build` (build system)
  - **Examples**: `feat(db): add sites table migration`, `fix(api): correct event collection validation`, `docs(readme): update installation guide`
- **Dependencies**: Minimal external deps - elem-go for HTML, modernc.org/sqlite, useragent parsing

## Development Rules
- **NO implementation without specifications** - check `specs/` directory first, request updates if work is unspecified
- **Testing**: Target 80% coverage, use testify framework, unit + integration + e2e testing required
- **Deployment**: Single static binary with embedded assets, SQLite only external dependency
- Environment setup requires `direnv allow` for Goose migrations to work
