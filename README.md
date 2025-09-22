# MTA GTFS Wrapper in Go

The MTA API returns all of their data in GTFS format, making it annoyingly difficult to read as a developer. Additionally, this bunches up all of the different lines together since there are 7 endpoints total, and there are 36 subway lines total.

This API has several endpoints that you can hit:
**/feed?route={some subway line}**
This will pull up the entire data from each of the 7 endpoints in a json format. Example use case below:
/feed?route=B returns all of the feed data associated with route B.

### /arrivals?route={some_subway_line}

This will return all trains in all directions for a particular subway line for a given time period. The time period is defined in the config/constants.go file. It is the **ArrivalWindowMinutes** constant. Right now this is set to 20 minutes, but you can change it if you desire.

### /arrivals?route={some_subway_line}&station={some_subway_station}

This will return all trains in all directions for a particular subway station for a given time period. The given time period is 20 minutes in the config/constants.go file.
Examples:
/arrivals?route=L&station=bedford-av
This will give you all arrivals for the L train for the station Bedford Avenue

### /arrivals?route=ALL&station={some_subway_station}

This will return all of the trains in the 20 minute window we specified in the configuration file. This will return all of the trains that stop in both directions for a station.
Examples:
/arrivals?route=ALL&station=times-sq-42-st

### /arrivals?route={a_route_or_ALL}&station={some_subway_station}&direction={N_or_S}

This will return all of the trains in the 20 minute window that we specified in the configuration file going either North or South. If you do route=ALL, it will return all N or S bound trains for a particular station.

### /arrivals?route={a_route_or_ALL}&direction={N_or_S}

This will return for any route you choose all northbound or southbound trains.

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
