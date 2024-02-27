package statusLine

import (
	"errors"
	"fmt"
	app "github.com/rolandwarburton/ptv-status-line/internal"
	"strings"
	"time"
)

func RoutesAction(routeName string) ([]app.Route, error) {
	routes, _ := app.GetRoutes(routeName)

	// guard against no routes
	if len(routes) == 0 {
		return []app.Route{}, nil
	}

	return routes, nil
}

func StopsAction(stopName string, routeName string) ([]app.Stop, error) {
	routes, err := app.GetRoutes(routeName)
	if err != nil || len(routes) < 1 {
		return nil, fmt.Errorf("no route found for route %s", routeName)
	}
	route := routes[0]

	// get the stops
	stops, err := app.GetStops(route.RouteID, "", stopName)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to get routes")
	}

	return stops, nil
}

func DeparturesAction(routeName string, stopName string, directionName string, departuresCount int, timezone string) ([]app.Departure, error) {
	if stopName == "" || routeName == "" || directionName == "" {
		return nil, fmt.Errorf(
			"missing required information: "+
				"stopName=%q, "+
				"routeName=%q, "+
				"directionName=%q",
			stopName, routeName, directionName,
		)
	}

	routes, err := app.GetRoutes(routeName)
	if err != nil || len(routes) != 1 {
		if len(routes) > 1 {
			return nil, fmt.Errorf("too many routes returned for route \"%s\"", routeName)
		}
		return nil, fmt.Errorf("no route found for route %s", routeName)
	}
	route := routes[0]

	stops, err := app.GetStops(route.RouteID, "", stopName)
	if err != nil || len(stops) < 1 {
		return nil, fmt.Errorf("no route found for route %s", routeName)
	}
	stop := stops[0]

	// get the departures for a stop on a route
	departures, err := app.GetDepartures(stop.StopID, route.RouteID, "")
	if err != nil {
		return nil, errors.New("failed to get departures")
	}

	directions, err := app.GetDirections(route.RouteID)
	if err != nil || len(routes) < 1 {
		return nil, fmt.Errorf("no direction found for route %s", routeName)
	}

	// get the valid directions
	var validDirections []string
	var foundDirection *app.Direction

	// get all the directions as a string
	for _, direction := range directions {
		validDirections = append(validDirections, direction.DirectionName)
	}

	// look for the direction
	for _, direction := range directions {
		if strings.Contains(direction.DirectionName, directionName) {
			foundDirection = &direction
			break
		}
	}

	if foundDirection == nil {
		return nil, fmt.Errorf("no direction found for route %s. Valid directions are: %v", routeName, strings.Join(validDirections, ", "))
	}

	// get the next N departures in a certain direction
	departuresTowardsDirection, err := app.GetNextDepartureTowards(departures, foundDirection.DirectionID, departuresCount, timezone)
	if err != nil {
		return nil, errors.New("failed to get departures in specific direction")
	}

	nextDepartures := []app.Departure{}
	for i := 0; i < len(departuresTowardsDirection); i++ {
		if err == nil {
			layout := "2006-01-02T15:04:05Z"
			departureTime, err := time.Parse(layout, departuresTowardsDirection[i].ScheduledDepartureUTC)
			if err == nil {
				formattedTime := departureTime.Format("02-01-2006 03:04 PM")
				departuresTowardsDirection[i].ScheduledDepartureUTC = formattedTime
				nextDepartures = append(nextDepartures, departuresTowardsDirection[i])
			}
		}
	}

	return nextDepartures, nil
}

func DirectionsAction(routeName string) ([]app.Direction, error) {
	if routeName == "" {
		return nil, errors.New("route name not provided")
	}

	routes, err := app.GetRoutes(routeName)
	if err != nil || len(routes) < 1 {
		return nil, fmt.Errorf("no route found for route %s", routeName)
	}
	route := routes[0]

	directions, _ := app.GetDirections(route.RouteID)

	return directions, nil
}
