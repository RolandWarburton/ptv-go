# PTV Status Line Test

Small tool to grab information from the PTV Timetable API.

Swagger documentation can be found [here](https://timetableapi.ptv.vic.gov.au/swagger/ui/index).

## Examples

```bash
# gets all the data as json
ptv-status-line routes

# format a route as a string
ptv-status-line routes Belgrave

# print information about the Ringwood stop on the Belgrave line
ptv-status-line stops --stop Ringwood 2

# print a stops attributes as a string
ptv-status-line stops --route "Belgrave" Ringwood

# print the directions for a route
# you can get the route ID from `directions routes "ROUTE NAME"`
ptv-status-line directions Belgrave

# print the next departures from a station in a direction
ptv-status-line departures \
--count 1 \
--direction "Flinders" \
--route "Belgrave" \
--stop "Ringwood" \
--timezone "Australia/Sydney" \
--format "ScheduledDepartureUTC"
```

## Formatting

By default each command will output a JSON response.
If the `--format` flag is encountered the output will be formatted
to the formats space separated names.

Formatting names can be found in `types.go`.

```bash
# example of using formatting
# print the just direction ID and the direction name
ptv-status-line directions --format "DirectionID DirectionName" --delimiter " -> " Belgrave

# OUTPUT:
# 3 -> Belgrave
# 1 -> City (Flinders Street)
```

## Limitations

The program does not support extracting nested values.
For example `--format "StopTicket.Zone"` will not work.

To get these values you can parse the values using a JSON processor such as [jq](https://github.com/jqlang/jq).

```bash
# get the zone for Ringwood Station on the Belgrave line
ptv-status-line stops --route Belgrave Ringwood | jq -r '.[0].stop_ticket.zone'
```
