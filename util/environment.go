package util

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars(envFilePath string) {
	mode := os.Getenv("CATBOOK_MODE")
	if mode == "development" || mode == "" {
		err := godotenv.Load(envFilePath)

		if err != nil {
			panic(err)
		}
	}
}

func DeleteEnvVars() {
	os.Unsetenv("DB_PROV")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_SSL")
}
