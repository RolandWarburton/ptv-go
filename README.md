# PTV Status Line Test

Exploring the PTV API to get the next scheduled departure for trains.

## Examples

Print the Belgrave route.

```bash
# gets all the data as json
ptv-status-line routes Belgrave

# format the output as a string
ptv-status-line --format "RouteID RouteGtfsID RouteName" --delimiter " - " routes Belgrave
```
