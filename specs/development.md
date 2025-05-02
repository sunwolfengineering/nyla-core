# Development Guide

## Environment Setup

### Required Tools
- Go 1.24 or later
- Node.js 22 LTS or later (for build tools only)
- SQLite 3.39.0 or later
- Docker (optional, for container testing)

### Local Development

1. **Clone Repository**
   ```bash
   git clone https://github.com/joepurdy/nyla.git
   cd nyla
   ```

2. **Initialize Development Database**
   ```bash
   make init-db
   ```
   This will:
   - Create SQLite database in `data/nyla.db`
   - Run all migrations
   - Load sample data for testing

3. **Start Development Server**
   ```bash
   make dev
   ```
   This starts:
   - API server on http://localhost:3000
   - Live reload for template changes
   - Development database
   - Server-sent events endpoint

### Test Data
Development environment includes sample data for:
- Page views
- Custom events
- Multiple sites
- Various time periods

## Git Workflow

### Branch Naming
- Feature branches: `feature/description`
- Bug fixes: `fix/description`
- Documentation: `docs/description`
- Linear tickets: Use auto-generated branch names

### Commit Messages
Follow conventional commits:
```
type(scope): description

[optional body]

[optional footer]
```

Types:
- feat: New feature
- fix: Bug fix
- docs: Documentation
- chore: Maintenance
- test: Test updates
- refactor: Code restructuring

### Pull Request Process
1. Create branch from main
2. Implement changes
3. Ensure tests pass locally
4. Submit PR with:
   - Clear description
   - Link to Linear ticket
   - Test coverage
   - Migration scripts if needed
5. Address review feedback
6. Squash and merge

## CI Pipeline

CI/CD implementation is planned for a future ticket. The pipeline will include:

### Quality Checks
- Go test coverage (minimum 80%)
- Go linting with golangci-lint
- SQL migrations validation
- Template syntax checking
- Security scanning with gosec

## Testing Strategy

### Unit Testing
- Go tests for business logic
- Table-driven tests for data processing
- Mocked external dependencies
- Coverage requirements enforced in CI

### Integration Testing
- API endpoint testing
- Database operations
- Event processing pipeline
- Real-time update system

### End-to-End Testing
- Core user journeys
- Dashboard functionality
- Data collection flow
- Export/import operations

### Performance Testing
- Request latency benchmarks
- Database query optimization
- Memory usage monitoring
- Load testing critical paths

## Development Commands

```bash
# Start development server
make dev

# Run tests
make test

# Run linters
make lint

# Build binary
make build

# Build container
make container

# Generate mocks
make generate

# Clean build artifacts
make clean
```

## Debugging

### Local Development
- API logs to stdout/stderr
- SQLite database in `data/nyla.db`
- Template cache in `tmp/templates`
- Debug logs with `LOG_LEVEL=debug`

### Production Issues
- Health check endpoint: `/health`
- Metrics endpoint: `/metrics`
- Structured logging with correlation IDs
- Error reporting via sentry (optional)

## Code Organization

```
nyla/
├── api/           # API handlers and middleware
├── cmd/           # Command line tools
├── internal/      # Internal packages
├── migrations/    # Database migrations
├── static/        # Static assets
├── templates/     # HTML templates
└── web/           # Web interface handlers
``` 