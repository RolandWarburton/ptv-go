package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	app "github.com/rolandwarburton/ptv-status-line/pkg"
	"github.com/urfave/cli/v2"
)

func prettyPrint(data any) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}

func writeToJSONFile(data any) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")

	// write the routes to a file
	file, _ := os.Create("routes.json")
	defer file.Close()
	file.Write(jsonData)
}

func printNextDepartures(totalDepartures int) {
	// get the departures for a stop on a route
	departures, err := app.GetDepartures(1016, 2, "?expand=All&include_geopath=true")
	if err != nil {
		fmt.Println(err)
		return
	}

	// get the next N departures in a certain direction
	nextBWDepartures, err := app.GetNextDepartureTowards(departures, 1, totalDepartures)
	if err != nil {
		fmt.Println(err)
		return
	}

	nextDepartures := []string{}
	for i := 0; i < len(nextBWDepartures); i++ {
		if err == nil {
			layout := "2006-01-02T15:04:05Z"
			departureTime, err := time.Parse(layout, nextBWDepartures[i].ScheduledDepartureUTC)
			if err != nil {
				nextDepartures = append(nextDepartures, "ERROR")
			} else {
				formattedTime := departureTime.Format("2-1-2006 3:4 PM")
				nextDepartures = append(nextDepartures, formattedTime)
			}
		}
	}
	fmt.Println(nextDepartures)
}

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

	// format as a string
	// Example --format "RouteID RouteName"
	route := routes[0]
	formatArgs := strings.Split(format, " ")
	result := ""
	for i, arg := range formatArgs {
		// dynamically access the fields of the Route
		val := reflect.ValueOf(route)
		field := val.FieldByName(arg)
		if field.IsValid() && i < len(formatArgs)-1 {
			result += fmt.Sprintf("%v%s", field.Interface(), delimiter)
		} else {
			result += fmt.Sprintf("%v", field.Interface())
		}
	}
	fmt.Println(result)
	return nil
}

func stopsAction(cCtx *cli.Context, stopName string, format string, delimiter string) error {
	// ensure a route ID is given
	stopID := cCtx.Args().First()
	if stopID == "" {
		return errors.New("please specify a route ID")
	}
	var v int
	var err error
	if v, err = strconv.Atoi(stopID); err != nil {
		return errors.New("please specify a valid route ID number")
	}

	// get the stops
	stops, err := app.GetStops(v, "", stopName)
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

	// format as a string
	stop := stops[0]
	formatArgs := strings.Split(format, " ")
	result := ""
	for i, arg := range formatArgs {
		// dynamically access the fields of the Route
		val := reflect.ValueOf(stop)
		field := val.FieldByName(arg)
		if field.IsValid() && i < len(formatArgs)-1 {
			result += fmt.Sprintf("%v%s", field.Interface(), delimiter)
		} else {
			result += fmt.Sprintf("%v", field.Interface())
		}
	}
	fmt.Println(result)

	return nil
}

func departuresAction(_ *cli.Context, routeID int, stopID int, direction int, departuresCount int, format string, delimiter string) error {
	// get the departures for a stop on a route
	departures, err := app.GetDepartures(stopID, routeID, "")
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to get departures")
	}

	// get the next N departures in a certain direction
	departuresTowardsDirection, err := app.GetNextDepartureTowards(departures, direction, departuresCount)
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
		prettyPrint(departuresTowardsDirection)
		return nil
	}

	for _, departure := range nextDepartures {
		formatArgs := strings.Split(format, " ")
		result := ""
		for i, arg := range formatArgs {
			// dynamically access the fields of the Route
			val := reflect.ValueOf(departure)
			field := val.FieldByName(arg)
			if field.IsValid() && i < len(formatArgs)-1 {
				result += fmt.Sprintf("%v%s", field.Interface(), delimiter)
			} else {
				result += fmt.Sprintf("%v", field.Interface())
			}
		}
		fmt.Println(result)
	}
	return nil
}

func directionsAction(cCtx *cli.Context, format string, delimiter string) error {
	arg1 := cCtx.Args().First()
	if arg1 == "" {
		return errors.New("route ID not provided")
	}
	var routeID int
	var err error
	if routeID, err = strconv.Atoi(arg1); err != nil {
		return errors.New("failed to parse route ID")
	}
	directions, _ := app.GetDirections(routeID)
	// print as json if no formatting is given
	if format == "" {
		prettyPrint(directions)
		return nil
	}

	// format as a string
	for _, direction := range directions {
		formatArgs := strings.Split(format, " ")
		result := ""
		for i, arg := range formatArgs {
			// dynamically access the fields of the Route
			val := reflect.ValueOf(direction)
			field := val.FieldByName(arg)
			if field.IsValid() && i < len(formatArgs)-1 {
				result += fmt.Sprintf("%v%s", field.Interface(), delimiter)
			} else {
				result += fmt.Sprintf("%v", field.Interface())
			}
		}
		fmt.Println(result)
	}
	return nil
}

func main() {
	var format string
	var delimiter string
	var stopName string
	var departuresCount int
	var routeID int
	var stopID int
	var directionID int

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
					Name:        "stop",
					Value:       "",
					Usage:       "Filter a specific stop by the station name",
					Destination: &stopName,
				}),
				Action: func(c *cli.Context) error {
					return stopsAction(c, stopName, format, delimiter)
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
					return departuresAction(c, routeID, stopID, directionID, departuresCount, format, delimiter)
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
