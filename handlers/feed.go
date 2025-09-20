package handlers

import (
	"encoding/json"
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

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	route := r.URL.Query().Get("route")
	if route == "" {
		http.Error(w, "Missing route parameter", http.StatusBadRequest)
		return
	}

	url, err := buildFeedURL(route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "failed to fetch MTA feed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		http.Error(w, "failed to decode protobuf: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
