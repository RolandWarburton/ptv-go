package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	app "github.com/rolandwarburton/ptv-status-line/pkg"
	"github.com/urfave/cli/v2"
)

func routeAction(cCtx *cli.Context) ([]app.Route, error) {
	// get routes
	routes, _ := app.GetRoutes(cCtx.Args().First())

	// guard against no routes
	if len(routes) == 0 {
		return []app.Route{}, nil
	}

	return routes, nil
}

func stopsAction(cCtx *cli.Context, routeName string, format string, delimiter string) ([]app.Stop, error) {
	// ensure a route ID is given
	stopName := cCtx.Args().First()

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

func departuresAction(_ *cli.Context, routeName string, stopName string, directionName string, departuresCount int, format string, delimiter string, timezone string) ([]app.Departure, error) {
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

func directionsAction(cCtx *cli.Context, format string, delimiter string) ([]app.Direction, error) {
	routeName := cCtx.Args().First()
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

func main() {
	var format string
	var delimiter string
	var departuresCount int
	var routeName string
	var stopName string
	var directionName string
	var timezone string

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "format",
			Value:       "",
			Usage:       "format the output",
			Destination: &format,
		},
		&cli.StringFlag{
			Name:        "delimiter",
			Value:       " ",
			Usage:       "delimiter between format arguments",
			Destination: &delimiter,
		},
		&cli.StringFlag{
			Name:        "timezone",
			Value:       "Australia/Sydney",
			Usage:       "specify timezone for dates",
			Destination: &timezone,
		},
	}

	app := &cli.App{
		Name:                 "ptv-status-line",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "routes",
				Usage: "explore routes",
				Flags: flags,
				Action: func(c *cli.Context) error {
					routes, err := routeAction(c)
					PrintResult[app.Route](routes, format, delimiter, "Australia/Sydney")
					if err != nil {
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
			{
				Name:  "stops",
				Usage: "explore stops",
				Flags: append(flags, &cli.StringFlag{
					Name:        "route",
					Value:       "",
					Usage:       "Specify the route the station is on",
					Destination: &routeName,
				}),
				Action: func(c *cli.Context) error {
					stops, err := stopsAction(c, routeName, format, delimiter)
					PrintResult[app.Stop](stops, format, delimiter, "Australia/Sydney")
					if err != nil {
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
			{
				Name:  "departures",
				Usage: "explore stops",
				Flags: append(flags,
					&cli.IntFlag{
						Name:        "count",
						Value:       1,
						Usage:       "The next N trains departing",
						Destination: &departuresCount,
					},
					&cli.StringFlag{
						Name:        "route",
						Value:       "",
						Usage:       "The route ID",
						Destination: &routeName,
					},
					&cli.StringFlag{
						Name:        "stop",
						Value:       "",
						Usage:       "The stop ID",
						Destination: &stopName,
					},
					&cli.StringFlag{
						Name:        "direction",
						Value:       "",
						Usage:       "The direction ID",
						Destination: &directionName,
					},
				),
				Action: func(c *cli.Context) error {
					departures, err := departuresAction(c, routeName, stopName, directionName, departuresCount, format, delimiter, timezone)
					PrintResult[app.Departure](departures, format, delimiter, "Australia/Sydney")
					if err != nil {
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
			{
				Name:  "directions",
				Usage: "explore directions",
				Flags: flags,
				Action: func(c *cli.Context) error {
					directions, err := directionsAction(c, format, delimiter)
					PrintResult[app.Direction](directions, format, delimiter, "Australia/Sydney")
					if err != nil {
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
