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
   - Dual licensing (GPL v3 + Commercial)

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

Nyla Analytics Core is dual-licensed under:

### Open Source License
**GNU General Public License v3 (GPL v3)** - for open source projects and community use. This copyleft license ensures that modifications remain open source, promoting transparency and community contributions.

### Commercial License
**Commercial License** - for businesses that want to integrate Nyla Core into proprietary applications without GPL obligations. Contact us for commercial licensing terms and pricing.

This dual licensing model allows:
- **Open source projects**: Free use under GPL v3
- **Commercial projects**: Proprietary use with commercial license
- **SaaS providers**: Flexibility to choose appropriate license

## Contributing

To contribute to Nyla Analytics Core:

### 📋 Before Contributing
1. Review the [Contributor License Agreement (CLA)](CLA.md)
2. Sign the CLA in your first pull request
3. Review existing specifications and code

### 🔄 Contribution Process
1. Discuss significant changes in GitHub issues first
2. Fork the repository and create a feature branch
3. Make your changes with clear, tested code
4. Submit a pull request with detailed description
5. Include CLA acceptance in your PR
6. Respond to code review feedback

### 📚 Types of Contributions
- **Code**: Core analytics features, bug fixes, performance improvements
- **Documentation**: Specification updates, guides, examples
- **Testing**: Unit tests, integration tests, performance benchmarks
- **Security**: Security fixes, vulnerability reports

### ⚖️ Licensing Note
All contributions are subject to the project's dual licensing model. By contributing, you grant rights for both GPL v3 and commercial licensing as outlined in the [CLA](CLA.md). 