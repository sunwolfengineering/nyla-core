# nyla
GDPR compliant privacy focused web analytics

## Overview

Nyla is a self-hosted web analytics platform that prioritizes user privacy and data protection. Built with simplicity and performance in mind, it provides essential analytics capabilities without compromising user privacy or website performance.

## Key Features

- Privacy-first design with GDPR/CCPA compliance
- Lightweight JavaScript tracker (<5KB gzipped)
- Real-time analytics via server-sent events
- Simple self-hosted deployment
- No cookies required
- IP anonymization by default
- Configurable data retention
- Server-side rendering with progressive enhancement

## Technical Stack

- **Backend**: Go with SQLite
- **Frontend**: HTMX + Hyperscript
- **Styling**: TailwindUI components
- **Real-time**: Server-sent events
- **Deployment**: Single binary

## Technical Specifications

The complete technical specifications for Nyla can be found in the [SPECS.md](SPECS.md) file. These specifications cover:

- [Architecture Overview](specs/architecture-overview.md)
- [API Specification](specs/api-specification.md)
- [JavaScript Tracker](specs/js-tracker-specification.md)
- [Database Schema](specs/database-schema.md)
- [Deployment Guide](specs/deployment.md)
- [Development Guide](specs/development.md)

## Development

See the [Development Guide](specs/development.md) for detailed setup instructions.

### Quick Start

```bash
# Clone repository
git clone https://github.com/joepurdy/nyla.git
cd nyla

# Initialize development database
make init-db

# Start development server
make dev
```

### Requirements

- Go 1.24 or later
- Node.js 22 LTS or later (build tools only)
- SQLite 3.39.0 or later
- Docker (optional, for container testing)

## Implementation Status

The project is currently in active development:

- ✅ Core architecture and technical decisions
- ✅ Development workflow and tooling
- 🚧 Initial implementation
- 📝 Documentation and examples