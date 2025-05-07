# Nyla API Development Guide

This guide covers local development for the Nyla API service. For full project setup, see the main development spec at [`specs/development.md`](../specs/development.md).

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

This project uses [direnv](https://direnv.net/) to automatically load environment variables from the `.envrc` file in the **`api/` directory**. These variables are required for Goose migrations and other development tasks in the API service.

**After cloning the repository, you must change into the `api/` directory and run:**
```bash
cd api
# (or ensure you are in the api directory)
direnv allow
```
This will trust the `.envrc` file in `api/` and ensure all required environment variables (such as `GOOSE_DRIVER` and `GOOSE_DBSTRING`) are set for your shell session.

If you skip this step, you may see errors from Goose about missing drivers, database strings, or migration files when working in the `api/` directory.

---

## Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/joepurdy/nyla.git
   cd nyla/api
   ```

2. **Trust Direnv Config**
   - From the `api/` directory:
   ```bash
   direnv allow
   ```

2. **Install Dependencies**
   - From the `api/` directory:
   ```bash
   go mod download
   ```

3. **Initialize the Database**
   - From the `api/` directory:
     ```bash
     make migrate
     make seed
     ```
   - This will create `nyla.db`, run all migrations, and seed with sample data.

---

## Building and Running the API

- **Build the API binary:**
  - From the `api/` directory:
  ```bash
  make nyla-api
  ```
- **Run the API:**
  - From the `api/` directory:
  ```bash
  ./nyla-api
  ```

---

## Database Management

- **Run migrations:**
  - From the `api/` directory:
  ```bash
  make migrate
  ```
- **Reset all migrations:**
  - From the `api/` directory:
  ```bash
  make migrate-reset
  ```
- **Check migration status:**
  - From the `api/` directory:
  ```bash
  make migrate-status
  ```
- **Seed the database:**
  - From the `api/` directory:
  ```bash
  make seed
  ```

---

## Troubleshooting

- **Missing Goose:**
  - Install with: `go install github.com/pressly/goose/v3/cmd/goose@latest`
- **Database file issues:**
  - Ensure you have write permissions in the `api/` directory.
  - Delete `nyla.db` and re-run `make migrate` if migrations fail.
- **Build errors:**
  - Check Go version (`go version`).
  - Run `go mod tidy` to clean up dependencies.
- **Port conflicts:**
  - If the API fails to start, ensure port 9876 (or your configured port) is free.
- **Seeding errors:**
  - Ensure `data/dump.data` exists and is formatted correctly.

---

## Reference
- Main project development guide: [`../specs/development.md`](../specs/development.md)
- For issues not covered here, check the main guide or open an issue. 