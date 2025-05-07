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
	"fmt"
	"log"
	"net/http"

	"github.com/mileusna/useragent"
)

// Version is provided at compile time
var Version = "devel"

var (
	events *Events = &Events{}
)

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

	http.HandleFunc("GET /v1/collect", getCollectV1)

	fmt.Println("listening on :9876")
	http.ListenAndServe(":9876", nil)
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

	ip, _ := ipFromRequest([]string{"X-Forwarded-For", "X-Real-IP"}, r)
	// geoInfo, _ := getGeoInfo(ip.String()) // TODO: add geo info
	hash, _ := generatePrivateIDHash(ip.String(), r.UserAgent(), r.Host, siteID)

	// Store event
	if err := events.Add(payload, hash, ua, nil); err != nil {
		fmt.Println("error adding event:", err) // TODO: error handling strategy is not yet defined
	}

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
