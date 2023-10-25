package db

import (
	"database/sql"
	"errors"
)

func CreateTables(database *sql.DB) (errs error) {
	_, err := database.Query(`
		CREATE TABLE IF NOT EXISTS users (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
			username VARCHAR(50) UNIQUE NOT NULL,
			passhash VARCHAR(97) NOT NULL,
			email VARCHAR(319) UNIQUE NOT NULL
		);
	`)
	errs = errors.Join(errs, err)

	_, err = database.Query(`
		CREATE TABLE IF NOT EXISTS cats (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
			owner_id uuid NOT NULL,
			username VARCHAR(50) UNIQUE NOT NULL,
			name VARCHAR(50) NOT NULL,
			follower_count INT DEFAULT 0,
			CONSTRAINT fk_owner
				FOREIGN KEY (owner_id)
					REFERENCES users(id)
		);
	`)
	errs = errors.Join(errs, err)

	return errs
}
