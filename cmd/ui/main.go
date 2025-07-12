package main

import (
	"net/http"
	"os"

	"github.com/joepurdy/nyla/pkg/handlers"
)

func main() {
	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://api.localhost"
	}

	h := &handlers.UIHandlers{APIBaseURL: apiBaseURL}
	http.HandleFunc("/", h.DashboardHandler)
	http.ListenAndServe(":8080", nil)
}
