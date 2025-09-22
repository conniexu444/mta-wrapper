package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/conniexu444/mta-wrapper/config"
	"github.com/conniexu444/mta-wrapper/models"
	"github.com/conniexu444/mta-wrapper/services"
)

func ArrivalsHandler(w http.ResponseWriter, r *http.Request) {
	route := strings.ToUpper(r.URL.Query().Get("route"))
	if route == "" {
		http.Error(w, "missing route parameter", http.StatusBadRequest)
		return
	}

	feed, err := services.FetchFeed(route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	arrivals := []models.Arrival{}
	now := time.Now().Unix()
	cutoff := time.Now().Add(config.ArrivalWindowMinutes * time.Minute).Unix()

	for _, entity := range feed.Entity {
		if entity.TripUpdate != nil {
			for _, stu := range entity.TripUpdate.StopTimeUpdate {
				arrTime := stu.Arrival.GetTime()
				if arrTime > now && arrTime <= cutoff {
					arrivals = append(arrivals, models.Arrival{
						Route:  route,
						StopID: stu.GetStopId(),
						TripID: entity.TripUpdate.Trip.GetTripId(),
						Time:   time.Unix(arrTime, 0),
					})
				}
			}
		}
	}

	sort.Slice(arrivals, func(i, j int) bool {
		return arrivals[i].Time.Before(arrivals[j].Time)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(arrivals)
}
