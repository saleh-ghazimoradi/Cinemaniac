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
	UpdateMovie(ctx context.Context, id int64, input *dto.UpdateMovie) (*domain.Movie, map[string]string, error)
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

func (m *movieService) fetchMovie(ctx context.Context, id int64) (*domain.Movie, error) {
	return m.movieRepository.GetMovieById(ctx, id)
}

func (m *movieService) UpdateMovie(ctx context.Context, id int64, input *dto.UpdateMovie) (*domain.Movie, map[string]string, error) {
	var updatedMovie *domain.Movie
	var validationErrors map[string]string

	err := m.txService.WithTx(ctx, func(tx *sql.Tx) error {
		// Get repository with transaction
		txRepo := m.movieRepository.WithTx(ctx, tx)

		// 1. Fetch existing movie
		movie, err := txRepo.GetMovieById(ctx, id)
		if err != nil {
			return err
		}

		// 2. Apply updates
		if input.Title != nil {
			movie.Title = *input.Title
		}
		if input.Year != nil {
			movie.Year = *input.Year
		}
		if input.Runtime != nil {
			movie.Runtime = *input.Runtime
		}
		if input.Genres != nil {
			movie.Genres = input.Genres
		}

		// 3. Validate
		v := validator.New()
		domain.ValidateMovie(v, movie)
		if !v.Valid() {
			validationErrors = v.Errors
			return errors.New("validation failed")
		}

		// 4. Update
		updatedMovie, err = txRepo.UpdateMovie(ctx, movie)
		return err
	})

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, nil, err
		case validationErrors != nil:
			return nil, validationErrors, nil
		default:
			return nil, nil, err
		}
	}

	return updatedMovie, nil, nil
}

func NewMovieService(movieRepository repository.MovieRepository, txService transaction.TxService) MovieService {
	return &movieService{
		movieRepository: movieRepository,
		txService:       txService,
	}
}
