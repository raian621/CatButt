package main

import (
	"fmt"
	"os"

	"catbook.com/server"
	"catbook.com/util"
)

func main() {
	util.LoadEnvVars(".env")
	database, err := server.NewDatabase()
	util.DeleteDatabaseEnvVars()

	if err != nil {
		panic(err)
	}

	r := server.NewRouter(database)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if host == "" || port == "" {
		r.Run()
	} else {
		r.Run(fmt.Sprintf("%s:%s", host, port))
	}
}
