package repository

import (
	"context"
	"database/sql"
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
	return nil, nil
}

func (m *movieRepository) GetMovies(ctx context.Context) ([]*domain.Movie, error) {
	return nil, nil
}

func (m *movieRepository) UpdateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	return nil, nil
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
