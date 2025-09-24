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
	"github.com/conniexu444/mta-wrapper/utils"
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
			suggestions := utils.SuggestStations(station, config.StationStops, 3)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":        "unknown station: " + station,
				"did_you_mean": suggestions,
			})
			return
		}
		if direction == "N" || direction == "S" {
			for _, id := range stopIDs {
				if strings.HasSuffix(id, direction) {
					stationStopIDs = append(stationStopIDs, id)
				}
			}
		} else {
			stationStopIDs = stopIDs
		}
	}

	if route == "ALL" && station != "" {
		arrivals := services.FetchArrivalsForStation(services.AllRoutes, stationStopIDs, direction)
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
					stopID := stu.GetStopId()

					if station != "" && !contains(stationStopIDs, stopID) {
						continue
					}

					if (direction == "N" || direction == "S") && !strings.HasSuffix(stopID, direction) {
						continue
					}

					arrivals = append(arrivals, models.Arrival{
						Route:  route,
						StopID: stopID,
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
