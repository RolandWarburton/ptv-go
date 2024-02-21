# NOTES

## Links

https://timetableapi.ptv.vic.gov.au/swagger/ui/index#/Stops

Get all the routes

```js
const requestString = '/v3/routes';
```

Get route types

```js
const requestString = '/v3/route_types';
```

Train route type is 0

SC station routes

```js
const requestString = '/v3/runs/routes/1181';
```

Get ever train route

```js
const requestString = '/v3/routes?route_types=0';
```

Belgrave route id 2
Belgrave route type 0

Get the Belgrave line stops

```js
const requestString = '/v3/stops/route/2/route_type/0';
```

Bayswater stop id 1016
SC stop id 1181

Get departures from Bayswater

```js
const requestString = '/v3/departures/route_type/0/stop/1016';
```

Get direction IDs from Bayswater

```js
const requestString = '/v3/directions/route/2';
```

BW -> SC 1
