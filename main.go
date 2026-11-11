package main

import (
	"log"
	"net"
	"net/http"

	pb "github.com/bed787/golang-homework/gen/moviesv1"
	"github.com/bed787/golang-homework/internal/api"
	"github.com/bed787/golang-homework/internal/config"
	"github.com/bed787/golang-homework/internal/movie"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	service := movie.NewService(cfg.MovieBaseURL)

	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMovieServiceServer(grpcServer, api.NewGRPCServer(service))
	go func() {
		log.Printf("gRPC: %s", cfg.GRPCAddr)
		log.Fatal(grpcServer.Serve(listener))
	}()

	log.Printf("HTTP: %s; Swagger: /swagger/index.html", cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, api.NewHTTPHandler(service)))
}
