package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conniexu444/mta-wrapper/services"
)

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	route := r.URL.Query().Get("route")
	if route == "" {
		http.Error(w, "Missing route parameter", http.StatusBadRequest)
		return
	}

	feed, err := services.FetchFeed(route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
