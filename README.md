# Nyla Analytics Core
GDPR compliant privacy focused web analytics - Open Source Edition

## Overview

Nyla Analytics Core is the open-source foundation of a privacy-first web analytics platform. It provides essential analytics capabilities for self-hosted deployments, focusing on simplicity, performance, and user privacy without compromising website performance.

This is the core open-source version that powers single-site analytics. For multi-site management and advanced features, see our commercial offering.

## Core Features

- **Privacy-first design** with GDPR/CCPA compliance
- **Lightweight JavaScript tracker** (<5KB gzipped)
- **Single-site analytics** - perfect for personal sites and small projects
- **Real-time visitor tracking** with basic dashboard
- **Simple self-hosted deployment** - single binary, SQLite database
- **No cookies required** - respects user privacy by default
- **IP anonymization** by default
- **Configurable data retention** policies
- **Essential metrics** - pageviews, visitors, top pages, referrers
- **Dual licensed** - GPL v3 (open source) + Commercial licensing available

## Technical Stack

- **Backend**: Go 1.24+ with SQLite database
- **Frontend**: Server-rendered HTML with basic JavaScript
- **Deployment**: Single binary with embedded assets
- **Real-time**: Basic live visitor updates
- **Storage**: SQLite for zero-config setup

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
git clone https://github.com/joepurdy/nyla-core.git
cd nyla-core

# Set up environment
cp .envrc.example .envrc
direnv allow

# Initialize development database
make migrate
make seed

# Build and run
make nyla-core
./bin/nyla-core
```

### Requirements

- **Go 1.24+** - Core application development
- **SQLite 3.39.0+** - Database storage
- **direnv** - Environment variable management
- **Docker** (optional) - Container testing and deployment

## Implementation Status

Nyla Analytics Core is in active development:

- ✅ Open-source core architecture designed
- ✅ Core specifications and development workflow
- ✅ Database schema and API design
- 🚧 Core analytics engine implementation
- 🚧 Basic dashboard interface
- 🚧 JavaScript tracker development
- 📝 Documentation and deployment guides

## License

Nyla Analytics Core is **dual-licensed**:

### 📖 Open Source License
Licensed under the **GNU General Public License v3 (GPL v3)** for open source use - see the [LICENSE](LICENSE) file for details.

### 💼 Commercial License
Commercial licensing available for businesses that want to integrate Nyla Core into proprietary applications without GPL v3 obligations.

**Choose the license that fits your needs:**
- **Open source projects** → Use GPL v3 (free)
- **Commercial/proprietary projects** → Contact us for commercial licensing
- **SaaS providers** → Either license works, depending on your requirements

📄 **Learn more**: See [LICENSING.md](LICENSING.md) for detailed licensing guide

📧 **Commercial licensing**: [licensing@getnyla.app](mailto:licensing@getnyla.app)