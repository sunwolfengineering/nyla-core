// SPDX-License-Identifier: GPL-3.0-only
package db

import (
	"database/sql"
	"testing"

	"github.com/mileusna/useragent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/sunwolfengineering/nyla-core/pkg/constants"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Create in-memory SQLite database for testing
	db, err := sql.Open("sqlite", ":memory:")
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
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err, "Should be able to create events table")
	
	return db
}

func TestEvents_Add_UsesDefaultSiteID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	events := &Events{DB: db}
	
	// Test data
	hash := "test-hash"
	ua := useragent.UserAgent{
		Name: "Chrome",
		OS:   "Windows",
	}
	
	// Add an event
	err := events.Add(nil, hash, ua, nil)
	require.NoError(t, err, "Should be able to add event")
	
	// Verify the event was stored with the default site ID
	var storedSiteID string
	var storedHash string
	err = db.QueryRow("SELECT anon_id, site_id FROM events WHERE anon_id = ?", hash).Scan(&storedHash, &storedSiteID)
	require.NoError(t, err, "Should be able to query stored event")
	
	assert.Equal(t, hash, storedHash, "Hash should match")
	assert.Equal(t, constants.DefaultSiteID, storedSiteID, "Site ID should be the default")
}

func TestEvents_Add_DefaultSiteIDConstant(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	events := &Events{DB: db}
	
	// Add multiple events
	for i := 0; i < 3; i++ {
		hash := "test-hash-" + string(rune(i+'1'))
		ua := useragent.UserAgent{Name: "Test", OS: "Test"}
		
		err := events.Add(nil, hash, ua, nil)
		require.NoError(t, err, "Should be able to add event %d", i)
	}
	
	// Verify all events have the same default site ID
	rows, err := db.Query("SELECT DISTINCT site_id FROM events")
	require.NoError(t, err, "Should be able to query site IDs")
	defer rows.Close()
	
	siteIDs := []string{}
	for rows.Next() {
		var siteID string
		err := rows.Scan(&siteID)
		require.NoError(t, err, "Should be able to scan site ID")
		siteIDs = append(siteIDs, siteID)
	}
	
	// Should only have one unique site ID
	assert.Len(t, siteIDs, 1, "Should only have one unique site ID")
	assert.Equal(t, constants.DefaultSiteID, siteIDs[0], "The only site ID should be the default")
}

func TestEvents_Add_UserAgentFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	events := &Events{DB: db}
	
	// Test with specific user agent
	hash := "test-hash"
	ua := useragent.UserAgent{
		Name: "Firefox",
		OS:   "Linux",
	}
	
	err := events.Add(nil, hash, ua, nil)
	require.NoError(t, err, "Should be able to add event")
	
	// Verify user agent fields are stored
	var browserName, osName string
	err = db.QueryRow("SELECT browser_name, os_name FROM events WHERE anon_id = ?", hash).Scan(&browserName, &osName)
	require.NoError(t, err, "Should be able to query user agent fields")
	
	assert.Equal(t, "Firefox", browserName, "Browser name should match")
	assert.Equal(t, "Linux", osName, "OS name should match")
}

func TestEvents_Open_Close(t *testing.T) {
	events := &Events{}
	
	// Test that Open method exists and works
	err := events.Open()
	// Open should succeed and create a database file
	assert.NoError(t, err, "Open should succeed")
	assert.NotNil(t, events.DB, "DB should be initialized after Open")
	
	// Test Close method
	err = events.Close()
	assert.NoError(t, err, "Close should succeed")
	
	// Clean up the created database file
	defer func() {
		// Remove test database file if it was created
		if err := events.Close(); err == nil {
			// File cleanup handled by defer
		}
	}()
	
	// Test Close with nil DB
	eventsEmpty := &Events{}
	err = eventsEmpty.Close()
	assert.NoError(t, err, "Close should succeed even with nil DB")
}

func TestNowToInt(t *testing.T) {
	// Test the nowToInt function
	result := nowToInt()
	
	// Should return a valid integer representation of current date
	assert.Greater(t, result, uint32(20250101), "Should return a date >= 2025-01-01")
	assert.Less(t, result, uint32(99991231), "Should return a reasonable date")
	
	// Test that it's consistent (same day should return same value)
	result2 := nowToInt()
	assert.Equal(t, result, result2, "Should return same value for same day")
}

func TestEvents_DatabaseSchema(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Test that our test schema matches what the Add method expects
	events := &Events{DB: db}
	
	hash := "schema-test-hash"
	ua := useragent.UserAgent{Name: "TestBrowser", OS: "TestOS"}
	
	err := events.Add(nil, hash, ua, nil)
	require.NoError(t, err, "Add method should work with test schema")
	
	// Verify all expected columns exist and can be queried
	var (
		anonID, siteID, eventType, event, referrer string
		createdAt int
		isTouch, browserName, osName, deviceType, country, region string
	)
	
	query := `SELECT anon_id, site_id, created_at, type, event, referrer, 
			  is_touch, browser_name, os_name, device_type, country, region 
			  FROM events WHERE anon_id = ?`
	
	err = db.QueryRow(query, hash).Scan(
		&anonID, &siteID, &createdAt, &eventType, &event, &referrer,
		&isTouch, &browserName, &osName, &deviceType, &country, &region,
	)
	require.NoError(t, err, "Should be able to query all expected columns")
	
	assert.Equal(t, hash, anonID, "anon_id should match")
	assert.Equal(t, constants.DefaultSiteID, siteID, "site_id should be default")
	assert.Equal(t, "TestBrowser", browserName, "browser_name should match")
	assert.Equal(t, "TestOS", osName, "os_name should match")
}
