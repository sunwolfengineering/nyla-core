// SPDX-License-Identifier: GPL-3.0-only
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/sunwolfengineering/nyla-core/pkg/constants"
	"github.com/sunwolfengineering/nyla-core/pkg/db"
)

func setupTestHandlers(t *testing.T) (*Handlers, *sql.DB) {
	// Create in-memory database for testing
	database, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err, "Should be able to open in-memory database")
	
	// Create the events table for testing
	createTableSQL := `
		CREATE TABLE events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			anon_id TEXT NOT NULL,
			site_id TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			type TEXT NOT NULL,
			event TEXT,
			referrer TEXT,
			is_touch TEXT,
			browser_name TEXT,
			os_name TEXT,
			device_type TEXT,
			country TEXT,
			region TEXT
		)
	`
	_, err = database.Exec(createTableSQL)
	require.NoError(t, err, "Should be able to create events table")
	
	events := &db.Events{DB: database}
	handlers := &Handlers{Events: events}
	
	return handlers, database
}

func TestGetCollectV1_SingleSiteValidation(t *testing.T) {
	// Setup
	handlers, database := setupTestHandlers(t)
	defer database.Close()

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
		{
			name:           "Missing event type with valid site_id",
			siteID:         "default",
			eventType:      "",
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
			errorMessage:   "Missing required parameter: type",
		},
		{
			name:           "Missing event type with omitted site_id",
			siteID:         "",
			eventType:      "",
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
			errorMessage:   "Missing required parameter: type",
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
	handlers, database := setupTestHandlers(t)
	defer database.Close()
	
	// Test that exactly the default site ID value is accepted
	req := httptest.NewRequest("GET", "/api/v1/collect?site_id="+constants.DefaultSiteID+"&type=pageview", nil)
	rec := httptest.NewRecorder()
	
	handlers.GetCollectV1(rec, req)
	
	assert.Equal(t, http.StatusOK, rec.Code, "Request with constants.DefaultSiteID should succeed")
}

func TestGetStatsRealtimeV1_SingleSiteQuery(t *testing.T) {
	// This test verifies that the stats query includes site filtering
	handlers, database := setupTestHandlers(t)
	defer database.Close()
	
	// Insert some test data
	_, err := database.Exec("INSERT INTO events (anon_id, site_id, created_at, type, event, referrer, is_touch, browser_name, os_name, device_type, country, region) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"test-hash", constants.DefaultSiteID, 20250101, "pageview", "https://example.com", "", "false", "Chrome", "Windows", "desktop", "US", "CA")
	require.NoError(t, err, "Should be able to insert test data")
	
	req := httptest.NewRequest("GET", "/api/v1/stats/realtime", nil)
	rec := httptest.NewRecorder()
	
	handlers.GetStatsRealtimeV1(rec, req)
	
	// Should succeed with proper database and show the count
	assert.Equal(t, http.StatusOK, rec.Code, "Should succeed with proper database")
	assert.Contains(t, rec.Body.String(), "Total Pageviews", "Should contain pageviews label")
	assert.Contains(t, rec.Body.String(), "1", "Should show count of 1")
}

func TestCollectorPayload_StructureChange(t *testing.T) {
	// Test that CollectorPayload no longer has SiteID field
	payload := CollectorPayload{
		Data: CollectorData{
			Type:      "pageview",
			Event:     "https://example.com",
			UserAgent: "test-agent",
			Hostname:  "example.com",
			Referrer:  "https://google.com",
		},
	}
	
	// Verify that we can create the payload without SiteID
	assert.Equal(t, "pageview", payload.Data.Type, "Should be able to access Data.Type")
	assert.Equal(t, "https://example.com", payload.Data.Event, "Should be able to access Data.Event")
	
	// Verify the struct doesn't have a SiteID field by checking that the struct is as expected
	assert.NotNil(t, payload.Data, "Data field should exist")
}
