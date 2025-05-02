# Nyla Analytics - Database Schema Specification

## Overview

The Nyla Analytics database schema is designed for efficient analytics data storage and querying using SQLite. The schema prioritizes query performance for common analytics operations while maintaining data integrity and supporting privacy requirements.

## Schema Version Control

```sql
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Core Tables

### Sites

```sql
CREATE TABLE sites (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    settings JSON NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_sites_created_at ON sites(created_at);
```

### Events

```sql
CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL,
    type TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    url TEXT,
    title TEXT,
    referrer TEXT,
    session_id TEXT,
    metadata JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(site_id) REFERENCES sites(id)
) STRICT;

-- Indexes for common queries
CREATE INDEX idx_events_site_timestamp ON events(site_id, timestamp);
CREATE INDEX idx_events_type_timestamp ON events(type, timestamp);
CREATE INDEX idx_events_session ON events(session_id, timestamp);
CREATE INDEX idx_events_url ON events(url, timestamp);
```

### Sessions

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    site_id TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration INTEGER, -- in seconds
    pages_viewed INTEGER DEFAULT 0,
    entry_page TEXT,
    exit_page TEXT,
    referrer TEXT,
    metadata JSON,
    FOREIGN KEY(site_id) REFERENCES sites(id)
) STRICT;

CREATE INDEX idx_sessions_site_time ON sessions(site_id, started_at);
CREATE INDEX idx_sessions_duration ON sessions(duration);
```

### Aggregates

```sql
CREATE TABLE daily_aggregates (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL,
    date DATE NOT NULL,
    pageviews INTEGER DEFAULT 0,
    unique_visitors INTEGER DEFAULT 0,
    total_sessions INTEGER DEFAULT 0,
    avg_session_duration REAL,
    bounce_rate REAL,
    metadata JSON,
    FOREIGN KEY(site_id) REFERENCES sites(id),
    UNIQUE(site_id, date)
) STRICT;

CREATE INDEX idx_daily_aggregates_site_date 
    ON daily_aggregates(site_id, date);
```

## Privacy & Retention

### Data Retention

```sql
CREATE TABLE retention_policies (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL,
    data_type TEXT NOT NULL, -- 'events', 'sessions', etc.
    retention_days INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(site_id) REFERENCES sites(id),
    UNIQUE(site_id, data_type)
);
```

### Privacy Logs

```sql
CREATE TABLE privacy_logs (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL,
    action TEXT NOT NULL, -- 'anonymize', 'delete', etc.
    data_type TEXT NOT NULL,
    identifier TEXT, -- session_id, url, etc.
    performed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata JSON,
    FOREIGN KEY(site_id) REFERENCES sites(id)
);

CREATE INDEX idx_privacy_logs_site_action 
    ON privacy_logs(site_id, action, performed_at);
```

## Views

### Active Visitors

```sql
CREATE VIEW active_visitors AS
SELECT 
    site_id,
    COUNT(DISTINCT session_id) as visitor_count
FROM events
WHERE timestamp >= datetime('now', '-30 minutes')
GROUP BY site_id;
```

### Popular Pages

```sql
CREATE VIEW popular_pages AS
SELECT 
    site_id,
    url,
    COUNT(*) as pageviews,
    COUNT(DISTINCT session_id) as unique_views
FROM events
WHERE 
    type = 'pageview'
    AND timestamp >= datetime('now', '-24 hours')
GROUP BY site_id, url;
```

## Triggers

### Updated Timestamp

```sql
CREATE TRIGGER sites_updated_at
AFTER UPDATE ON sites
BEGIN
    UPDATE sites 
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;
```

### Session Updates

```sql
CREATE TRIGGER update_session_stats
AFTER INSERT ON events
WHEN NEW.type = 'pageview'
BEGIN
    UPDATE sessions
    SET 
        pages_viewed = pages_viewed + 1,
        ended_at = NEW.timestamp,
        duration = CAST(
            (strftime('%s', NEW.timestamp) - 
             strftime('%s', started_at)) AS INTEGER
        )
    WHERE id = NEW.session_id;
END;
```

## Functions

### Privacy Helpers

```sql
-- Anonymize URL parameters
CREATE FUNCTION anonymize_url(url TEXT)
RETURNS TEXT
BEGIN
    -- Strip query parameters containing sensitive patterns
    RETURN regexp_replace(
        url,
        '([?&](email|token|key)=[^&]*)',
        '\\1=REDACTED'
    );
END;

-- Calculate retention window
CREATE FUNCTION retention_window(site_id TEXT, data_type TEXT)
RETURNS TIMESTAMP
BEGIN
    RETURN datetime(
        'now',
        '-' || (
            SELECT retention_days 
            FROM retention_policies
            WHERE site_id = site_id 
            AND data_type = data_type
        ) || ' days'
    );
END;
```

## Maintenance

### Cleanup Jobs

```sql
-- Delete expired events
DELETE FROM events
WHERE timestamp < retention_window(site_id, 'events');

-- Delete expired sessions
DELETE FROM sessions
WHERE ended_at < retention_window(site_id, 'sessions');

-- Vacuum database periodically
VACUUM;
```

### Optimization

```sql
-- Analyze tables for query optimization
ANALYZE events;
ANALYZE sessions;
ANALYZE daily_aggregates;

-- Reindex for fragmented indexes
REINDEX idx_events_site_timestamp;
REINDEX idx_sessions_site_time;
```

## Performance Considerations

1. WAL Journal Mode
```sql
PRAGMA journal_mode = WAL;
```

2. Memory Settings
```sql
PRAGMA cache_size = -2000; -- 2MB cache
PRAGMA temp_store = MEMORY;
```

3. Busy Timeout
```sql
PRAGMA busy_timeout = 5000;
```

4. Foreign Key Support
```sql
PRAGMA foreign_keys = ON;
```

## Backup Configuration

```sql
-- Backup settings
PRAGMA auto_vacuum = INCREMENTAL;
PRAGMA page_size = 4096;
```

## Migration Guidelines

1. Always use transactions for schema changes
2. Add new columns as nullable or with defaults
3. Create new indexes before dropping old ones
4. Maintain backward compatibility where possible
5. Include down migrations for rollback support

## Data Types

- Use TEXT for IDs and URLs
- Use INTEGER for counts and flags
- Use REAL for percentages and ratios
- Use TIMESTAMP for dates (ISO8601 strings)
- Use JSON for flexible metadata storage

## Security

1. Table and column permissions
2. Input validation
3. Prepared statements
4. Regular backup verification
5. Encryption at rest support 