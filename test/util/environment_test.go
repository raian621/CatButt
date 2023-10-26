package util_test

import (
	"os"
	"testing"

	"catbook.com/util"
)

func TestLoadEnvVarsDevelopment(t *testing.T) {
	// by default, .env file will not exist
	os.Setenv("MODE", "development")
	_, err := os.Open(".env")

	shouldPanic := os.IsNotExist(err)
	if shouldPanic {
		t.Log(".env file not found, LoadEnvVars should panic.")
	}
	defer func() {
		err := recover()
		if shouldPanic && err == nil {
			t.Error("Function should have panicked.")
		}
		// else if !shouldPanic && err != nil {
		// 	t.Error("Function should not have panicked.", err)
		// }
	}()

	util.LoadEnvVars(".env")
}
