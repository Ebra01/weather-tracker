package main

import (
	"fmt"

	weatherv1 "example/weather-tracker/pb/weather/v1"
)

func main() {
	req := &weatherv1.GetWeatherRequest{
		Latitude:  37.7749,
		Longitude: -122.4194,
	}

	fmt.Printf("Client executing request for  Latitude & Longitude: %f, %f\n", req.GetLatitude(), req.GetLongitude())
}
