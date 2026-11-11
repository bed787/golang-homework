package movie

import "errors"

type Movie struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Year     int32  `json:"year"`
	Director string `json:"director"`
}

type LoadFailure struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

type LoadResult struct {
	Movies   []Movie
	Failures []LoadFailure
}

var ErrInvalidRange = errors.New("invalid range: expected 1 <= from <= to")
var ErrNotFound = errors.New("movie not found")
