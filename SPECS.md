# Nyla Analytics Technical Specifications

## Overview

Nyla Analytics is a privacy-focused, self-hosted web analytics platform designed for simplicity and performance. These specifications outline the technical architecture, implementation details, and operational considerations for the platform.

## Core Specifications

### [Architecture Overview](specs/architecture-overview.md)
- System components and high-level design
- Technical stack decisions and rationale
- Infrastructure approach
- Development workflow
- Security and privacy considerations
- System architecture diagrams
- Component interactions
- Data flow patterns

### [API Specification](specs/api-specification.md)
- Hypermedia-driven interface design
- Event collection endpoints
- Real-time analytics interface
- Server-sent events integration
- Error handling and rate limiting
- Authentication and authorization
- Progressive enhancement strategy

### [JavaScript Tracker](specs/js-tracker-specification.md)
- Lightweight client implementation (<5KB gzipped)
- Installation and configuration
- Automatic and custom event tracking
- Privacy features and compliance
- Browser compatibility
- Performance optimization
- Integration guides

### [Database Schema](specs/database-schema.md)
- SQLite table structure and indexes
- Data models and relationships
- Privacy and retention policies
- Query optimization
- Maintenance procedures
- Backup configuration

### [Deployment](specs/deployment.md)
- Container configuration
- Environment setup
- Resource requirements
- Monitoring and logging
- Backup procedures
- Maintenance tasks
- Security hardening

### [Development Guide](specs/development.md)
- Environment setup
- Required tools and versions
- Local development workflow
- Git workflow and conventions
- Testing strategy
- CI/CD pipeline
- Debugging guidelines
- Code organization

## Domain Configuration

The [Domain Conventions](.cursor/rules/domain-conventions.mdc) rule defines the domain structure for all Nyla services:

- `getnyla.app` - Root domain
- `app.getnyla.app` - Main application interface
- `dashboard.getnyla.app` - Analytics dashboard
- `api.getnyla.app` - API endpoints
- `cdn.getnyla.app` - Content delivery network

## Key Design Principles

1. **Privacy-First**
   - IP anonymization by default
   - Minimal data collection
   - Configurable retention
   - GDPR/CCPA compliance
   - No cross-site tracking
   - No cookies required

2. **Simplicity**
   - Single binary deployment
   - SQLite for storage
   - Hypermedia-driven interface
   - Minimal JavaScript
   - Progressive enhancement

3. **Performance**
   - Efficient data storage
   - Real-time capabilities via SSE
   - Optimized queries
   - Small client footprint (<5KB)
   - Caching strategies

4. **Self-Hosted**
   - Easy installation
   - Simple backup process
   - Minimal dependencies
   - Clear upgrade path
   - Resource efficient

## Technical Stack

1. **Backend**
   - Go for server implementation
   - SQLite for data storage
   - HTML-over-the-wire with HTMX
   - Server-sent events for real-time

2. **Frontend**
   - HTMX + Hyperscript
   - TailwindUI components
   - Server-side state management
   - Progressive enhancement
   - Dark/light themes

3. **Build Tools**
   - esbuild for JS
   - PostCSS for Tailwind
   - go:embed for assets
   - GitHub Actions for CI/CD

## Implementation Status

The specifications are currently in development, with the following status:

- ✅ Core architecture decisions
- ✅ API design patterns
- ✅ Database schema
- ✅ Deployment strategy
- ✅ Development workflow
- ✅ Technical stack decisions
- 🚧 HTML templates and components
- 🚧 Integration examples
- 📝 Documentation

## Future Considerations

Areas planned for future specification:

1. **Multi-Site Support**
   - Site isolation
   - Resource quotas
   - Cross-site analytics

2. **Team Management**
   - Role-based access
   - Audit logging
   - Team permissions

3. **Advanced Analytics**
   - Custom dashboards
   - Data export
   - Advanced filtering
   - Custom metrics

4. **Integration Ecosystem**
   - CMS plugins
   - Framework integrations
   - Export/import tools

## Contributing

To contribute to these specifications:

1. Review the existing specifications
2. Discuss major changes in issues first
3. Submit pull requests with specification updates
4. Ensure cross-reference consistency
5. Update implementation status 