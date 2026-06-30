package main

import (
	"context"
	weatherv1 "example/weather-tracker/pb/weather/v1"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	weatherv1.UnimplementedWeatherServiceServer
}

func (s *server) GetWeather(ctx context.Context, in *weatherv1.GetWeatherRequest) (*weatherv1.GetWeatherResponse, error) {

	var temperature float64 = 35.0 // Hardcoded mock value - I'll replace it with OpenMeteo on the integration step.

	fmt.Printf("Request Temperature for Latitude (%v) and Longitude (%v) - Got %.2f C\n", in.Latitude, in.Longitude, temperature)

	return &weatherv1.GetWeatherResponse{
		Temperature: float64(temperature),
	}, nil
}

func main() {

	port := ":50051"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	weatherv1.RegisterWeatherServiceServer(grpcServer, &server{})

	log.Printf("gRPC Weather Server is running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
