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
	station := r.URL.Query().Get("station")
	direction := strings.ToUpper(r.URL.Query().Get("direction"))

	if route == "" {
		http.Error(w, "missing route parameter", http.StatusBadRequest)
		return
	}

	var stationStopIDs []string
	if station != "" {
		stopIDs, ok := config.StationStops[station]
		if !ok {
			http.Error(w, "unknown station: "+station, http.StatusBadRequest)
			return
		}

		if direction == "N" || direction == "S" {
			filtered := []string{}
			for _, id := range stopIDs {
				if strings.HasSuffix(id, direction) {
					filtered = append(filtered, id)
				}
			}
			stationStopIDs = filtered
		} else {
			stationStopIDs = stopIDs
		}
	}

	if route == "ALL" && station != "" {
		stopIDs, ok := config.StationStops[station]
		if !ok {
			http.Error(w, "unknown station: "+station, http.StatusBadRequest)
			return
		}

		arrivals := services.FetchArrivalsForStation(services.AllRoutes, stopIDs)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(arrivals)
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

					if station != "" && !contains(stationStopIDs, stu.GetStopId()) {
						continue
					}

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

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
