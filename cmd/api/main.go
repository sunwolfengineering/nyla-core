// nyla-api - GDPR compliant privacy focused web analytics
// Copyright (C) 2024 Joe Purdy
// mailto:nyla AT purdy DOT dev
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/mileusna/useragent"

	"github.com/joepurdy/nyla/pkg/db"
	"github.com/joepurdy/nyla/pkg/geo"
	"github.com/joepurdy/nyla/pkg/handlers"
	"github.com/joepurdy/nyla/pkg/hash"
)

// Version is provided at compile time
var Version = "devel"

var (
	events *db.Events = &db.Events{}

	// CORS configuration via environment variables
	corsAllowedOrigins   = getEnvDefault("CORS_ALLOWED_ORIGINS", "https://localhost")
	corsAllowedHeaders   = getEnvDefault("CORS_ALLOWED_HEADERS", "Content-Type,HX-Request,HX-Target,HX-Current-URL,HX-Trigger,HX-Trigger-Name,HX-History-Restore-Request")
	corsExposedHeaders   = getEnvDefault("CORS_EXPOSED_HEADERS", "HX-Redirect,HX-Location,HX-Push,HX-Refresh,HX-Trigger,HX-Trigger-After-Settle,HX-Trigger-After-Swap")
	corsAllowCredentials = getEnvDefault("CORS_ALLOW_CREDENTIALS", "true")
)

// getEnvDefault returns the value of the environment variable or a default
func getEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// corsMiddleware adds CORS headers and handles preflight requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && (corsAllowedOrigins == "*" || strings.Contains(corsAllowedOrigins, origin)) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		} else if corsAllowedOrigins == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", corsAllowedHeaders)
		w.Header().Set("Access-Control-Expose-Headers", corsExposedHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", corsAllowCredentials)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type CollectorData struct {
	Type      string `json:"type"`
	Event     string `json:"event"`
	UserAgent string `json:"ua"`
	Hostname  string `json:"hostname"`
	Referrer  string `json:"referrer"`
}

type CollectorPayload struct {
	SiteID string        `json:"site_id"`
	Data   CollectorData `json:"data"`
}

func main() {
	fmt.Println("nyla version:", Version)

	if err := events.Open(); err != nil {
		log.Fatal(err)
	}

	h := &handlers.Handlers{Events: events}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/collect", h.GetCollectV1)
	mux.HandleFunc("GET /v1/stats/realtime", h.GetStatsRealtimeV1)

	// Wrap with CORS middleware
	handler := corsMiddleware(mux)

	fmt.Println("listening on :9876")
	http.ListenAndServe(":9876", handler)
}

// getCollectV1 handles GET /v1/collect for single event collection via query parameters.
func getCollectV1(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	siteID := r.URL.Query().Get("site_id")
	eventType := r.URL.Query().Get("type")
	url := r.URL.Query().Get("url")
	referrer := r.URL.Query().Get("referrer")

	// Minimal required fields
	if siteID == "" || eventType == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Build CollectorPayload (reuse legacy struct)
	payload := CollectorPayload{
		SiteID: siteID,
		Data: CollectorData{
			Type:      eventType,
			Event:     url, // For MVP, store url in Event field
			UserAgent: r.UserAgent(),
			Hostname:  r.Host,
			Referrer:  referrer,
		},
	}

	ua := useragent.Parse(r.UserAgent())

	ip, _ := geo.IPFromRequest([]string{"X-Forwarded-For", "X-Real-IP"}, r)
	// geoInfo, _ := geo.GetGeoInfo(ip.String()) // TODO: add geo info
	hashVal, _ := hash.GeneratePrivateIDHash(ip.String(), r.UserAgent(), r.Host, siteID)

	// Store event
	if err := events.Add(payload, hashVal, ua, nil); err != nil {
		log.Println("error adding event:", err) // TODO: error handling strategy is not yet defined
	}

	log.Println("event added", payload) // TODO: better logging strategy is needed
	// Return 1x1 transparent GIF
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{ // 1x1 transparent GIF
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00,
		0x01, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0x21, 0xF9, 0x04, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
		0x01, 0x00, 0x3B,
	})
}

// getStatsRealtimeV1 handles GET /v1/stats/realtime for realtime stats.
func getStatsRealtimeV1(w http.ResponseWriter, r *http.Request) {
	// Query total pageviews
	db := events.DB
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM events WHERE type = 'pageview'").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<div class='error'>Error: %v</div>", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fragment := elem.Div(attrs.Props{"class": "bg-white rounded-lg shadow p-6"},
		elem.Div(attrs.Props{"class": "text-sm text-gray-500"}, elem.Text("Total Pageviews")),
		elem.Div(attrs.Props{"class": "text-2xl font-bold text-indigo-700 mt-2"}, elem.Text(fmt.Sprintf("%d", count))),
	).Render()
	w.Write([]byte(fragment))
}
