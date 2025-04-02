package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/dto"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/transaction"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/validator"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
)

type MovieService interface {
	CreateMovie(ctx context.Context, input *dto.Movie) (*domain.Movie, error)
	GetMovieById(ctx context.Context, id int64) (*domain.Movie, error)
	GetMovies(ctx context.Context) ([]*domain.Movie, error)
	UpdateMovie(ctx context.Context, id int64, input *dto.UpdateMovie) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

type movieService struct {
	movieRepository repository.MovieRepository
	txService       transaction.TxService
}

func (m *movieService) CreateMovie(ctx context.Context, input *dto.Movie) (*domain.Movie, error) {
	v := validator.New()

	movie := &domain.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	domain.ValidateMovie(v, movie)

	if err := v.GetValidationError(); err != nil {
		slg.Logger.Error("validation failed", "errors", v.Errors)
		return nil, err
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
		return nil, err
	}

	return createdMovie, nil
}

func (m *movieService) GetMovieById(ctx context.Context, id int64) (*domain.Movie, error) {
	movie, err := m.movieRepository.GetMovieById(ctx, id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *movieService) GetMovies(ctx context.Context) ([]*domain.Movie, error) {
	return m.movieRepository.GetMovies(ctx)
}

func (m *movieService) fetchMovie(ctx context.Context, id int64) (*domain.Movie, error) {
	return m.movieRepository.GetMovieById(ctx, id)
}

func (m *movieService) UpdateMovie(ctx context.Context, id int64, input *dto.UpdateMovie) (*domain.Movie, error) {
	var updatedMovie *domain.Movie

	err := m.txService.WithTx(ctx, func(tx *sql.Tx) error {
		txRepo := m.movieRepository.WithTx(ctx, tx)

		movie, err := txRepo.GetMovieById(ctx, id)
		if err != nil {
			return err
		}

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

		v := validator.New()
		domain.ValidateMovie(v, movie)
		if err = v.GetValidationError(); err != nil {
			return err
		}

		updatedMovie, err = txRepo.UpdateMovie(ctx, movie)
		return err
	})

	if err != nil {
		return nil, err
	}

	return updatedMovie, nil
}

func (m *movieService) DeleteMovie(ctx context.Context, id int64) error {
	return m.txService.WithTx(ctx, func(tx *sql.Tx) error {
		txRepo := m.movieRepository.WithTx(ctx, tx)
		return txRepo.DeleteMovie(ctx, id)
	})
}

func NewMovieService(movieRepository repository.MovieRepository, txService transaction.TxService) MovieService {
	return &movieService{
		movieRepository: movieRepository,
		txService:       txService,
	}
}
