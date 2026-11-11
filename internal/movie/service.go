package movie

import (
	"context"
	"errors"
	"math"
)

func (s *Service) LoadMovies(ctx context.Context, from, to int64) (LoadResult, error) {
	if from < 1 || to < from {
		return LoadResult{}, ErrInvalidRange
	}
	result := LoadResult{Movies: []Movie{}}
	var loadErr error
	for id := from; ; id++ {
		m, err := s.fetchMovie(ctx, id)
		if ctx.Err() != nil {
			return result, ctx.Err()
		}
		if err != nil {
			result.Failures = append(result.Failures, LoadFailure{ID: id, Message: err.Error()})
			if loadErr == nil || errors.Is(loadErr, ErrNotFound) {
				loadErr = err
			}
		} else {
			result.Movies = append(result.Movies, m)
		}
		if id == to {
			break
		}
	}
	if len(result.Movies) == 0 {
		return result, loadErr
	}
	return result, nil
}

func (s *Service) GetLastMovie(ctx context.Context) (Movie, error) {
	last, err := s.fetchMovie(ctx, 1)
	if err != nil {
		return Movie{}, err
	}
	left, right := int64(1), int64(2)
	for {
		m, err := s.fetchMovie(ctx, right)
		if errors.Is(err, ErrNotFound) {
			break
		}
		if err != nil {
			return Movie{}, err
		}
		last = m
		left = right
		if right == math.MaxInt64 {
			return last, nil
		}
		if right > math.MaxInt64/2 {
			right = math.MaxInt64
		} else {
			right *= 2
		}
	}
	for right-left > 1 {
		mid := left + (right-left)/2
		m, err := s.fetchMovie(ctx, mid)
		if errors.Is(err, ErrNotFound) {
			right = mid
		} else if err != nil {
			return Movie{}, err
		} else {
			last = m
			left = mid
		}
	}
	return last, nil
}
