-- +goose Up
CREATE TABLE IF NOT EXISTS events (
  anon_id TEXT NOT NULL,
  site_id TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  type TEXT NOT NULL,
  event TEXT NOT NULL,
  referrer TEXT NOT NULL,
  is_touch INTEGER NOT NULL,
  browser_name TEXT NOT NULL,
  os_name TEXT NOT NULL,
  device_type TEXT NOT NULL,
  country TEXT NOT NULL,
  region TEXT NOT NULL,
  timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS events; 