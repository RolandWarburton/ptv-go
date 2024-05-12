package ptvgo

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
	Status Status  `json:"status"`
}

type Ticket struct {
	TicketType       string `json:"ticket_type"`
	Zone             string `json:"zone"`
	IsFreeFareZone   bool   `json:"is_free_fare_zone"`
	TicketMachine    bool   `json:"ticket_machine"`
	TicketChecks     bool   `json:"ticket_checks"`
	VLineReservation bool   `json:"vline_reservation"`
	TicketZones      []int  `json:"ticket_zones"`
}

type Stop struct {
	DisruptionIds []string `json:"disruption_ids"`
	StopSuburb    string   `json:"stop_suburb"`
	RouteType     int      `json:"route_type"`
	StopLatitude  float64  `json:"stop_latitude"`
	StopLongitude float64  `json:"stop_longitude"`
	StopSequence  int      `json:"stop_sequence"`
	StopTicket    Ticket   `json:"stop_ticket"`
	StopID        int      `json:"stop_id"`
	StopName      string   `json:"stop_name"`
	StopLandmark  string   `json:"stop_landmark"`
}

type StopResponse struct {
	Stops  []Stop `json:"stops"`
	Status Status `json:"status"`
}

type Direction struct {
	RouteDirectionDescription string `json:"route_direction_description"`
	DirectionID               int    `json:"direction_id"`
	DirectionName             string `json:"direction_name"`
	RouteID                   int    `json:"route_id"`
	RouteType                 int    `json:"route_type"`
}

type DirectionsResponse struct {
	Directions []Direction `json:"directions"`
	Status     Status      `json:"status"`
}
