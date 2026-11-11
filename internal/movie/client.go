package movie

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	baseURL string
	client  *http.Client
}

func NewService(baseURL string) *Service {
	return &Service{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{},
	}
}

func (s *Service) fetchMovie(ctx context.Context, id int64) (Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/%d/info.0.json", s.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Movie{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return Movie{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Movie{}, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return Movie{}, fmt.Errorf("movie service returned HTTP %d", resp.StatusCode)
	}
	var m Movie
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return Movie{}, fmt.Errorf("invalid movie JSON: %w", err)
	}
	return m, nil
}
