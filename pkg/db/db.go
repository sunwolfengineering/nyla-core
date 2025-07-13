package db

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/mileusna/useragent"
	_ "modernc.org/sqlite"

	"github.com/sunwolfengineering/nyla-core/pkg/constants"
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

	// Always use the default site ID for single-site architecture
	_, err := e.DB.ExecContext(
		context.Background(),
		q,
		hash,
		constants.DefaultSiteID, // Always use default site
		nowToInt(),
		"pageview", // Default type for now - should be extracted from payload
		"", // URL - should be extracted from payload
		"", // Referrer - should be extracted from payload
		"false", // is_touch
		ua.Name, 
		ua.OS, 
		"not-implemented", // device_type
		"not-implemented", // country
		"not-implemented", // region
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
