package repository

import "database/sql"

func exec(db *sql.DB, tx *sql.Tx) *sql.Tx {
	if tx != nil {
		return tx
	}
	tx, _ = db.Begin()
	
	return tx
}
