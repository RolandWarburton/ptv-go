package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	app "github.com/rolandwarburton/ptv-status-line/pkg"
	"github.com/urfave/cli/v2"
)

func writeToJSONFile(data any) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")

	// pretty print like so
	// fmt.Println(string(jsonData))

	// write the routes to a file
	file, _ := os.Create("routes.json")
	defer file.Close()
	file.Write(jsonData)
}

func printStops() {
	// get the departures for a stop on a route
	routes, err := app.GetStops(2, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	// pretty print like so
	jsonData, _ := json.MarshalIndent(routes, "", "  ")
	fmt.Println(string(jsonData))

	// write the routes to a file
	file, _ := os.Create("stops.json")
	defer file.Close()
	file.Write(jsonData)
}

func printRoutes() {
	// get the departures for a stop on a route
	routes, err := app.GetRoutes("Belgrave")
	if err != nil {
		fmt.Println(err)
		return
	}

	// pretty print like so
	jsonData, _ := json.MarshalIndent(routes, "", "  ")
	fmt.Println(string(jsonData))
}

func printNextTwoDepartures() {
	// get the departures for a stop on a route
	departures, err := app.GetDepartures(1016, 2, "?expand=All&include_geopath=true")
	if err != nil {
		fmt.Println(err)
		return
	}

	// get the next N departures in a certain direction
	nextBWDepartures, err := app.GetNextDepartureTowards(departures, 1, 2)
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

func main() {
	app := &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
