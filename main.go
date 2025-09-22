package main

import (
	"fmt"
	"net/http"

	"github.com/conniexu444/mta-wrapper/handlers"
)

func main() {
	// Register endpoints
	http.HandleFunc("/feed", handlers.FeedHandler)
	http.HandleFunc("/arrivals", handlers.ArrivalsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server")
	}
}
