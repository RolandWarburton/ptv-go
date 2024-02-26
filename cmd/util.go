package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

func WriteToJSONFile(data any) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")

	// write the routes to a file
	file, _ := os.Create("routes.json")
	defer file.Close()
	file.Write(jsonData)
}

// format output as a string
func PrintFormatted[Type any](data []Type, format string, delimiter string) {
	// Example --format "RouteID RouteName"
	formatArgs := strings.Split(format, " ")
	result := ""
	for i, item := range data {
		for j, arg := range formatArgs {
			// dynamically access the fields of the Route
			val := reflect.ValueOf(item)
			field := val.FieldByName(arg)
			if field.IsValid() && j < len(formatArgs)-1 {
				result += fmt.Sprintf("%v%s", field.Interface(), delimiter)
			} else {
				result += fmt.Sprintf("%v", field.Interface())
			}
		}
		if i < len(data)-1 {
			result += "\n"
		}
	}
	fmt.Println(result)
}

func ConvertDatesToTimezone[Type interface{}](data []Type, timezone string) ([]Type, error) {
	var location *time.Location
	var err error
	if timezone == "UTC" {
		location = time.UTC
	} else {
		location, err = time.LoadLocation(timezone)
		if err != nil {
			fmt.Println("Error loading location:", err)
			return nil, err
		}
	}

	slice := reflect.ValueOf(data)

	// check a slice is being given
	if slice.Kind() != reflect.Slice {
		return nil, errors.New("a slice was not given")
	}

	// create a copy of the data
	// we can overwrite these elements later as we iterate over the slice
	newData := reflect.MakeSlice(reflect.TypeOf(data), slice.Len(), slice.Len())

	//for each element in the reflected slice
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)

		// skip if this is not an interface that we can print
		if !elem.CanInterface() {
			newData.Index(i).Set(elem)
			continue
		}

		// get the field in this interface
		field := elem.FieldByName("ScheduledDepartureUTC")

		// skip if this is not the ScheduledDepartureUTC field
		if !field.IsValid() {
			newData.Index(i).Set(elem)
			continue
		}

		// // parse the date string into a date
		layout := "02-01-2006 03:04 PM"
		departureDate, err := time.Parse(layout, field.Interface().(string))
		if err != nil {
			fmt.Println(err)
			field.SetString("date conversion failed")
			newData.Index(i).Set(elem)
			continue
		}

		// convert the timezone
		departureDateTZ := departureDate.In(location)
		field.SetString(departureDateTZ.Format("02-01-2006 03:04 PM"))

		// overwrite the newData interface reflection with this new element
		newData.Index(i).Set(elem)
	}

	// reconstruct the newData reflection as the initial data type
	var result []Type
	for i := 0; i < newData.Len(); i++ {
		result = append(result, newData.Index(i).Interface().(Type))
	}

	return result, nil
}

func PrettyPrint[Type interface{}](data []Type, timezone string) {
	// convert the date to the specified timezone
	data, _ = ConvertDatesToTimezone[Type](data, timezone)

	// pretty print
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}
