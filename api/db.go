package main

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

func (e *Events) Add(payload CollectorPayload, hash string, ua useragent.UserAgent, geo *GeoInfo) error {
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

	_, err := e.DB.ExecContext(
		context.Background(),
		q,
		hash,
		payload.SiteID,
		nowToInt(),
		payload.Data.Type,
		payload.Data.Event,
		payload.Data.Referrer,
		"false",
		ua.Name,
		ua.OS,
		"not-implemented",
		geo.Country,
		geo.RegionName,
	)

	return err
}

func nowToInt() uint32 {
	now := time.Now().Format("20060102")
	i, err := strconv.ParseInt(now, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return uint32(i)
}
