package main

import (
	"context"
	"database/sql"
	"errors"
	weatherv1 "example/weather-tracker/pb/weather/v1"
	"log"
	"net"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mmcloughlin/geohash"
	"google.golang.org/grpc"
)

type result struct {
	Geohash     string
	Temperature float64
	Humidity    float64
	Elevation   float64
}

type server struct {
	weatherv1.UnimplementedWeatherServiceServer
	db *sql.DB
}

func openDB(connString string) (*sql.DB, error) {

	/*

		-- DATABASE TABLE SCHEMA --

		TABLE: weather
		id          - INTEGER AUTO INCREMENT
		geohash     - Varchar(12)
		temperature - Numeric(4, 2) - 2 point precision
		humidity    - Numeric(4, 2) - 2 point precision
		elevation   - Numeric(4, 2) - 2 point precision

	*/

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil

}

func (s *server) GetTemperature(ctx context.Context, lat, long float64) (result, bool) {

	var res = result{}

	hash := geohash.EncodeWithPrecision(lat, long, 6)

	res.Geohash = hash

	query := `
	SELECT temperature, humidity, elevation
	FROM weather
	WHERE geohash = $1`

	err := s.db.QueryRowContext(ctx, query, hash).Scan(&res.Temperature, &res.Humidity, &res.Elevation)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, false
		}
		return res, false
	}

	return res, true

}

func (s *server) SaveTemperature(ctx context.Context, result result) error {

	query := `
	INSERT INTO 
	weather (geohash, temperature, humidity, elevation)
	VALUES ($1, $2, $3, $4)`

	_, err := s.db.ExecContext(ctx, query, result.Geohash, result.Temperature, result.Humidity, result.Elevation)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) GetWeather(ctx context.Context, in *weatherv1.GetWeatherRequest) (*weatherv1.GetWeatherResponse, error) {

	var result result
	// Check the database with lat and long to get the value.
	result, ok := s.GetTemperature(ctx, in.Latitude, in.Longitude)
	if ok {
		log.Println("Cache Hit - Retrieving data from database...")
		return &weatherv1.GetWeatherResponse{
			Temperature: float64(result.Temperature),
			Humidity:    float64(result.Humidity),
			Elevation:   float64(result.Elevation),
		}, nil
	}

	// If not available in database - get the value from OpenMeteo
	log.Println("Cache Miss - Retrieving data from OpenMeteo API...")
	// Hardcoded mock value - I'll replace it with OpenMeteo on the integration step.
	result.Temperature = 35.0
	result.Humidity = 30.0
	result.Elevation = 5.0

	// Save new value to database
	err := s.SaveTemperature(ctx, result)
	if err != nil {
		log.Printf("Failed to save result to database %v\n", err)
	}

	log.Printf("Request Temperature for Latitude (%v) and Longitude (%v) - Got %.2f C\n", in.Latitude, in.Longitude, result.Temperature)

	return &weatherv1.GetWeatherResponse{
		Temperature: float64(result.Temperature),
		Humidity:    float64(result.Humidity),
		Elevation:   float64(result.Elevation),
	}, nil
}

func main() {

	connStr := "postgres://web:pass@localhost:5432/weatherapp?sslmode=disable"

	db, err := openDB(connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		os.Exit(1)
	}

	log.Println("Connected to PostgreSQL database: weatherapp")

	port := ":50051"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	weatherv1.RegisterWeatherServiceServer(grpcServer, &server{db: db})

	log.Printf("gRPC Weather Server is running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
