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
   - Real-time capabilities
   - Optimized queries
   - Small client footprint
   - Caching strategies

4. **Self-Hosted**
   - Easy installation
   - Simple backup process
   - Minimal dependencies
   - Clear upgrade path
   - Resource efficient

## Implementation Status

The specifications are currently in development, with the following status:

- ✅ Core architecture decisions
- ✅ API design patterns
- ✅ Database schema
- ✅ Deployment strategy
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