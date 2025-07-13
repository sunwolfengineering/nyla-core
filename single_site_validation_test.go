// SPDX-License-Identifier: GPL-3.0-only
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/sunwolfengineering/nyla-core/internal/server"
	"github.com/sunwolfengineering/nyla-core/pkg/constants"
	"github.com/sunwolfengineering/nyla-core/pkg/db"
)

// TestSingleSiteArchitectureCompliance validates the complete single-site enforcement
// across all components that were modified in NYLA-35
func TestSingleSiteArchitectureCompliance(t *testing.T) {
	// Setup test server with database
	database, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err, "Should be able to create test database")
	defer database.Close()
	
	// Create events table
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
	srv := server.New(events)
	
	t.Run("Constant Value Enforcement", func(t *testing.T) {
		// Verify the default site ID constant has the expected value
		assert.Equal(t, "default", constants.DefaultSiteID, "DefaultSiteID constant must be 'default'")
	})
	
	t.Run("HTTP API Validation", func(t *testing.T) {
		testCases := []struct {
			name           string
			query          string
			expectStatus   int
			expectJSON     bool
			expectErrorMsg string
		}{
			{
				name:         "Accept default site_id",
				query:        "?site_id=default&type=pageview",
				expectStatus: http.StatusOK,
				expectJSON:   false,
			},
			{
				name:         "Accept omitted site_id",
				query:        "?type=pageview",
				expectStatus: http.StatusOK,
				expectJSON:   false,
			},
			{
				name:           "Reject invalid site_id",
				query:          "?site_id=another-site&type=pageview",
				expectStatus:   http.StatusBadRequest,
				expectJSON:     true,
				expectErrorMsg: "Invalid site_id. This instance only supports site_id='default'",
			},
			{
				name:         "Accept empty string site_id (treated same as omitted)",
				query:        "?site_id=&type=pageview",
				expectStatus: http.StatusOK,
				expectJSON:   false,
			},
		}
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest("GET", "/api/v1/collect"+tc.query, nil)
				rec := httptest.NewRecorder()
				
				srv.ServeHTTP(rec, req)
				
				assert.Equal(t, tc.expectStatus, rec.Code, "Status code should match")
				
				if tc.expectJSON {
					assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
					
					var response map[string]string
					err := json.Unmarshal(rec.Body.Bytes(), &response)
					require.NoError(t, err, "Response should be valid JSON")
					assert.Equal(t, tc.expectErrorMsg, response["error"], "Error message should match")
				} else {
					assert.Equal(t, "image/gif", rec.Header().Get("Content-Type"))
				}
			})
		}
	})
	
	t.Run("Database Storage Enforcement", func(t *testing.T) {
		// Clear any existing data
		_, err := database.Exec("DELETE FROM events")
		require.NoError(t, err, "Should be able to clear events")
		
		// Make successful requests
		validRequests := []string{
			"/api/v1/collect?site_id=default&type=pageview&url=test1",
			"/api/v1/collect?type=pageview&url=test2",
			"/api/v1/collect?site_id=default&type=event&url=test3",
		}
		
		for _, url := range validRequests {
			req := httptest.NewRequest("GET", url, nil)
			rec := httptest.NewRecorder()
			srv.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code, "Request should succeed: %s", url)
		}
		
		// Verify all stored events have the default site_id
		rows, err := database.Query("SELECT site_id FROM events")
		require.NoError(t, err, "Should be able to query events")
		defer rows.Close()
		
		for rows.Next() {
			var siteID string
			err := rows.Scan(&siteID)
			require.NoError(t, err, "Should be able to scan site_id")
			assert.Equal(t, constants.DefaultSiteID, siteID, "All stored events should have default site_id")
		}
		
		// Verify we have exactly the expected number of events
		var count int
		err = database.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
		require.NoError(t, err, "Should be able to count events")
		assert.Equal(t, len(validRequests), count, "Should have stored all valid requests")
	})
	
	t.Run("Statistics Query Filtering", func(t *testing.T) {
		// Insert mixed data to test filtering
		_, err := database.Exec("DELETE FROM events")
		require.NoError(t, err, "Should be able to clear events")
		
		// Insert events with default site_id
		for i := 0; i < 5; i++ {
			_, err := database.Exec(`INSERT INTO events 
				(anon_id, site_id, created_at, type, event, referrer, is_touch, browser_name, os_name, device_type, country, region) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				"hash-"+string(rune(i+'1')), constants.DefaultSiteID, 20250101, "pageview", 
				"https://example.com", "", "false", "Chrome", "Windows", "desktop", "US", "CA")
			require.NoError(t, err, "Should be able to insert default site event")
		}
		
		// Insert events with non-default site_id (simulating hypothetical future data)
		for i := 0; i < 3; i++ {
			_, err := database.Exec(`INSERT INTO events 
				(anon_id, site_id, created_at, type, event, referrer, is_touch, browser_name, os_name, device_type, country, region) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				"other-hash-"+string(rune(i+'1')), "other-site", 20250101, "pageview", 
				"https://other.com", "", "false", "Chrome", "Windows", "desktop", "US", "CA")
			require.NoError(t, err, "Should be able to insert other site event")
		}
		
		// Test stats endpoint - should only count default site events
		req := httptest.NewRequest("GET", "/api/v1/stats/realtime", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		
		assert.Equal(t, http.StatusOK, rec.Code, "Stats request should succeed")
		
		body := rec.Body.String()
		assert.Contains(t, body, "5", "Should show count of 5 (default site events only)")
		assert.NotContains(t, body, "8", "Should not show total count including other sites")
	})
	
	t.Run("End-to-End Single Site Workflow", func(t *testing.T) {
		// Clear events
		_, err = database.Exec("DELETE FROM events")
		require.NoError(t, err, "Should be able to clear events")
		
		// 1. Collect some events
		collectReq := httptest.NewRequest("GET", "/api/v1/collect?type=pageview&url=https://example.com/test", nil)
		collectRec := httptest.NewRecorder()
		srv.ServeHTTP(collectRec, collectReq)
		assert.Equal(t, http.StatusOK, collectRec.Code, "Collection should succeed")
		
		// 2. Get stats
		statsReq := httptest.NewRequest("GET", "/api/v1/stats/realtime", nil)
		statsRec := httptest.NewRecorder()
		srv.ServeHTTP(statsRec, statsReq)
		assert.Equal(t, http.StatusOK, statsRec.Code, "Stats should succeed")
		assert.Contains(t, statsRec.Body.String(), "1", "Should show 1 collected event")
		
		// 3. Verify dashboard loads
		dashboardReq := httptest.NewRequest("GET", "/", nil)
		dashboardRec := httptest.NewRecorder()
		srv.ServeHTTP(dashboardRec, dashboardReq)
		assert.Equal(t, http.StatusOK, dashboardRec.Code, "Dashboard should load")
		assert.Contains(t, dashboardRec.Body.String(), "Nyla Analytics", "Dashboard should contain title")
		
		// 4. Verify database consistency
		var storedSiteID string
		err = database.QueryRow("SELECT site_id FROM events LIMIT 1").Scan(&storedSiteID)
		require.NoError(t, err, "Should have stored event")
		assert.Equal(t, constants.DefaultSiteID, storedSiteID, "Stored event should have default site_id")
	})
}
