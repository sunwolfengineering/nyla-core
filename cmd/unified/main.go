package main

import (
	"fmt"
	"log"

	"github.com/joepurdy/nyla/internal/server"
	"github.com/joepurdy/nyla/pkg/db"
)

// Version is provided at compile time
var Version = "devel"

func main() {
	fmt.Println("nyla unified server version:", Version)

	events := &db.Events{}
	if err := events.Open(); err != nil {
		log.Fatal(err)
	}

	srv := server.New(events)
	
	fmt.Println("unified server listening on :8080")
	fmt.Println("API routes available at /api/v1/*")
	fmt.Println("UI dashboard available at /")
	
	if err := srv.ListenAndServe(":8080"); err != nil {
		log.Fatal(err)
	}
}
