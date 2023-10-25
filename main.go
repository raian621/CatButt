package main

import (
	"fmt"
	"os"

	"catbook.com/db"
	"catbook.com/server"
	"catbook.com/util"
)

func main() {
	util.LoadEnvVars(".env")

	dbParams := &db.DatabaseParams{
		Provider: os.Getenv("DB_PROV"),
		Hostname: os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSL"),
	}

	util.DeleteEnvVars()

	database, err := db.ConnectToDB(dbParams)
	dbParams = nil
	if err != nil {
		fmt.Println("Database offline!")
		panic(err)
	}

	err = db.CreateTables(database)
	if err != nil {
		panic(err)
	}

	server.Start(database)
}
