package api

import (
	"context"
	"errors"

	pb "github.com/bed787/golang-homework/gen/moviesv1"
	"github.com/bed787/golang-homework/internal/movie"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedMovieServiceServer
	service *movie.Service
}

func NewGRPCServer(s *movie.Service) *GRPCServer {
	return &GRPCServer{service: s}
}

func (s *GRPCServer) LoadMovies(ctx context.Context, req *pb.LoadMoviesRequest) (*pb.LoadMoviesResponse, error) {
	result, err := s.service.LoadMovies(ctx, req.GetFrom(), req.GetTo())
	if err != nil {
		return nil, grpcError(err)
	}
	response := &pb.LoadMoviesResponse{}
	for _, m := range result.Movies {
		response.Movies = append(response.Movies, toProto(m))
	}
	for _, f := range result.Failures {
		response.Failures = append(response.Failures, &pb.LoadFailure{Id: f.ID, Message: f.Message})
	}
	return response, nil
}

func (s *GRPCServer) GetLastMovie(ctx context.Context, _ *pb.GetLastMovieRequest) (*pb.Movie, error) {
	m, err := s.service.GetLastMovie(ctx)
	if err != nil {
		return nil, grpcError(err)
	}
	return toProto(m), nil
}

func toProto(m movie.Movie) *pb.Movie {
	return &pb.Movie{Id: m.ID, Title: m.Title, Year: m.Year, Director: m.Director}
}

func grpcError(err error) error {
	code := codes.Unavailable
	switch {
	case errors.Is(err, movie.ErrInvalidRange):
		code = codes.InvalidArgument
	case errors.Is(err, movie.ErrNotFound):
		code = codes.NotFound
	case errors.Is(err, context.DeadlineExceeded):
		code = codes.DeadlineExceeded
	case errors.Is(err, context.Canceled):
		code = codes.Canceled
	}
	return status.Error(code, err.Error())
}
