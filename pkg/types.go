package app

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
