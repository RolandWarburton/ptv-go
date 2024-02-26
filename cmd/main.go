package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	app "github.com/rolandwarburton/ptv-status-line/pkg"
	"github.com/urfave/cli/v2"
)

func routeAction(cCtx *cli.Context, format string, delimiter string) error {
	// get routes
	routes, _ := app.GetRoutes(cCtx.Args().First())

	// guard against no routes
	if len(routes) == 0 {
		fmt.Println("[]")
		return nil
	}

	// if not formatting print as JSON
	if !cCtx.IsSet("format") {
		jsonData, _ := json.MarshalIndent(routes, "", "  ")
		fmt.Println(string(jsonData))
		return nil
	}

	PrintFormatted[app.Route](routes, format, delimiter)
	return nil
}

func stopsAction(cCtx *cli.Context, routeName string, format string, delimiter string) error {
	// ensure a route ID is given
	stopName := cCtx.Args().First()
	if stopName == "" {
		return errors.New("please specify a stop name")
	}

	routes, err := app.GetRoutes(routeName)
	if err != nil || len(routes) < 1 {
		return fmt.Errorf("no route found for route %s", routeName)
	}
	route := routes[0]

	// get the stops
	stops, err := app.GetStops(route.RouteID, "", stopName)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to get routes")
	}

	// print as json if no formatting is given
	if format == "" {
		jsonData, _ := json.MarshalIndent(stops, "", "  ")
		fmt.Println(string(jsonData))
		return nil
	}

	PrintFormatted[app.Stop](stops, format, delimiter)

	return nil
}

func departuresAction(_ *cli.Context, routeID int, stopID int, direction int, departuresCount int, format string, delimiter string, timezone string) error {
	// get the departures for a stop on a route
	departures, err := app.GetDepartures(stopID, routeID, "")
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to get departures")
	}

	// get the next N departures in a certain direction
	departuresTowardsDirection, err := app.GetNextDepartureTowards(departures, direction, departuresCount, timezone)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to get departures in specific direction")
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

	if format == "" {
		PrettyPrint[app.Departure](departuresTowardsDirection, timezone)
		return nil
	}

	PrintFormatted[app.Departure](nextDepartures, format, delimiter)
	return nil
}

func directionsAction(cCtx *cli.Context, format string, delimiter string) error {
	routeName := cCtx.Args().First()
	if routeName == "" {
		return errors.New("route ID not provided")
	}

	routes, err := app.GetRoutes(routeName)
	if err != nil || len(routes) < 1 {
		return fmt.Errorf("no route found for route %s", routeName)
	}
	route := routes[0]

	directions, _ := app.GetDirections(route.RouteID)

	// print as json if no formatting is given
	if format == "" {
		PrettyPrint(directions, "Australia/Sydney")
		return nil
	}

	// format as a string
	PrintFormatted[app.Direction](directions, format, delimiter)
	return nil
}

func main() {
	var format string
	var delimiter string
	var departuresCount int
	var routeID int
	var routeName string
	var stopID int
	var directionID int
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
					return routeAction(c, format, delimiter)
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
					return stopsAction(c, routeName, format, delimiter)
				},
			},
			{
				Name:  "departures",
				Usage: "explore stops",
				Flags: append(flags,
					&cli.IntFlag{
						Name:        "count",
						Value:       -1,
						Usage:       "The next N trains departing",
						Destination: &departuresCount,
					},
					&cli.IntFlag{
						Name:        "route",
						Value:       -1,
						Usage:       "The route ID",
						Destination: &routeID,
					},
					&cli.IntFlag{
						Name:        "stop",
						Value:       -1,
						Usage:       "The stop ID",
						Destination: &stopID,
					},
					&cli.IntFlag{
						Name:        "direction",
						Value:       -1,
						Usage:       "The direction ID",
						Destination: &directionID,
					},
				),
				Action: func(c *cli.Context) error {
					return departuresAction(c, routeID, stopID, directionID, departuresCount, format, delimiter, timezone)
				},
			},
			{
				Name:  "directions",
				Usage: "explore directions",
				Flags: flags,
				Action: func(c *cli.Context) error {
					directionsAction(c, format, delimiter)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
