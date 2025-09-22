package handlers

import (
	"net/http"
	"strings"

	"github.com/conniexu444/mta-wrapper/services"
)

func ArrivalsHandler(w http.ResponseWriter, r *http.Request) {
	route := strings.ToUpper(r.URL.Query().Get("route"))
	if route == "" {
		http.Error(w, "missing route parameter", http.StatusBadRequest)
		return
	}

	_, err := services.FetchFeed(route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
