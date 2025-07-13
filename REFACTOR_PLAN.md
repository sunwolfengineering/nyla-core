# Nyla Analytics Core - Refactor Plan

**Target**: Transform current codebase into GPL v3 dual-licensed open-core analytics engine with single-site support and integrated dashboard.

## 🎯 Executive Summary

**Goal**: Refactor Nyla into a single-binary, single-site analytics core that:
- Combines API server and dashboard into one process
- Enforces single site (`site_id = "default"`) throughout
- Implements clean extension points for commercial features
- Provides essential analytics: pageviews, sessions, real-time visitors
- Uses SQLite with proper migrations and schema from specs
- Achieves 80%+ test coverage with full CI/CD

## 📊 Current State Analysis

### ❌ Problems to Fix
- **Two binaries**: `cmd/api` + `cmd/ui` violates single-binary requirement
- **Multi-site leakage**: Code accepts arbitrary `site_id` values
- **Outdated schema**: Database doesn't match `specs/database-schema.md`
- **Incomplete dashboard**: Only basic real-time card implemented
- **No extensions**: No plugin interfaces for commercial features
- **License mismatch**: AGPL headers instead of GPL v3
- **Minimal testing**: Low coverage, no CI pipeline

### ✅ Current Strengths
- Go backend with SQLite foundation
- Basic event collection working
- elem-go templating system functional
- Repository structure established

## 🏗️ Target Architecture

```
nyla-core/
├── cmd/nyla-core/           # Single binary entry point
├── internal/
│   ├── server/              # HTTP routing, middleware, static assets
│   ├── storage/             # SQLite wrapper, migrations, queries
│   ├── services/            # Business logic (events, stats, sessions)
│   ├── dashboard/           # HTML rendering, template helpers
│   └── plugins/             # Extension registry and feature flags
├── pkg/core/                # Public interfaces for commercial extensions
├── web/                     # Embedded templates and static assets
├── migrations/              # Numbered SQL migration files
└── tests/                   # Unit, integration, and e2e tests
```

## 📋 Refactor Phases

### Phase 1: Repository Restructuring (Week 1)
**Goal**: Reorganize codebase and remove multi-site assumptions

#### Tasks:
1. **Merge binaries**
   - Delete `cmd/api/` and `cmd/ui/`
   - Create `cmd/nyla-core/main.go` as single entry point
   - Move HTTP server logic to `internal/server/`

2. **Enforce single-site**
   - Add constant `const DefaultSiteID = "default"` in central location
   - Remove multi-site parameters from all functions
   - Add validation: reject any `site_id != "default"` with HTTP 400

3. **Package reorganization**
   - Move business logic to `internal/services/`
   - Create `internal/storage/` for database operations
   - Move UI logic to `internal/dashboard/`
   - Create `pkg/core/` for public extension interfaces

#### Acceptance Criteria:
- [ ] Single `cmd/nyla-core/main.go` builds successfully
- [ ] All code references single default site
- [ ] Package structure matches target architecture
- [ ] No multi-site logic remains in codebase

### Phase 2: Database Layer (Week 1-2)
**Goal**: Implement proper SQLite schema with migrations

#### Tasks:
1. **Schema implementation**
   - Create `migrations/001_initial_schema.sql` from `specs/database-schema.md`
   - Implement automatic migration runner in `internal/storage/migrate.go`
   - Add SQLite PRAGMA settings (WAL mode, foreign keys, etc.)

2. **Storage layer**
   - Create `internal/storage/db.go` with connection management
   - Implement typed methods: `InsertEvent()`, `GetRealtimeStats()`, etc.
   - Add proper error handling and connection pooling

3. **Data consistency**
   - Ensure default site configuration exists on startup
   - Add database integrity checks
   - Implement backup/restore utilities

#### Acceptance Criteria:
- [ ] Database schema matches specifications exactly
- [ ] Migrations run automatically on startup
- [ ] All storage operations use typed methods
- [ ] SQLite configured with optimal settings for analytics workload

### Phase 3: Core Services (Week 2-3)
**Goal**: Implement business logic for analytics features

#### Tasks:
1. **Event service** (`internal/services/events/`)
   - Event validation and sanitization
   - IP anonymization (configurable)
   - Batch processing for high-volume inserts
   - Session detection and management

2. **Statistics service** (`internal/services/stats/`)
   - Real-time visitor queries (30-minute window)
   - Historical pageview aggregations
   - Top pages and referrers
   - Session duration calculations

3. **Settings service** (`internal/services/settings/`)
   - Site configuration management
   - Privacy settings (retention, anonymization)
   - Basic customization options

#### Acceptance Criteria:
- [ ] Events processed with proper validation
- [ ] Real-time statistics calculated accurately
- [ ] Privacy controls functional (IP anonymization, retention)
- [ ] All services have comprehensive unit tests

### Phase 4: Integrated Dashboard (Week 3)
**Goal**: Embed dashboard into main server process

#### Tasks:
1. **Template system**
   - Convert elem-go templates to `html/template` format
   - Embed templates using `go:embed` in `web/templates/`
   - Add template hot-reloading for development

2. **Dashboard components**
   - Real-time visitor count card
   - Pageviews chart (24h, 7d, 30d views)
   - Top pages table with sorting
   - Basic settings page

3. **Real-time updates**
   - Server-Sent Events (SSE) endpoint at `/api/updates`
   - WebSocket fallback for older browsers
   - Automatic reconnection on connection loss

#### Acceptance Criteria:
- [ ] Dashboard accessible at `/dashboard` endpoint
- [ ] All essential metrics displayed correctly
- [ ] Real-time updates working via SSE
- [ ] Mobile-responsive design
- [ ] Settings page functional

### Phase 5: Extension Points (Week 4)
**Goal**: Create clean interfaces for commercial features

#### Tasks:
1. **Plugin architecture**
   - Define `Plugin` interface in `pkg/core/plugin.go`
   - Create plugin registry in `internal/plugins/registry.go`
   - Add plugin lifecycle management (init, start, stop)

2. **Extension hooks**
   - `AfterEventStored(event *Event)` - for custom processing
   - `DashboardWidgets() []Widget` - for custom dashboard components
   - `CustomMetrics() []Metric` - for additional analytics

3. **Feature flags**
   - Environment-based feature toggles
   - Build-tag separation for commercial features
   - Runtime detection of available features

#### Acceptance Criteria:
- [ ] Plugin interface defined and documented
- [ ] Extension points working with mock plugins
- [ ] Feature flags controllable via environment variables
- [ ] Commercial build can extend without modifying core

### Phase 6: Licensing Compliance (Week 4)
**Goal**: Ensure GPL v3 compliance throughout codebase

#### Tasks:
1. **License headers**
   - Replace all AGPL headers with GPL v3
   - Add SPDX identifiers: `// SPDX-License-Identifier: GPL-3.0-only`
   - Ensure consistency across all Go files

2. **Legal compliance**
   - Update copyright notices
   - Add license notice to binary output
   - Create license scanning CI job

3. **Documentation updates**
   - Update all references to GPL v3 in docs
   - Add commercial licensing information
   - Create contributor guidelines with CLA

#### Acceptance Criteria:
- [ ] All source files have correct GPL v3 headers
- [ ] License scanner passes in CI
- [ ] Binary displays correct license information
- [ ] Documentation reflects dual licensing model

### Phase 7: Testing & CI (Week 4-5)
**Goal**: Achieve comprehensive test coverage and automated validation

#### Tasks:
1. **Unit tests**
   - Event validation and processing
   - Statistics calculations
   - Database operations
   - Template rendering
   - Target: 80%+ coverage

2. **Integration tests**
   - End-to-end event collection
   - Dashboard rendering
   - SSE functionality
   - Database migrations

3. **CI/CD pipeline**
   - GitHub Actions workflow
   - Automated testing on multiple Go versions
   - License compliance checking
   - Binary builds for multiple platforms

#### Acceptance Criteria:
- [ ] 80%+ test coverage across all packages
- [ ] Integration tests verify core workflows
- [ ] CI pipeline passes on all supported platforms
- [ ] Automated releases with proper versioning

## 🚀 Implementation Guidelines

### Code Organization Principles
1. **Dependency Injection**: Use interfaces for all external dependencies
2. **Single Responsibility**: Each package has one clear purpose
3. **Testability**: All business logic easily unit testable
4. **Extension Points**: Well-defined interfaces for commercial features

### Database Guidelines
1. **Migrations Only**: Never modify schema directly
2. **Backwards Compatibility**: Support upgrade paths
3. **Performance**: Proper indexing for analytics queries
4. **Integrity**: Foreign key constraints enforced

### API Design
1. **RESTful**: Follow REST conventions for data endpoints
2. **Hypermedia**: Dashboard uses HTMX for interactivity
3. **Versioning**: All API endpoints versioned (`/v1/`)
4. **Documentation**: OpenAPI specs for all endpoints

### Security Requirements
1. **Input Validation**: All user inputs validated and sanitized
2. **Rate Limiting**: Basic protection against abuse
3. **Privacy**: IP anonymization enabled by default
4. **CORS**: Properly configured for cross-origin requests

## 📚 Testing Strategy

### Unit Tests (70% of coverage)
- Individual function behavior
- Business logic validation
- Error handling
- Edge cases

### Integration Tests (20% of coverage)
- Database operations
- HTTP endpoint behavior
- Template rendering
- Real-time updates

### End-to-End Tests (10% of coverage)
- Full user workflows
- Performance under load
- Browser compatibility
- Mobile responsiveness

## 🎉 Success Criteria

### Functional Requirements
- [ ] Single binary deployment
- [ ] Single site analytics working
- [ ] Basic dashboard functional
- [ ] Real-time updates operational
- [ ] Privacy controls working

### Non-Functional Requirements
- [ ] 80%+ test coverage
- [ ] GPL v3 compliance verified
- [ ] Performance: 1000+ events/second
- [ ] Memory usage: <100MB typical
- [ ] Build time: <30 seconds

### Business Requirements
- [ ] Extension points for commercial features
- [ ] Clear upgrade path to commercial version
- [ ] Documentation for developers
- [ ] Easy self-hosting process

## 📖 Documentation Updates Required

1. **README.md**: Update installation and usage instructions
2. **DEVELOPMENT.md**: New build and test procedures
3. **SPECS.md**: Architecture changes and extension points
4. **API.md**: Complete endpoint documentation
5. **MIGRATION.md**: Upgrade guide from previous versions

---

**Next Steps**: Begin Phase 1 by creating the new directory structure and merging the binary entry points. Each phase should be completed in order, with acceptance criteria verified before proceeding.

This plan provides a comprehensive roadmap for transforming Nyla into a clean, testable, and extensible open-core analytics platform while maintaining the essential simplicity that makes it attractive for self-hosting.
