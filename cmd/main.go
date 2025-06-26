package main

import (
	"fmt"
	"log"
	"media-service/internal/database"
	"net/http"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	log.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
