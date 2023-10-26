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
			requiredEnvVars := []string{
				"DB_PROV", "DB_HOST", "DB_USER", "DB_PASS",
				"DB_NAME", "DB_PORT", "DB_SSL",
			}
			for _, envVar := range requiredEnvVars {
				if os.Getenv(envVar) == "" {
					panic(err)
				}
			}
		}
	}
}

func DeleteDatabaseEnvVars() {
	os.Unsetenv("DB_PROV")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_SSL")
}
