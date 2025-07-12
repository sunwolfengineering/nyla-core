package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/mileusna/useragent"

	"github.com/joepurdy/nyla/pkg/db"
	"github.com/joepurdy/nyla/pkg/geo"
	"github.com/joepurdy/nyla/pkg/hash"
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

type Handlers struct {
	Events *db.Events
}

func (h *Handlers) GetCollectV1(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	siteID := r.URL.Query().Get("site_id")
	eventType := r.URL.Query().Get("type")
	url := r.URL.Query().Get("url")
	referrer := r.URL.Query().Get("referrer")

	if siteID == "" || eventType == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	payload := CollectorPayload{
		SiteID: siteID,
		Data: CollectorData{
			Type:      eventType,
			Event:     url,
			UserAgent: r.UserAgent(),
			Hostname:  r.Host,
			Referrer:  referrer,
		},
	}

	ua := useragent.Parse(r.UserAgent())
	ip, _ := geo.IPFromRequest([]string{"X-Forwarded-For", "X-Real-IP"}, r)
	hashVal, _ := hash.GeneratePrivateIDHash(ip.String(), r.UserAgent(), r.Host, siteID)

	if err := h.Events.Add(payload, hashVal, ua, nil); err != nil {
		log.Println("error adding event:", err)
	}

	log.Println("event added", payload)
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00,
		0x01, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0x21, 0xF9, 0x04, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
		0x01, 0x00, 0x3B,
	})
}

func (h *Handlers) GetStatsRealtimeV1(w http.ResponseWriter, r *http.Request) {
	db := h.Events.DB
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
