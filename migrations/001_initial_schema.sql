-- Nyla Analytics Core - Initial Schema
-- Version: 001
-- Applied: Single-site architecture with full schema from specification

-- Schema version tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Site Configuration (Core - Single Site Only)
CREATE TABLE site_config (
    id TEXT PRIMARY KEY DEFAULT 'default',
    name TEXT NOT NULL DEFAULT 'My Site',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    settings TEXT NOT NULL DEFAULT '{}',
    CHECK (id = 'default') -- Enforce single site in core
);

-- Insert default site configuration
INSERT INTO site_config (id, name) VALUES ('default', 'My Site')
ON CONFLICT DO NOTHING;

-- Events
CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    type TEXT NOT NULL,
    timestamp TEXT NOT NULL,
    url TEXT,
    title TEXT,
    referrer TEXT,
    session_id TEXT,
    metadata TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default') -- Enforce single site in core
) STRICT;

-- Indexes for common queries (optimized for single site)
CREATE INDEX idx_events_timestamp ON events(timestamp);
CREATE INDEX idx_events_type_timestamp ON events(type, timestamp);
CREATE INDEX idx_events_session ON events(session_id, timestamp);
CREATE INDEX idx_events_url ON events(url, timestamp);

-- Sessions
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    started_at TEXT NOT NULL,
    ended_at TEXT,
    duration INTEGER, -- in seconds
    pages_viewed INTEGER DEFAULT 0,
    entry_page TEXT,
    exit_page TEXT,
    referrer TEXT,
    metadata TEXT,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default') -- Enforce single site in core
) STRICT;

CREATE INDEX idx_sessions_time ON sessions(started_at);
CREATE INDEX idx_sessions_duration ON sessions(duration);

-- Aggregates
CREATE TABLE daily_aggregates (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    date TEXT NOT NULL,
    pageviews INTEGER DEFAULT 0,
    unique_visitors INTEGER DEFAULT 0,
    total_sessions INTEGER DEFAULT 0,
    avg_session_duration REAL,
    bounce_rate REAL,
    metadata TEXT,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default'), -- Enforce single site in core
    UNIQUE(site_id, date)
) STRICT;

CREATE INDEX idx_daily_aggregates_date ON daily_aggregates(date);

-- Data Retention Policies
CREATE TABLE retention_policies (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    data_type TEXT NOT NULL, -- 'events', 'sessions', etc.
    retention_days INTEGER NOT NULL DEFAULT 90,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default'), -- Enforce single site in core
    UNIQUE(site_id, data_type)
);

-- Insert default retention policies
INSERT INTO retention_policies (site_id, data_type, retention_days) VALUES 
    ('default', 'events', 90),
    ('default', 'sessions', 90)
ON CONFLICT DO NOTHING;

-- Privacy Logs
CREATE TABLE privacy_logs (
    id INTEGER PRIMARY KEY,
    site_id TEXT NOT NULL DEFAULT 'default',
    action TEXT NOT NULL, -- 'anonymize', 'delete', etc.
    data_type TEXT NOT NULL,
    identifier TEXT, -- session_id, url, etc.
    performed_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata TEXT,
    FOREIGN KEY(site_id) REFERENCES site_config(id),
    CHECK (site_id = 'default') -- Enforce single site in core
);

CREATE INDEX idx_privacy_logs_action 
    ON privacy_logs(action, performed_at);

-- Views
CREATE VIEW active_visitors AS
SELECT 
    site_id,
    COUNT(DISTINCT session_id) as visitor_count
FROM events
WHERE timestamp >= datetime('now', '-30 minutes')
GROUP BY site_id;

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

-- Triggers
CREATE TRIGGER site_config_updated_at
AFTER UPDATE ON site_config
BEGIN
    UPDATE site_config 
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

CREATE TRIGGER update_session_stats
AFTER INSERT ON events
WHEN NEW.type = 'pageview'
BEGIN
    INSERT INTO sessions (
        id, 
        site_id, 
        started_at, 
        ended_at,
        pages_viewed,
        entry_page
    ) VALUES (
        NEW.session_id,
        NEW.site_id,
        NEW.timestamp,
        NEW.timestamp,
        1,
        NEW.url
    )
    ON CONFLICT(id) DO UPDATE SET 
        pages_viewed = pages_viewed + 1,
        ended_at = NEW.timestamp,
        exit_page = NEW.url,
        duration = CAST(
            (strftime('%s', NEW.timestamp) - 
             strftime('%s', started_at)) AS INTEGER
        );
END;

-- Record this migration
INSERT INTO schema_migrations (version) VALUES (1);
