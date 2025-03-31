package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/dto"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/validator"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
)

type MovieService interface {
	CreateMovie(ctx context.Context, input *dto.Movie) (*domain.Movie, map[string]string, error)
}

type movieService struct {
	movieRepository repository.MovieRepository
}

func (m *movieService) CreateMovie(ctx context.Context, input *dto.Movie) (*domain.Movie, map[string]string, error) {
	v := validator.New()

	movie := &domain.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	domain.ValidateMovie(v, movie)

	if !v.Valid() {
		slg.Logger.Error("validation failed", "errors", v.Errors)
		return nil, v.Errors, errors.New("validation failed")
	}

	movie, err := m.movieRepository.CreateMovie(ctx, movie)
	if err != nil {
		slg.Logger.Error("error creating movie", "error", err)
		return nil, nil, errors.New("error creating movie")
	}

	return movie, nil, nil
}

func NewMovieService(movieRepository repository.MovieRepository) MovieService {
	return &movieService{
		movieRepository: movieRepository,
	}
}
