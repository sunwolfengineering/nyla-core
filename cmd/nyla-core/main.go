package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/sunwolfengineering/nyla-core/internal/server"
	"github.com/sunwolfengineering/nyla-core/internal/storage"
)

func main() {
	var port = flag.String("port", "8080", "port to listen on")
	flag.Parse()

	// Initialize database with built-in migrations
	db, err := storage.NewDB("nyla.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create unified server
	srv := server.New(db)

	fmt.Printf("🚀 nyla-core server starting on port %s\n", *port)
	fmt.Printf("📊 Dashboard: http://localhost:%s\n", *port)
	fmt.Printf("🔗 API: http://localhost:%s/api/v1\n", *port)

	log.Fatal(http.ListenAndServe(":"+*port, srv))
}
