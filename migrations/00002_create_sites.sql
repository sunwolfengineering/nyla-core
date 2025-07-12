-- +goose Up
CREATE TABLE sites (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    settings JSON NOT NULL DEFAULT '{}'
);

-- +goose StatementBegin
CREATE INDEX idx_sites_created_at ON sites(created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER sites_updated_at
AFTER UPDATE ON sites
FOR EACH ROW
BEGIN
    UPDATE sites 
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS sites_updated_at;
DROP INDEX IF EXISTS idx_sites_created_at;
DROP TABLE IF EXISTS sites;
