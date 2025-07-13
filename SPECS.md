# Nyla Analytics Core - Technical Specifications

## Overview

Nyla Analytics Core is the open-source foundation of a privacy-focused web analytics platform. These specifications define the core analytics engine, event collection system, and minimal self-hosted UI that forms the basis for both self-hosted deployments and the commercial SaaS offering.

## Open-Core Architecture

Nyla Core is a complete, self-hosted web analytics platform:

- **Event Collection**: Real-time analytics tracking
- **Analytics Engine**: Core metrics and insights
- **Dashboard Interface**: Simple, focused analytics view
- **Self-Hosted Deployment**: Single-binary installation

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
- Event collection endpoints
- Analytics data retrieval
- Authentication and access control
- Rate limiting and privacy controls
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
- SQLite schema for analytics data
- Event storage and session tracking
- Privacy controls and data retention
- Aggregation tables and performance optimization
- Backup and maintenance procedures

### [Deployment](specs/deployment.md)
- Self-hosted container deployment
- Single-binary architecture
- Environment configuration
- Resource requirements and sizing
- Monitoring and logging setup
- Backup procedures
- Security considerations

### [Development Guide](specs/development.md)
- Development environment setup
- Project structure and organization
- Local development workflow
- Git workflow and conventions
- Testing strategy and patterns
- Build pipeline and tooling
- Debugging guidelines
- Contribution guidelines

## Domain Configuration

For self-hosted deployments, Core supports:
- Custom domain configuration
- Local analytics interface
- API endpoint customization

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
   - Minimal core interface
   - Essential analytics only
   - Zero-config installation

3. **Performance**
   - Efficient event collection
   - Basic real-time updates
   - Optimized core queries
   - Small tracker footprint (<5KB)
   - Resource efficient

4. **Self-Hosted First**
   - One-command deployment
   - File-based configuration
   - Single dependency (SQLite)
   - Simple upgrade process
   - Minimal resource usage

5. **Community-Focused**
   - Stable public APIs
   - Extensible architecture
   - Clear separation of concerns
   - Open development process
   - Apache 2.0 licensing

## Technical Stack

1. **Backend**
   - Go 1.24+ for analytics engine
   - SQLite for data storage
   - REST API for event collection
   - HTML interface

2. **Frontend**
   - Minimal JavaScript (<5KB tracker)
   - Server-rendered HTML
   - Real-time updates
   - Analytics dashboard

3. **Build & Deployment**
   - Single Go binary
   - Embedded static assets
   - Container-first deployment
   - Zero external dependencies

## Implementation Status

Components under development:

- ✅ Architecture design
- ✅ API patterns
- ✅ Database schema
- ✅ Deployment strategy
- ✅ Development workflow
- 🚧 Analytics engine
- 🚧 Dashboard interface
- 🚧 JavaScript tracker
- 📝 Documentation

## Core Features

- Event collection and storage
- Pageview and session analytics
- Real-time dashboard
- Self-hosted deployment
- API key authentication
- Data export (JSON/CSV)
- Privacy controls and GDPR compliance
- Single-binary installation

## Licensing

- Apache 2.0 (permissive open source)

## Contributing

To contribute to the specifications:

1. Review existing specifications
2. Discuss changes in GitHub issues
3. Submit pull requests with clear descriptions
4. Update implementation status and tests
5. Follow conventional commit guidelines 