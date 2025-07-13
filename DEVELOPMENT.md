# Nyla Core Development Guide

This guide covers local development for the Nyla Core analytics engine. For complete project specifications, see [`specs/development.md`](specs/development.md).

---

## Prerequisites
- **Go** 1.24 or later
- **SQLite** 3.39.0 or later
- **Goose** (for migrations): https://github.com/pressly/goose
- **Git** (for versioning)
- **direnv** (for environment variable management): https://direnv.net/

Optional:
- **Docker** (for containerized development/testing)

---

## Environment Setup: direnv

This project uses [direnv](https://direnv.net/) to automatically load environment variables from the `.envrc` file. These variables are required for Goose migrations and other development tasks in the core analytics service.

### Initial Setup

1. **Copy the example environment file:**
   ```bash
   cp .envrc.example .envrc
   ```

2. **Review and customize the environment variables** in `.envrc` according to your local setup.

3. **Trust the direnv configuration:**
   ```bash
   direnv allow
   ```

This will trust the `.envrc` file and ensure all required environment variables are set for your shell session.

If you skip this step, you may see errors from Goose about missing drivers, database strings, or migration files.

### Environment Variables Reference

#### Required Variables (for Goose migrations)
- `GOOSE_DRIVER`: Database driver (default: `sqlite3`)
- `GOOSE_DBSTRING`: Database connection string (default: `./nyla.db`)
- `GOOSE_MIGRATION_DIR`: Migration files directory (default: `./migrations`)

#### Core Server Configuration
- `PORT`: Core server port (default: `8080`)
- `BASE_URL`: Base URL for core server (default: `http://localhost:8080`)

#### CORS Configuration
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed origins (default: `http://localhost:8080,https://localhost`)
- `CORS_ALLOWED_HEADERS`: Allowed request headers for HTMX integration
- `CORS_EXPOSED_HEADERS`: Headers exposed to the browser
- `CORS_ALLOW_CREDENTIALS`: Whether to allow credentials (default: `true`)

#### Development Settings
- `NYLA_ENV`: Environment mode (default: `development`)
- `NYLA_LOG_LEVEL`: Logging level (default: `debug`)

#### Optional Configuration
- `GEOIP_API_KEY`: API key for geo IP service (if using external provider)
- `GEOIP_PROVIDER`: Geo IP provider name
- `NYLA_SITE_ID`: Default site ID for development
- `NYLA_SAMPLING_RATE`: Event sampling rate (1.0 = 100%)

### Updating Environment Variables

After modifying `.envrc`, run:
```bash
direnv allow
```

You can tigger an env reload with `direnv reload` or by simply navigating out of and back into the project directory.

### Configuration Files for Different Environments

For different development scenarios, you may want to create environment-specific configuration files:

- `.envrc.local` - Your personal local overrides (add to .gitignore)
- `.envrc.docker` - Docker-specific configuration
- `.envrc.test` - Test environment configuration

Copy and source the appropriate file as needed:
```bash
cp .envrc.docker .envrc
direnv allow
```

---

## Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/joepurdy/nyla-core.git
   cd nyla-core
   ```

2. **Set up Environment Configuration**
   ```bash
   cp .envrc.example .envrc
   # Review and customize .envrc as needed
   direnv allow
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Initialize the Database**
     ```bash
     make migrate
     make seed
     ```
   - This will create `nyla.db`, run all migrations, and seed with sample data.

---

## Building and Running the Core Server

- **Build the core analytics binary:**
  ```bash
  make nyla-core
  ```
- **Run the core server:**
  ```bash
  ./bin/nyla-core
  ```
  This starts the unified analytics server with both API and dashboard interface.

---

## Database Management

- **Run migrations:**
  ```bash
  make migrate
  ```
- **Reset all migrations:**
  ```bash
  make migrate-reset
  ```
- **Check migration status:**
  ```bash
  make migrate-status
  ```
- **Seed the database:**
  ```bash
  make seed
  ```

---

## Troubleshooting

- **Missing Goose:**
  - Install with: `go install github.com/pressly/goose/v3/cmd/goose@latest`
- **Environment variable errors:**
  - Ensure you've copied `.envrc.example` to `.envrc` and run `direnv allow`
  - Check that all required environment variables are set: `env | grep GOOSE`
  - If direnv isn't working, manually source the file: `source .envrc`
- **Database file issues:**
  - Ensure you have write permissions in the `nyla` directory.
  - Delete `nyla.db` and re-run `make migrate` if migrations fail.
- **Build errors:**
  - Check Go version (`go version`).
  - Run `go mod tidy` to clean up dependencies.
- **Port conflicts:**
  - If the core server fails to start, ensure port 8080 (or your configured port) is free.
  - Modify the PORT in `.envrc` if needed and run `direnv reload`.
- **CORS issues:**
  - Update `CORS_ALLOWED_ORIGINS` in `.envrc` to include your development URLs.
  - Ensure the core server is running for local development.
- **Seeding errors:**
  - Ensure `data/dump.data` exists and is formatted correctly.

---

## Reference
- Main project development guide: [`/specs/development.md`](specs/development.md)
- For issues not covered here, check the main guide or open an issue. 