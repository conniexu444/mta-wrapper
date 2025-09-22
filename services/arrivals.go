package services

import (
	"sync"
	"time"

	"github.com/conniexu444/mta-wrapper/config"
	"github.com/conniexu444/mta-wrapper/models"
)

var AllRoutes = []string{
	"1", "2", "3", "4", "5", "6", "S",
	"A", "C", "E",
	"N", "Q", "R", "W",
	"B", "D", "F", "M",
	"L", "G",
	"J", "Z",
}

func FetchArrivalsForStation(routes []string, stopIDs []string) []models.Arrival {
	now := time.Now().Unix()
	cutoff := time.Now().Add(config.ArrivalWindowMinutes * time.Minute).Unix()

	results := make(chan []models.Arrival, len(routes))
	var wg sync.WaitGroup

	for _, r := range routes {
		wg.Add(1)
		go func(rt string) {
			defer wg.Done()
			feed, err := FetchFeed(rt)
			if err != nil {
				return
			}

			localArrivals := []models.Arrival{}
			for _, entity := range feed.Entity {
				if entity.TripUpdate != nil {
					for _, stu := range entity.TripUpdate.StopTimeUpdate {
						arrTime := stu.Arrival.GetTime()
						if arrTime > now && arrTime <= cutoff && contains(stopIDs, stu.GetStopId()) {
							localArrivals = append(localArrivals, models.Arrival{
								Route:  rt,
								StopID: stu.GetStopId(),
								TripID: entity.TripUpdate.Trip.GetTripId(),
								Time:   time.Unix(arrTime, 0),
							})
						}
					}
				}
			}
			results <- localArrivals
		}(r)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var arrivals []models.Arrival
	for res := range results {
		arrivals = append(arrivals, res...)
	}
	return arrivals
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
