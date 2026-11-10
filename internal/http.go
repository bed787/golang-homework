package internal

import (
	"net/http"
	"strconv"
	"log"
)

func loadMovies(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	from, err := strconv.Atoi(fromStr)
	if err != nil || from < 1 {
		http.Error(w, "Invalid from", http.StatusBadRequest)
		return
	}

	toStr := r.URL.Query().Get("to")

	to, err := strconv.Atoi(toStr)
	if err != nil || from > to {
		http.Error(w, "Invalid to", http.StatusBadRequest)
		return
	}

	client := NewMovieClient()
	if err != nil {
		log.Printf("GetMovie Failed: %v", err)
		http.Error(w, "Failed to fetch", http.StatusBadGateway)
	}
}

func RegisterRouters() {
	http.HandleFunc("GET /load-movies", loadMovies)
}
