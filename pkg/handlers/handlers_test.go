// SPDX-License-Identifier: GPL-3.0-only
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sunwolfengineering/nyla-core/internal/storage"
	"github.com/sunwolfengineering/nyla-core/pkg/constants"
)

func setupTestHandlers(t *testing.T) (*Handlers, *storage.DB) {
	// Create temporary database for testing
	dbPath := "test_handlers.db"
	t.Cleanup(func() {
		os.Remove(dbPath)
		os.Remove(dbPath + "-shm")
		os.Remove(dbPath + "-wal")
	})
	
	// Create temporary migrations directory for testing
	migrationsDir := setupTestMigrations(t)
	t.Cleanup(func() {
		os.RemoveAll(migrationsDir)
	})
	
	db, err := storage.NewDBWithMigrations(dbPath, migrationsDir)
	require.NoError(t, err, "Should be able to create test database")
	
	handlers := &Handlers{DB: db}
	
	return handlers, db
}

// setupTestMigrations creates a temporary migrations directory with the minimal schema
func setupTestMigrations(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "test_migrations")
	require.NoError(t, err)
	
	// Create a simplified test migration
	migration := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE site_config (
			id TEXT PRIMARY KEY DEFAULT 'default',
			name TEXT NOT NULL DEFAULT 'My Site',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			settings TEXT NOT NULL DEFAULT '{}',
			CHECK (id = 'default')
		);
		
		INSERT INTO site_config (id, name) VALUES ('default', 'My Site')
		ON CONFLICT DO NOTHING;
		
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
			CHECK (site_id = 'default')
		) STRICT;
		
		CREATE INDEX idx_events_timestamp ON events(timestamp);
		CREATE INDEX idx_events_type_timestamp ON events(type, timestamp);
		CREATE INDEX idx_events_session ON events(session_id, timestamp);
		CREATE INDEX idx_events_url ON events(url, timestamp);
		
		CREATE TABLE sessions (
			id TEXT PRIMARY KEY,
			site_id TEXT NOT NULL DEFAULT 'default',
			started_at TEXT NOT NULL,
			ended_at TEXT,
			duration INTEGER,
			pages_viewed INTEGER DEFAULT 0,
			entry_page TEXT,
			exit_page TEXT,
			referrer TEXT,
			metadata TEXT,
			FOREIGN KEY(site_id) REFERENCES site_config(id),
			CHECK (site_id = 'default')
		) STRICT;
		
		CREATE INDEX idx_sessions_time ON sessions(started_at);
		CREATE INDEX idx_sessions_duration ON sessions(duration);
		
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
		
		INSERT INTO schema_migrations (version) VALUES (1);
	`
	
	migrationPath := t.TempDir() + "/001_test_schema.sql"
	err = os.WriteFile(migrationPath, []byte(migration), 0644)
	require.NoError(t, err)
	
	return tempDir
}

func TestGetCollectV1_SingleSiteValidation(t *testing.T) {
	// Setup
	handlers, db := setupTestHandlers(t)
	defer db.Close()

	tests := []struct {
		name           string
		siteID         string
		eventType      string
		expectedStatus int
		expectJSON     bool
		errorMessage   string
	}{
		{
			name:           "Valid default site_id",
			siteID:         "default",
			eventType:      "pageview",
			expectedStatus: http.StatusOK,
			expectJSON:     false,
		},
		{
			name:           "Omitted site_id (should default to valid)",
			siteID:         "",
			eventType:      "pageview",
			expectedStatus: http.StatusOK,
			expectJSON:     false,
		},
		{
			name:           "Invalid site_id",
			siteID:         "invalid",
			eventType:      "pageview",
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
			errorMessage:   "Invalid site_id. This instance only supports site_id='default'",
		},
		{
			name:           "Invalid site_id (different value)",
			siteID:         "another-site",
			eventType:      "pageview",
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
			errorMessage:   "Invalid site_id. This instance only supports site_id='default'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request URL
			url := "/api/v1/collect"
			params := []string{}
			if tt.siteID != "" {
				params = append(params, "site_id="+tt.siteID)
			}
			if tt.eventType != "" {
				params = append(params, "type="+tt.eventType)
			}
			// Always add a URL parameter to make it a valid request
			params = append(params, "url=https://example.com")
			
			if len(params) > 0 {
				url += "?" + strings.Join(params, "&")
			}

			// Create request
			req := httptest.NewRequest("GET", url, nil)
			rec := httptest.NewRecorder()

			// Execute handler
			handlers.GetCollectV1(rec, req)

			// Verify status code
			assert.Equal(t, tt.expectedStatus, rec.Code, "Status code should match expected")

			if tt.expectJSON {
				// Verify JSON error response
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"), "Content-Type should be application/json for error responses")
				
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err, "Should be able to parse JSON response")
				
				assert.Equal(t, tt.errorMessage, response["error"], "Error message should match expected")
			} else {
				// Verify GIF response for successful requests
				assert.Equal(t, "image/gif", rec.Header().Get("Content-Type"), "Content-Type should be image/gif for successful requests")
				assert.Equal(t, 43, rec.Body.Len(), "GIF response should be 43 bytes")
			}
		})
	}
}

func TestGetCollectV1_DefaultSiteIDUsage(t *testing.T) {
	// Test that the handler uses constants.DefaultSiteID internally
	assert.Equal(t, "default", constants.DefaultSiteID, "Test assumes DefaultSiteID is 'default'")
	
	// This test verifies that our validation logic correctly uses the constant
	handlers, db := setupTestHandlers(t)
	defer db.Close()
	
	// Test that exactly the default site ID value is accepted
	req := httptest.NewRequest("GET", "/api/v1/collect?site_id="+constants.DefaultSiteID+"&type=pageview&url=https://example.com", nil)
	rec := httptest.NewRecorder()
	
	handlers.GetCollectV1(rec, req)
	
	assert.Equal(t, http.StatusOK, rec.Code, "Request with constants.DefaultSiteID should succeed")
}

func TestGetStatsRealtimeV1_ReturnsStats(t *testing.T) {
	// This test verifies that the stats endpoint works with new storage
	handlers, db := setupTestHandlers(t)
	defer db.Close()
	
	req := httptest.NewRequest("GET", "/api/v1/stats/realtime", nil)
	rec := httptest.NewRecorder()
	
	handlers.GetStatsRealtimeV1(rec, req)
	
	// Should succeed with proper database
	assert.Equal(t, http.StatusOK, rec.Code, "Should succeed with proper database")
	assert.Contains(t, rec.Body.String(), "Total Pageviews", "Should contain pageviews label")
}

func TestCollectorPayload_Structure(t *testing.T) {
	// Test that CollectorPayload structure works as expected
	payload := CollectorPayload{
		Data: CollectorData{
			Type:      "pageview",
			Event:     "https://example.com",
			UserAgent: "test-agent",
			Hostname:  "example.com",
			Referrer:  "https://google.com",
		},
	}
	
	// Verify that we can create the payload
	assert.Equal(t, "pageview", payload.Data.Type, "Should be able to access Data.Type")
	assert.Equal(t, "https://example.com", payload.Data.Event, "Should be able to access Data.Event")
	
	// Verify the struct is as expected
	assert.NotNil(t, payload.Data, "Data field should exist")
}
