package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/mileusna/useragent"

	"github.com/sunwolfengineering/nyla-core/internal/storage"
	"github.com/sunwolfengineering/nyla-core/pkg/constants"
	"github.com/sunwolfengineering/nyla-core/pkg/geo"
	"github.com/sunwolfengineering/nyla-core/pkg/hash"
)

type CollectorData struct {
	Type      string `json:"type"`
	Event     string `json:"event"`
	UserAgent string `json:"ua"`
	Hostname  string `json:"hostname"`
	Referrer  string `json:"referrer"`
}

type CollectorPayload struct {
	Data CollectorData `json:"data"`
}

type Handlers struct {
	DB *storage.DB
}

func (h *Handlers) GetCollectV1(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	siteID := r.URL.Query().Get("site_id")
	eventType := r.URL.Query().Get("type")
	url := r.URL.Query().Get("url")
	referrer := r.URL.Query().Get("referrer")

	// Enforce single-site architecture
	if siteID != "" && siteID != constants.DefaultSiteID {
		response := map[string]string{
			"error": "Invalid site_id. This instance only supports site_id='default'",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set defaults if not provided
	if eventType == "" {
		eventType = "pageview"
	}
	if url == "" {
		url = r.Header.Get("Referer")
	}
	if referrer == "" {
		referrer = r.Header.Get("Referer")
	}

	// Parse user agent for metadata
	ua := useragent.Parse(r.UserAgent())
	ip, _ := geo.IPFromRequest([]string{"X-Forwarded-For", "X-Real-IP"}, r)
	sessionID, _ := hash.GeneratePrivateIDHash(ip.String(), r.UserAgent(), r.Host, constants.DefaultSiteID)

	// Create event using new storage API
	event := &storage.Event{
		Type:      eventType,
		Timestamp: time.Now(),
		URL:       url,
		Referrer:  referrer,
		SessionID: sessionID,
		Metadata: map[string]interface{}{
			"user_agent":   r.UserAgent(),
			"hostname":     r.Host,
			"browser_name": ua.Name,
			"os_name":      ua.OS,
			"is_bot":       ua.Bot,
		},
	}

	ctx := context.Background()
	if err := h.DB.InsertEvent(ctx, event); err != nil {
		log.Printf("Error inserting event: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Event added: %s %s", eventType, url)
	
	// Return 1x1 transparent GIF
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
	ctx := context.Background()
	stats, err := h.DB.GetRealtimeStats(ctx)
	if err != nil {
		log.Printf("Error getting realtime stats: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<div class='error'>Error: %v</div>", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fragment := elem.Div(attrs.Props{"class": "bg-white rounded-lg shadow p-6"},
		elem.Div(attrs.Props{"class": "text-sm text-gray-500"}, elem.Text("Total Pageviews")),
		elem.Div(attrs.Props{"class": "text-2xl font-bold text-indigo-700 mt-2"}, 
			elem.Text(fmt.Sprintf("%d", stats.PageviewsToday))),
	).Render()
	w.Write([]byte(fragment))
}
