package services

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	gtfs "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/conniexu444/mta-wrapper/config"
	"google.golang.org/protobuf/proto"
)

func buildFeedURL(route string) (string, error) {
	key := strings.ToUpper(route)
	feedPath, ok := config.RouteFeeds[key]
	if !ok {
		return "", fmt.Errorf("unknown route: %s", route)
	}
	return fmt.Sprintf("https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/%s", feedPath), nil
}

func FetchFeed(route string) (*gtfs.FeedMessage, error) {
	key := strings.ToUpper(route)
	_, ok := config.RouteFeeds[route]
	if !ok {
		return nil, fmt.Errorf("unknown route: %s", route)
	}
	url, _ := buildFeedURL(route)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, err
	}

	filtered := gtfs.FeedMessage{
		Header: feed.Header,
	}
	for _, entity := range feed.Entity {
		if entity.TripUpdate != nil && entity.TripUpdate.Trip != nil {
			if entity.TripUpdate.Trip.RouteId != nil && *entity.TripUpdate.Trip.RouteId == key {
				filtered.Entity = append(filtered.Entity, entity)
			}
		}
	}

	return feed, nil
}
