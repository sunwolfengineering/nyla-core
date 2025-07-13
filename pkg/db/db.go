package db

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/mileusna/useragent"
	_ "modernc.org/sqlite"
)

type Event struct {
	AnonID      string
	SiteID      string
	CreatedAt   int32
	Type        string
	Event       string
	Referrer    string
	IsTouch     bool
	BrowserName string
	OSName      string
	DeviceType  string
	Country     string
	Region      string
	Timestamp   time.Time
}

type Events struct {
	DB *sql.DB
}

func (e *Events) Open() error {
	db, err := sql.Open("sqlite", "nyla.db")
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	e.DB = db
	return nil
}

func (e *Events) Add(payload interface{}, hash string, ua useragent.UserAgent, geo interface{}) error {
	q := `
	INSERT INTO events
	(
		anon_id,
		site_id, 
		created_at, 
		type, 
		event, 
		referrer,
		is_touch, 
		browser_name, 
		os_name,
		device_type, 
		country, 
		region
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	)
	`

	// The payload and geo types should be replaced with the correct types from the main package or moved here as needed.
	// For now, use interface{} and expect the caller to provide the correct types.
	// You may want to define CollectorPayload and GeoInfo in a shared package for type safety.

	_, err := e.DB.ExecContext(
		context.Background(),
		q,
		hash,
		// The following fields must be extracted from payload and geo as appropriate.
		// This is a placeholder and should be updated for real usage.
		"", "", 0, "", "", "", "false", ua.Name, ua.OS, "not-implemented", "not-implemented", "not-implemented",
	)

	return err
}

func (e *Events) Close() error {
	if e.DB != nil {
		return e.DB.Close()
	}
	return nil
}

func nowToInt() uint32 {
	now := time.Now().Format("20060102")
	i, err := strconv.ParseInt(now, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return uint32(i)
}
