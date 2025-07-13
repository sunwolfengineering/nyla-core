package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"

	"github.com/sunwolfengineering/nyla-core/pkg/constants"
)

// DB represents the database connection and operations
type DB struct {
	conn *sql.DB
	path string
}

// Event represents an analytics event
type Event struct {
	ID        int64             `json:"id"`
	SiteID    string            `json:"site_id"`
	Type      string            `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	URL       string            `json:"url,omitempty"`
	Title     string            `json:"title,omitempty"`
	Referrer  string            `json:"referrer,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// Session represents a user session
type Session struct {
	ID               string            `json:"id"`
	SiteID           string            `json:"site_id"`
	StartedAt        time.Time         `json:"started_at"`
	EndedAt          *time.Time        `json:"ended_at,omitempty"`
	Duration         *int              `json:"duration,omitempty"`
	PagesViewed      int               `json:"pages_viewed"`
	EntryPage        string            `json:"entry_page,omitempty"`
	ExitPage         string            `json:"exit_page,omitempty"`
	Referrer         string            `json:"referrer,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// RealtimeStats represents real-time statistics
type RealtimeStats struct {
	ActiveVisitors int `json:"active_visitors"`
	PageviewsToday int `json:"pageviews_today"`
	TotalSessions  int `json:"total_sessions"`
}

// NewDB creates a new database connection with optimal SQLite settings
func NewDB(dbPath string) (*DB, error) {
	return NewDBWithMigrations(dbPath, "migrations")
}

// NewDBWithMigrations creates a new database connection and runs migrations from the specified path
func NewDBWithMigrations(dbPath, migrationsPath string) (*DB, error) {
	db := &DB{path: dbPath}
	
	if err := db.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	if err := db.configureDatabase(); err != nil {
		return nil, fmt.Errorf("failed to configure database: %w", err)
	}
	
	// Run migrations if migrations path is provided
	if migrationsPath != "" {
		migrationRunner := NewMigrationRunner(db.conn)
		if err := migrationRunner.Run(migrationsPath); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}
	
	return db, nil
}

// connect establishes the database connection
func (db *DB) connect() error {
	conn, err := sql.Open("sqlite", db.path)
	if err != nil {
		return err
	}
	
	if err := conn.Ping(); err != nil {
		conn.Close()
		return err
	}
	
	db.conn = conn
	return nil
}

// configureDatabase applies optimal SQLite settings for analytics workload
func (db *DB) configureDatabase() error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL",           // Write-Ahead Logging for better concurrency
		"PRAGMA synchronous = NORMAL",         // Good balance of performance and safety
		"PRAGMA cache_size = -2000",          // 2MB cache
		"PRAGMA temp_store = MEMORY",         // Store temp tables in memory
		"PRAGMA busy_timeout = 5000",         // 5 second timeout for busy database
		"PRAGMA foreign_keys = ON",           // Enable foreign key constraints
		"PRAGMA auto_vacuum = INCREMENTAL",   // Enable incremental vacuum
		"PRAGMA page_size = 4096",           // Optimal page size for most systems
		"PRAGMA mmap_size = 268435456",       // 256MB memory mapping
	}
	
	for _, pragma := range pragmas {
		if _, err := db.conn.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %s: %w", pragma, err)
		}
	}
	
	log.Println("Database configured with optimal settings")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// InsertEvent inserts a new event into the database
func (db *DB) InsertEvent(ctx context.Context, event *Event) error {
	// Always use default site ID for single-site architecture
	event.SiteID = constants.DefaultSiteID
	
	// Serialize metadata to JSON string
	var metadataJSON string
	if event.Metadata != nil {
		data, err := json.Marshal(event.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = string(data)
	}
	
	query := `
		INSERT INTO events (
			site_id, type, timestamp, url, title, referrer, session_id, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := db.conn.ExecContext(
		ctx, query,
		event.SiteID,
		event.Type,
		event.Timestamp.Format(time.RFC3339),
		event.URL,
		event.Title,
		event.Referrer,
		event.SessionID,
		metadataJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}
	
	event.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get inserted event ID: %w", err)
	}
	
	return nil
}

// GetRealtimeStats returns real-time analytics statistics
func (db *DB) GetRealtimeStats(ctx context.Context) (*RealtimeStats, error) {
	stats := &RealtimeStats{}
	
	// Get active visitors (last 30 minutes)
	// Use datetime comparison with ISO8601 timestamps
	err := db.conn.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT session_id) 
		FROM events 
		WHERE site_id = ? 
		AND timestamp >= datetime('now', '-30 minutes')
		AND session_id IS NOT NULL
	`, constants.DefaultSiteID).Scan(&stats.ActiveVisitors)
	if err != nil {
		return nil, fmt.Errorf("failed to get active visitors: %w", err)
	}
	
	// Get pageviews today
	err = db.conn.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM events 
		WHERE site_id = ? 
		AND type = 'pageview'
		AND date(timestamp) = date('now')
	`, constants.DefaultSiteID).Scan(&stats.PageviewsToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get pageviews today: %w", err)
	}
	
	// Get total sessions today
	err = db.conn.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM sessions 
		WHERE site_id = ? 
		AND date(started_at) = date('now')
	`, constants.DefaultSiteID).Scan(&stats.TotalSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to get total sessions: %w", err)
	}
	
	return stats, nil
}

// GetPopularPages returns the most popular pages in the last 24 hours
func (db *DB) GetPopularPages(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT url, COUNT(*) as pageviews, COUNT(DISTINCT session_id) as unique_views
		FROM events 
		WHERE site_id = ? 
		AND type = 'pageview'
		AND timestamp >= datetime('now', '-24 hours')
		AND url IS NOT NULL
		GROUP BY url
		ORDER BY pageviews DESC
		LIMIT ?`
	
	rows, err := db.conn.QueryContext(ctx, query, constants.DefaultSiteID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query popular pages: %w", err)
	}
	defer rows.Close()
	
	var pages []map[string]interface{}
	for rows.Next() {
		var url string
		var pageviews, uniqueViews int
		
		if err := rows.Scan(&url, &pageviews, &uniqueViews); err != nil {
			return nil, fmt.Errorf("failed to scan popular page row: %w", err)
		}
		
		pages = append(pages, map[string]interface{}{
			"url":           url,
			"pageviews":     pageviews,
			"unique_views":  uniqueViews,
		})
	}
	
	return pages, nil
}

// GetSessionByID retrieves a session by its ID
func (db *DB) GetSessionByID(ctx context.Context, sessionID string) (*Session, error) {
	session := &Session{}
	var metadataJSON sql.NullString
	var endedAtStr sql.NullString
	var startedAtStr string
	var duration sql.NullInt64
	var entryPage, exitPage, referrer sql.NullString
	
	query := `
		SELECT id, site_id, started_at, ended_at, duration, pages_viewed,
		       entry_page, exit_page, referrer, metadata
		FROM sessions 
		WHERE id = ? AND site_id = ?`
	
	err := db.conn.QueryRowContext(ctx, query, sessionID, constants.DefaultSiteID).Scan(
		&session.ID,
		&session.SiteID,
		&startedAtStr,
		&endedAtStr,
		&duration,
		&session.PagesViewed,
		&entryPage,
		&exitPage,
		&referrer,
		&metadataJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Session not found
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	
	// Parse timestamps
	session.StartedAt, err = time.Parse(time.RFC3339, startedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse started_at timestamp: %w", err)
	}
	
	// Handle nullable fields
	if endedAtStr.Valid {
		endedAt, err := time.Parse(time.RFC3339, endedAtStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ended_at timestamp: %w", err)
		}
		session.EndedAt = &endedAt
	}
	if duration.Valid {
		d := int(duration.Int64)
		session.Duration = &d
	}
	if entryPage.Valid {
		session.EntryPage = entryPage.String
	}
	if exitPage.Valid {
		session.ExitPage = exitPage.String
	}
	if referrer.Valid {
		session.Referrer = referrer.String
	}
	
	// Parse metadata JSON
	if metadataJSON.Valid && metadataJSON.String != "" {
		if err := json.Unmarshal([]byte(metadataJSON.String), &session.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal session metadata: %w", err)
		}
	}
	
	return session, nil
}

// Ping checks if the database connection is healthy
func (db *DB) Ping(ctx context.Context) error {
	return db.conn.PingContext(ctx)
}

// Cleanup performs database maintenance operations
func (db *DB) Cleanup(ctx context.Context) error {
	// Run ANALYZE to update table statistics
	if _, err := db.conn.ExecContext(ctx, "ANALYZE"); err != nil {
		return fmt.Errorf("failed to analyze database: %w", err)
	}
	
	// Run incremental vacuum
	if _, err := db.conn.ExecContext(ctx, "PRAGMA incremental_vacuum"); err != nil {
		return fmt.Errorf("failed to run incremental vacuum: %w", err)
	}
	
	log.Println("Database cleanup completed")
	return nil
}
