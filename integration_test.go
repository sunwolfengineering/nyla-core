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

func setupTestServer(t *testing.T) (*server.Server, *sql.DB) {
	// Create in-memory database
	database, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err, "Should be able to create test database")
	
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
	
	// Create Events instance
	events := &db.Events{DB: database}
	
	// Create server
	srv := server.New(events)
	
	return srv, database
}

func TestIntegration_CollectEndpoint_SingleSiteValidation(t *testing.T) {
	srv, database := setupTestServer(t)
	defer database.Close()
	
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectJSON     bool
	}{
		{
			name:           "Valid request with default site_id",
			url:            "/api/v1/collect?site_id=default&type=pageview&url=https://example.com",
			expectedStatus: http.StatusOK,
			expectJSON:     false,
		},
		{
			name:           "Valid request without site_id",
			url:            "/api/v1/collect?type=pageview&url=https://example.com",
			expectedStatus: http.StatusOK,
			expectJSON:     false,
		},
		{
			name:           "Invalid site_id should be rejected",
			url:            "/api/v1/collect?site_id=invalid&type=pageview&url=https://example.com",
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
		},
		{
			name:           "Missing type parameter",
			url:            "/api/v1/collect?site_id=default&url=https://example.com",
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			rec := httptest.NewRecorder()
			
			srv.ServeHTTP(rec, req)
			
			assert.Equal(t, tt.expectedStatus, rec.Code, "Status code should match")
			
			if tt.expectJSON {
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
				
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err, "Should be valid JSON")
				assert.Contains(t, response, "error", "Should contain error field")
			} else {
				assert.Equal(t, "image/gif", rec.Header().Get("Content-Type"))
				assert.Equal(t, 43, rec.Body.Len(), "Should return GIF")
			}
		})
	}
}

func TestIntegration_DatabaseStorage_DefaultSiteOnly(t *testing.T) {
	srv, database := setupTestServer(t)
	defer database.Close()
	
	// Make several successful requests
	requests := []string{
		"/api/v1/collect?site_id=default&type=pageview&url=https://example.com/page1",
		"/api/v1/collect?type=pageview&url=https://example.com/page2",
		"/api/v1/collect?site_id=default&type=pageview&url=https://example.com/page3",
	}
	
	for i, url := range requests {
		req := httptest.NewRequest("GET", url, nil)
		rec := httptest.NewRecorder()
		
		srv.ServeHTTP(rec, req)
		
		assert.Equal(t, http.StatusOK, rec.Code, "Request %d should succeed", i+1)
	}
	
	// Verify all events were stored with default site_id
	rows, err := database.Query("SELECT DISTINCT site_id FROM events")
	require.NoError(t, err, "Should be able to query site_ids")
	defer rows.Close()
	
	siteIDs := []string{}
	for rows.Next() {
		var siteID string
		err := rows.Scan(&siteID)
		require.NoError(t, err, "Should be able to scan site_id")
		siteIDs = append(siteIDs, siteID)
	}
	
	assert.Len(t, siteIDs, 1, "Should only have one unique site_id")
	assert.Equal(t, constants.DefaultSiteID, siteIDs[0], "Should be the default site_id")
	
	// Verify we have the expected number of events
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	require.NoError(t, err, "Should be able to count events")
	assert.Equal(t, len(requests), count, "Should have stored all successful requests")
}

func TestIntegration_StatsEndpoint_SingleSiteFiltering(t *testing.T) {
	srv, database := setupTestServer(t)
	defer database.Close()
	
	// Insert test data directly into database to simulate mixed data
	// (this tests that our query properly filters by site_id)
	insertSQL := `INSERT INTO events (anon_id, site_id, created_at, type, event, referrer, is_touch, browser_name, os_name, device_type, country, region)
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// Insert events for default site
	for i := 0; i < 3; i++ {
		_, err := database.Exec(insertSQL, 
			"hash-"+string(rune(i+'1')), constants.DefaultSiteID, 20250101, "pageview", 
			"https://example.com", "", "false", "Chrome", "Windows", "desktop", "US", "CA")
		require.NoError(t, err, "Should be able to insert test event %d", i+1)
	}
	
	// Insert events for hypothetical other site (to test filtering)
	for i := 0; i < 2; i++ {
		_, err := database.Exec(insertSQL, 
			"other-hash-"+string(rune(i+'1')), "other-site", 20250101, "pageview", 
			"https://other.com", "", "false", "Chrome", "Windows", "desktop", "US", "CA")
		require.NoError(t, err, "Should be able to insert other site event %d", i+1)
	}
	
	// Request stats
	req := httptest.NewRequest("GET", "/api/v1/stats/realtime", nil)
	rec := httptest.NewRecorder()
	
	srv.ServeHTTP(rec, req)
	
	assert.Equal(t, http.StatusOK, rec.Code, "Stats request should succeed")
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
	
	// The response should show only the 3 events for the default site, not all 5
	body := rec.Body.String()
	assert.Contains(t, body, "Total Pageviews", "Should contain pageviews label")
	assert.Contains(t, body, "3", "Should show count of 3 (default site events only)")
	
	// Verify the response shows exactly 3, proving the filtering works
	// (If filtering failed, it would show 5 total events)
	assert.NotContains(t, body, ">5<", "Should not show total count including other sites")
}

func TestIntegration_CORS_Headers(t *testing.T) {
	srv, database := setupTestServer(t)
	defer database.Close()
	
	req := httptest.NewRequest("GET", "/api/v1/collect?type=pageview", nil)
	rec := httptest.NewRecorder()
	
	srv.ServeHTTP(rec, req)
	
	// Verify CORS headers are present (from middleware)
	assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Methods"), "Should have CORS methods header")
	assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Headers"), "Should have CORS headers header")
}

func TestIntegration_Dashboard_Route(t *testing.T) {
	srv, database := setupTestServer(t)
	defer database.Close()
	
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	
	srv.ServeHTTP(rec, req)
	
	assert.Equal(t, http.StatusOK, rec.Code, "Dashboard should be accessible")
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
	
	body := rec.Body.String()
	assert.Contains(t, body, "Nyla Analytics", "Should contain page title")
	assert.Contains(t, body, "Dashboard", "Should contain dashboard content")
}
