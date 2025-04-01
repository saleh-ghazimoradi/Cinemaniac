package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/dto"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/transaction"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/validator"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
)

type MovieService interface {
	CreateMovie(ctx context.Context, input *dto.Movie) (*domain.Movie, map[string]string, error)
	GetMovieById(ctx context.Context, id int64) (*domain.Movie, error)
}

type movieService struct {
	movieRepository repository.MovieRepository
	txService       transaction.TxService
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

	var createdMovie *domain.Movie
	err := m.txService.WithTx(ctx, func(tx *sql.Tx) error {
		repo := m.movieRepository.WithTx(ctx, tx)
		var err error
		createdMovie, err = repo.CreateMovie(ctx, movie)
		return err
	})

	if err != nil {
		slg.Logger.Error("error creating movie", "error", err)
		return nil, nil, errors.New("error creating movie")
	}

	return createdMovie, nil, nil
}

func (m *movieService) GetMovieById(ctx context.Context, id int64) (*domain.Movie, error) {
	movie, err := m.movieRepository.GetMovieById(ctx, id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func NewMovieService(movieRepository repository.MovieRepository, txService transaction.TxService) MovieService {
	return &movieService{
		movieRepository: movieRepository,
		txService:       txService,
	}
}
