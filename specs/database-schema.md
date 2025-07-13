# Nyla Analytics - Database Schema Specification (Core)

## Overview

The Nyla Analytics Core database schema is designed for efficient single-site analytics data storage and querying using SQLite. The schema prioritizes query performance for common analytics operations while maintaining data integrity and supporting privacy requirements.

The Core schema provides efficient single-site analytics data storage and querying using SQLite.

## Schema Version Control

```sql
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Core Tables

### Site Configuration (Core - Single Site Only)

```sql
CREATE TABLE site_config (
    id TEXT PRIMARY KEY DEFAULT 'default',
    name TEXT NOT NULL DEFAULT 'My Site',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    settings JSON NOT NULL DEFAULT '{}',
    CHECK (id = 'default') -- Enforce single site in core
);

-- Insert default site configuration
INSERT INTO site_config (id, name) VALUES ('default', 'My Site')
ON CONFLICT DO NOTHING;
```

### Events

```sql
CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    type TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    url TEXT,
    title TEXT,
    referrer TEXT,
    session_id TEXT,
    metadata JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default') -- Enforce single site in core
) STRICT;

-- Indexes for common queries (optimized for single site)
CREATE INDEX idx_events_timestamp ON events(timestamp);
CREATE INDEX idx_events_type_timestamp ON events(type, timestamp);
CREATE INDEX idx_events_session ON events(session_id, timestamp);
CREATE INDEX idx_events_url ON events(url, timestamp);
```

### Sessions

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration INTEGER, -- in seconds
    pages_viewed INTEGER DEFAULT 0,
    entry_page TEXT,
    exit_page TEXT,
    referrer TEXT,
    metadata JSON,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default') -- Enforce single site in core
) STRICT;

CREATE INDEX idx_sessions_time ON sessions(started_at);
CREATE INDEX idx_sessions_duration ON sessions(duration);
```

### Aggregates

```sql
CREATE TABLE daily_aggregates (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    date DATE NOT NULL,
    pageviews INTEGER DEFAULT 0,
    unique_visitors INTEGER DEFAULT 0,
    total_sessions INTEGER DEFAULT 0,
    avg_session_duration REAL,
    bounce_rate REAL,
    metadata JSON,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default'), -- Enforce single site in core
    UNIQUE(site_id, date)
) STRICT;

CREATE INDEX idx_daily_aggregates_date ON daily_aggregates(date);
```

## Privacy & Retention

### Data Retention

```sql
CREATE TABLE retention_policies (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    data_type TEXT NOT NULL, -- 'events', 'sessions', etc.
    retention_days INTEGER NOT NULL DEFAULT 90,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default'), -- Enforce single site in core
    UNIQUE(site_id, data_type)
);

-- Insert default retention policies
INSERT INTO retention_policies (site_id, data_type, retention_days) VALUES 
    ('default', 'events', 90),
    ('default', 'sessions', 90)
ON CONFLICT DO NOTHING;
```

### Privacy Logs

```sql
CREATE TABLE privacy_logs (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    action TEXT NOT NULL, -- 'anonymize', 'delete', etc.
    data_type TEXT NOT NULL,
    identifier TEXT, -- session_id, url, etc.
    performed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata JSON,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default') -- Enforce single site in core
);

CREATE INDEX idx_privacy_logs_action 
    ON privacy_logs(action, performed_at);
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