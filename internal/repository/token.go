package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
)

type TokenRepository interface {
	Insert(ctx context.Context, token *domain.Token) error
	DeleteAllForUser(ctx context.Context, scope string, userId int64) error
	WithTx(ctx context.Context, tx *sql.Tx) TokenRepository
}

type tokenRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (t *tokenRepository) Insert(ctx context.Context, token *domain.Token) error {
	query := `
        INSERT INTO tokens (hash, user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.CTX.Timeout)
	defer cancel()

	_, err := exec(t.dbWrite, t.tx).ExecContext(ctx, query, args...)

	return err
}

func (t *tokenRepository) DeleteAllForUser(ctx context.Context, scope string, userId int64) error {
	query := `
        DELETE FROM tokens 
        WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.CTX.Timeout)
	defer cancel()

	_, err := exec(t.dbRead, t.tx).ExecContext(ctx, query, scope, userId)
	return err
}

func (t *tokenRepository) WithTx(ctx context.Context, tx *sql.Tx) TokenRepository {
	return &tokenRepository{
		dbWrite: t.dbWrite,
		dbRead:  t.dbRead,
		tx:      tx,
	}
}

func NewTokenRepository(dbWrite, dbRead *sql.DB) TokenRepository {
	return &tokenRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
