package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Current   struct {
		Time        string  `json:"time"`
		Temperature float64 `json:"temperature_2m"`
		Humidity    float64 `json:"relative_humidity_2m"`
	} `json:"current"`
}

func GetWeatherData(lat, long float64, result *Result) error {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m", lat, long)

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unable to connect to Open-meteo: Status (%d) %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var weather OpenMeteoResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return err
	}

	// Populate result with data from OpenMeteo
	result.Temperature = weather.Current.Temperature
	result.Elevation = weather.Elevation
	result.Humidity = weather.Current.Humidity

	return nil
}
