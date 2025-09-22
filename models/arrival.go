package models

import (
	"time"
)

type Arrival struct {
	Route  string    `json:"route"`
	StopID string    `json:"stop_id"`
	TripID string    `json:"trip_id"`
	Time   time.Time `json:"time"`
}
