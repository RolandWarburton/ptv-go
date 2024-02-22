package app

type Status struct {
	version string `json:"version"`
	health  int    `json:"health"`
}

type Departure struct {
	StopID                int    `json:"stop_id"`
	RouteID               int    `json:"route_id"`
	RunID                 int    `json:"run_id"`
	RunRef                string `json:"run_ref"`
	DirectionID           int    `json:"direction_id"`
	DisruptionIDs         []int  `json:"disruption_ids"`
	ScheduledDepartureUTC string `json:"scheduled_departure_utc"`
	EstimatedDepartureUTC string `json:"estimated_departure_utc"`
	AtPlatform            bool   `json:"at_platform"`
	PlatformNumber        string `json:"platform_number"`
	Flags                 string `json:"flags"`
	DepartureSequence     int    `json:"departure_sequence"`
}

type DepartureResponse struct {
	Departures []Departure `json:"departures"`
}

type RouteServiceStatus struct {
	description string `json:"description"`
	timestamp   string `json:"timestamp"`
}

type Route struct {
	routeServiceStatus RouteServiceStatus `json:"route_service_status"`
	routeType          int                `json:"route_type"`
	RouteID            int                `json:"route_id"`
	RouteName          string             `json:"route_name"`
	RouteNumber        string             `json:"route_number"`
	RouteGtfsID        string             `json:"route_gtfs_id"`
}

type RouteResponse struct {
	Routes []Route `json:"routes"`
  Status Status `json:"status"`
}
