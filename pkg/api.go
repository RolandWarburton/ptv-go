package app

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetUrl(request string) (string, error) {
	devId := os.Getenv("PTV_DEVID")
	key := os.Getenv("PTV_KEY")
	if key == "" || devId == "" {
		return "", errors.New("PTV_KEY or PTV_DEVID not set in environment")
	}
	baseURL := "http://timetableapi.ptv.vic.gov.au"

	if strings.Contains(request, "?") {
		request = request + "&"
	} else {
		request = request + "?"
	}
	raw := request + fmt.Sprintf("devid=%s", devId)
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(raw))
	signature := hex.EncodeToString(h.Sum(nil))
	url := fmt.Sprintf("%s%s&signature=%s", baseURL, raw, signature)
	// print the URL for debugging
	// fmt.Println(url)
	return url, nil
}

func PrintFormattedDate(dateObj time.Time) {
	formattedDate := dateObj.Format("15:04 PM")
	fmt.Println(formattedDate)
}

func GetRoutes(routeName string) ([]Route, error) {
	requestString := "/v3/routes?route_types=0"
	url, err := GetUrl(requestString)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error! Status: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response RouteResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error:", err)
	}

	routes := response.Routes

	// if we are not searching for a route name return all the routes
	if routeName == "" {
		return routes, nil
	}

	matchingRoutes := []Route{}

	for i := 0; i < len(routes); i++ {
		route := routes[i]
		if strings.Contains(route.RouteName, routeName) {
			matchingRoutes = append(matchingRoutes, route)
		}
	}

	return matchingRoutes, nil
}

func GetDepartures(stopID int, routeID int, queryParams string) ([]Departure, error) {
	requestString := fmt.Sprintf("/v3/departures/route_type/0/stop/%d/route/%d%s", stopID, routeID, queryParams)
	url, err := GetUrl(requestString)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error! Status: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response DepartureResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error:", err)
	}

	departures := response.Departures
	return departures, nil
}

func GetNextDepartureTowards(departures []Departure, directionID int, count int) ([]Departure, error) {
	now := time.Now()

	validDepartures := make([]Departure, 0)
	i := 0

	for _, departure := range departures {
		departureDateStr := departure.ScheduledDepartureUTC
		// departureDateStr, ok := departure["scheduled_departure_utc"].(string)
		// if !ok {
		// 	return nil, fmt.Errorf("failed to parse departure date")
		// }

		departureDate, err := time.Parse(time.RFC3339, departureDateStr)
		if err != nil {
			return nil, err
		}

		// if the train is not going in the direction we want skip it
		if int(departure.DirectionID) != directionID {
			continue
		}

		// if the train already departed skip it
		if departureDate.Before(now) {
			continue
		}

		// add the train to the list of departures to return
		validDepartures = append(validDepartures, departure)

		// if we have returned the number of departures required then return them
		i++
		if i == count {
			return validDepartures, nil
		}
	}

	return validDepartures, nil
}
