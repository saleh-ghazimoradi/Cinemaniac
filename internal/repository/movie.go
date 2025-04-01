package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	GetMovieById(ctx context.Context, id int64) (*domain.Movie, error)
	GetMovies(ctx context.Context) ([]*domain.Movie, error)
	UpdateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
	WithTx(ctx context.Context, tx *sql.Tx) MovieRepository
}

type movieRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (m *movieRepository) CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	query := `
        INSERT INTO movies (title, year, runtime, genres) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	err := exec(m.dbWrite, m.tx).QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (m *movieRepository) GetMovieById(ctx context.Context, id int64) (*domain.Movie, error) {
	query := `
        SELECT id, created_at, title, year, runtime, genres, version
        FROM movies
        WHERE id = $1`

	movie := &domain.Movie{}

	if err := exec(m.dbRead, m.tx).QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return movie, nil
}

func (m *movieRepository) GetMovies(ctx context.Context) ([]*domain.Movie, error) {
	return nil, nil
}

func (m *movieRepository) UpdateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	query := `
        UPDATE movies 
        SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
        WHERE id = $5
        RETURNING version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID}

	if err := exec(m.dbWrite, m.tx).QueryRowContext(ctx, query, args...).Scan(&movie.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return movie, nil
}

func (m *movieRepository) DeleteMovie(ctx context.Context, id int64) error {
	return nil
}

func (m *movieRepository) WithTx(ctx context.Context, tx *sql.Tx) MovieRepository {
	return &movieRepository{
		dbWrite: m.dbWrite,
		dbRead:  m.dbRead,
		tx:      tx,
	}
}

func NewMovieRepository(dbWrite, dbRead *sql.DB) MovieRepository {
	return &movieRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
