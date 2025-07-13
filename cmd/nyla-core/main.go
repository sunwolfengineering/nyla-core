package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/sunwolfengineering/nyla-core/internal/server"
	"github.com/sunwolfengineering/nyla-core/pkg/db"
)

func main() {
	var port = flag.String("port", "8080", "port to listen on")
	flag.Parse()

	// Initialize database
	events := &db.Events{}
	if err := events.Open(); err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer events.Close()

	// Create unified server
	srv := server.New(events)

	fmt.Printf("🚀 nyla-core server starting on port %s\n", *port)
	fmt.Printf("📊 Dashboard: http://localhost:%s\n", *port)
	fmt.Printf("🔗 API: http://localhost:%s/api/v1\n", *port)

	log.Fatal(http.ListenAndServe(":"+*port, srv))
}
