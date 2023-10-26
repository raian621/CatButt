package db

import (
	"database/sql"
	"errors"
)

func CreateTables(database *sql.DB) (errs error) {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id       UUID         PRIMARY KEY DEFAULT gen_random_uuid (),
			username VARCHAR(50)  UNIQUE NOT NULL,
			passhash VARCHAR(97)  NOT NULL,
			email    VARCHAR(319) UNIQUE NOT NULL
		);
	`)
	errs = errors.Join(errs, err)

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS cats (
			id             UUID        PRIMARY KEY DEFAULT gen_random_uuid (),
			owner_id       UUID        NOT NULL,
			username       VARCHAR(50) UNIQUE NOT NULL,
			name           VARCHAR(50) NOT NULL,
			follower_count INT         DEFAULT 0,
			CONSTRAINT fk_owner
				FOREIGN KEY (owner_id)
					REFERENCES users(id)
					ON DELETE CASCADE
		);
	`)
	errs = errors.Join(errs, err)

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id         CHAR(32)    PRIMARY KEY NOT NULL,
			user_id    UUID        NOT NULL,
			expires    TIMESTAMP   NOT NULL,
			os         VARCHAR(50) NOT NULL,
			browser    VARCHAR(50) NOT NULL,
			ip_address VARCHAR(39) NOT NULL,
			CONSTRAINT fk_session_user
				FOREIGN KEY (user_id)
				REFERENCES users(id)
				ON DELETE CASCADE
		);
	`)
	errs = errors.Join(errs, err)

	return errs
}
