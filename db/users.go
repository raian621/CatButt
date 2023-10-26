package db

import (
	"database/sql"
	"errors"

	"catbook.com/auth"
	"catbook.com/util"
)

type User struct {
	UserId   string
	Username string
	Passhash string
	Email    string
}

func CreateUser(userInfo *auth.UserRegistrationInfo, database *sql.DB) error {
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

func DeleteUser(userId string, database *sql.DB) error {
	_, err := database.Exec(`DELETE FROM users WHERE id=$1;`, userId)
	return err
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
		return false, nil
	}

	return util.MatchHash(creds.Password, passhash)
}

func GetUserByUsername(username string, database *sql.DB) (*User, error) {
	var user User
	row := database.QueryRow(`SELECT id, username, email FROM users WHERE username=$1;`, username)
	err := row.Scan(&user.UserId, &user.Username, &user.Email)

	return &user, err
}
