package db_test

import (
	"testing"

	"catbook.com/db"
)

func TestConnectToDB(t *testing.T) {
	_, err := db.ConnectToDB(&db.DatabaseParams{
		Provider: "postgres",
		Hostname: "localhost",
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
		Port:     "5432",
		SSLMode:  "disabled",
	})

	if err != nil {
		t.Error("ConnectToDB should have failed since the database is not set up.")
	}
}
