package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseParams struct {
	Provider string
	Hostname string
	Username string
	Password string
	Database string
	Port     string
	SSLMode  string
}

func databaseURL(dbParams *DatabaseParams) (url string) {
	url = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbParams.Username,
		dbParams.Password,
		dbParams.Hostname,
		dbParams.Port,
		dbParams.Database,
		dbParams.SSLMode,
	)

	return url
}

func ConnectToDB(dbParams *DatabaseParams) (database *sql.DB, err error) {
	url := databaseURL(dbParams)
	database, err = sql.Open(dbParams.Provider, url)
	if err != nil {
		fmt.Println(database)
	}

	return database, err
}
