package main

import (
	"fmt"
	"log"
	"os"

	statusLine "github.com/rolandwarburton/ptv-go/pkg"
	"github.com/urfave/cli/v2"
)

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
		Name:                 "ptv-go",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "routes",
				Usage: "explore routes",
				Flags: flags,
				Action: func(c *cli.Context) error {
					routeName := c.Args().First()
					routes, err := statusLine.RoutesAction(routeName)
					PrintResult[statusLine.Route](routes, format, delimiter, "Australia/Sydney")
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
					stopName := c.Args().First()
					stops, err := statusLine.StopsAction(stopName, routeName)
					PrintResult[statusLine.Stop](stops, format, delimiter, "Australia/Sydney")
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
				Action: func(_ *cli.Context) error {
					departures, err := statusLine.DeparturesAction(routeName, stopName, directionName, departuresCount, timezone)
					PrintResult[statusLine.Departure](departures, format, delimiter, "Australia/Sydney")
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
					routeName := c.Args().First()
					directions, err := statusLine.DirectionsAction(routeName)
					PrintResult[statusLine.Direction](directions, format, delimiter, "Australia/Sydney")
					if err != nil {
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
		},
	}

	devId := os.Getenv("PTV_DEVID")
	key := os.Getenv("PTV_KEY")
	if key == "" || devId == "" {
		fmt.Println("PTV_KEY or PTV_DEVID not set in environment")
		os.Exit(1)
		return
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
