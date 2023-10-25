package db

import (
	"database/sql"
	"errors"
	"fmt"

	"catbook.com/auth"
	"catbook.com/util"
)

func CreateUser(userInfo *auth.UserRegistrationInfo, database *sql.DB) (err error) {
	if UserOrEmailExists(userInfo.Username, userInfo.Email, database) {
		return errors.New("user already exists")
	}

	passhash, err := util.DefaultGenerateHash(userInfo.Password)

	if err != nil {
		return err
	}

	_, err = database.Exec(
		`INSERT INTO users (username, passhash, email) VALUES ($1, $2, $3);`,
		userInfo.Username, passhash, userInfo.Email,
	)

	return err
}

func DeleteUser(userCreds *auth.UserCredentials, database *sql.DB) error {
	return nil
}

func UserOrEmailExists(username, email string, database *sql.DB) bool {
	row := database.QueryRow(`
		SELECT username FROM users WHERE username=$1 OR email=$2;
	`, username, email)

	err := row.Scan()

	return err != sql.ErrNoRows
}

func ValidUserCredentials(creds *auth.UserCredentials, database *sql.DB) (bool, error) {
	row := database.QueryRow(`
		SELECT passhash FROM users WHERE username=$1;
	`, creds.Username)

	var passhash string
	if err := row.Scan(&passhash); err == sql.ErrNoRows {
		fmt.Println(passhash)

		return false, nil
	}

	fmt.Println(passhash)
	return util.MatchHash(creds.Password, passhash)
}
