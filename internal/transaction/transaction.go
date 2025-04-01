package transaction

import (
	"context"
	"database/sql"
	"errors"
)

type TxService interface {
	Begin(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	WithTx(ctx context.Context, fn func(*sql.Tx) error) error
}

type txService struct {
	db *sql.DB
}

func (t *txService) Begin(ctx context.Context) (*sql.Tx, error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to begin transaction: " + err.Error())
	}
	return tx, nil
}

func (t *txService) Commit(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}
	return nil
}

func (t *txService) Rollback(tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		return errors.New("failed to rollback transaction: " + err.Error())
	}
	return nil
}

func (t *txService) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}

	defer func() {
		if p := recover(); p != nil {
			_ = t.Rollback(tx)
			panic(p)
		}
	}()

	if err = fn(tx); err != nil {
		if rollbackErr := t.Rollback(tx); rollbackErr != nil {
			return errors.New("failed to rollback: " + rollbackErr.Error() + "; original error: " + err.Error())
		}
		return err
	}

	return t.Commit(tx)
}

func NewTXService(db *sql.DB) TxService {
	return &txService{
		db: db,
	}
}
