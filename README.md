# PTV Status Line Test

Exploring the PTV API to get the next scheduled departure for trains.

## Examples

Formatting names can be found in `types.go`.

```bash
# gets all the data as json
ptv-status-line routes Belgrave

# format a route as a string
ptv-status-line routes --format "RouteID RouteGtfsID RouteName" --delimiter " - "  Belgrave

# print information about the Ringwood stop on the Bayswater line
ptv-status-line stops --stop Ringwood 2

# print a stops attribute as a string
ptv-status-line stops --format "StopName" --delimiter " - " --stop Ringwood 2
```

## Examples (continued)

The program does not support extracting nested values.
For example `--format "StopTicket.Zone"` will not work.

To get these values you can parse the values using a JSON processor such as [jq](https://github.com/jqlang/jq).

```bash
# get the zone for Ringwood Station on the Belgrave line
ptv-status-line stops --stop Ringwood 2 |  jq -r '.[0].stop_ticket.zone'
```
