package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
)

type PermissionRepository interface {
	GetAllForUser(userID int64) (*domain.Permissions, error)
	AddForUser(userID int64, codes ...string) error
	WithTx(ctx context.Context, tx *sql.Tx) PermissionRepository
}

type permissionRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (p *permissionRepository) GetAllForUser(userID int64) (*domain.Permissions, error) {
	query := `
        SELECT permissions.code
        FROM permissions
        INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
        INNER JOIN users ON users_permissions.user_id = users.id
        WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig.CTX.Timeout)
	defer cancel()

	rows, err := exec(p.dbRead, p.tx).QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions domain.Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &permissions, nil

}

func (p *permissionRepository) AddForUser(userID int64, codes ...string) error {
	query := `
        INSERT INTO users_permissions
        SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig.CTX.Timeout)
	defer cancel()

	_, err := exec(p.dbWrite, p.tx).ExecContext(ctx, query, userID, pq.Array(codes))
	return err
}

func (p *permissionRepository) WithTx(ctx context.Context, tx *sql.Tx) PermissionRepository {
	return &permissionRepository{
		dbWrite: p.dbWrite,
		dbRead:  p.dbRead,
		tx:      tx,
	}
}

func NewPermissionRepository(dbWrite *sql.DB, dbRead *sql.DB) PermissionRepository {
	return &permissionRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
