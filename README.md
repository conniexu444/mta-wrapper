# MTA GTFS Wrapper in Go

The MTA API returns all of their data in GTFS format, making it annoyingly difficult to read as a developer. Additionally, this bunches up all of the different lines together since there are 7 endpoints total, and there are 36 subway lines total.

This API has several endpoints that you can hit:
**/feed?route={some subway line}**
This will pull up the entire data from each of the 7 endpoints in a json format. Example use case below:
/feed?route=B returns the following output:

````

returns JSON like:

```json
{
  "header": {
    "gtfs_realtime_version": "1.0",
    "timestamp": 1758503002
  },
  "entity": [
    {
      "id": "000001D",
      "trip_update": {
        "trip": {
          "trip_id": "121250_D..N07X002",
          "route_id": "D",
          "start_time": "20:12:28",
          "start_date": "20250921"
        }
      }
    },
    {
      "id": "000002D",
      "vehicle": {
        "trip": {
          "trip_id": "121250_D..N07X002",
          "route_id": "D",
          "start_time": "20:12:28",
          "start_date": "20250921"
        },
        "current_stop_sequence": 14,
        "stop_id": "R31N",
        "current_status": 1,
        "timestamp": 1758501774
      }
    },
    ...
  ]
}

````

### /arrivals?route={some_subway_line}

This will return all trains in all directions for a particular subway line for a given time period. The time period is defined in the config/constants.go file. It is the **ArrivalWindowMinutes** constant. Right now this is set to 20 minutes, but you can change it if you desire.

### /arrivals?route={some_subway_line}&station={some_subway_station}

This will return all trains in all directions for a particular subway station for a given time period. The given time period is 20 minutes in the config/constants.go file.
Examples:
/arrivals?route=L&station=bedford-av
This will give you all arrivals for the L train for the station Bedford Avenue

## How to run the API

### Prerequisites

- Go (https://go.dev/dl/)

```
git clone https://github.com/conniexu444/mta-wrapper.git
cd mta-wrapper

# Download dependencies
go mod tidy

# Run the server
go run main.go
```
