package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	weatherv1 "go-buf/gen/go/myapp/weather/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":8080"

type weatherService struct {
	weatherv1.UnimplementedWeatherServiceServer
}

func (w *weatherService) GetWeather(ctx context.Context, req *weatherv1.GetWeatherRequest) (*weatherv1.GetWeatherResponse, error) {
	return &weatherv1.GetWeatherResponse{
		Temperature: 10,
		Condition:   weatherv1.Condition_CONDITION_SUNNY,
	}, nil
}

func main() {
	listner, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	weatherv1.RegisterWeatherServiceServer(s, &weatherService{})

	reflection.Register(s)

	go func() {
		if err := s.Serve(listner); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
}
