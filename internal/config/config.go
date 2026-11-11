package config

import "os"

type Config struct {
	MovieBaseURL string
	HTTPAddr     string
	GRPCAddr     string
}

func Load() Config {
	return Config{
		MovieBaseURL: env("MOVIE_BASE_URL", "https://homeworksite.site"),
		HTTPAddr:     env("HTTP_ADDR", ":8080"),
		GRPCAddr:     env("GRPC_ADDR", ":50051"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
