package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bed787/golang-homework/docs"
	"github.com/bed787/golang-homework/internal/movie"
	swaggerFiles "github.com/swaggo/files/v2"
)

type HTTPHandler struct {
	service *movie.Service
}

func NewHTTPHandler(s *movie.Service) http.Handler {
	h := &HTTPHandler{service: s}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /load-movies", h.loadMovies)
	mux.HandleFunc("GET /last-movie", h.lastMovie)
	mux.HandleFunc("GET /swagger/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(docs.SwaggerHTML)
	})
	mux.HandleFunc("GET /swagger/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(docs.OpenAPI)
	})
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerFiles.FS))))
	return mux
}

func (h *HTTPHandler) loadMovies(w http.ResponseWriter, r *http.Request) {
	from, errFrom := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
	to, errTo := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
	if errFrom != nil || errTo != nil {
		writeError(w, movie.ErrInvalidRange)
		return
	}
	result, err := h.service.LoadMovies(r.Context(), from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	statusCode := http.StatusOK
	if len(result.Failures) > 0 {
		failures, _ := json.Marshal(result.Failures)
		w.Header().Set("X-Movie-Errors", string(failures))
		statusCode = http.StatusMultiStatus
	}
	writeJSON(w, statusCode, result.Movies)
}

func (h *HTTPHandler) lastMovie(w http.ResponseWriter, r *http.Request) {
	m, err := h.service.GetLastMovie(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func writeError(w http.ResponseWriter, err error) {
	statusCode := http.StatusBadGateway
	switch {
	case errors.Is(err, movie.ErrInvalidRange):
		statusCode = http.StatusBadRequest
	case errors.Is(err, movie.ErrNotFound):
		statusCode = http.StatusNotFound
	case errors.Is(err, context.DeadlineExceeded):
		statusCode = http.StatusGatewayTimeout
	case errors.Is(err, context.Canceled):
		statusCode = http.StatusRequestTimeout
	}
	writeJSON(w, statusCode, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(value)
}
