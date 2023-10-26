package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"catbook.com/server"
	"catbook.com/util"
)

func TestRegisterRoute(t *testing.T) {
	util.LoadEnvVars("../../.env")
	database, err := server.NewDatabase()
	defer database.Close()
	util.DeleteDatabaseEnvVars()

	if err != nil {
		t.Error(err)
	}

	r := server.NewRouter(database)
	w := httptest.NewRecorder()
	rawJSON := []byte(`{"username":"jeffbezos","password":"superyacht_420","email":"jeff.bezos@amazon.com"}`)
	bodyReader := bytes.NewReader(rawJSON)
	req, _ := http.NewRequest("POST", "/register", bodyReader)

	defer func() {
		if _, err := database.Exec(`DELETE FROM users WHERE username='jeffbezos';`); err != nil {
			t.Error(err)
		}
	}()
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("Server should have returned a 201 status code!")
	}

	bodyReader = bytes.NewReader(rawJSON)
	req, _ = http.NewRequest("POST", "/register", bodyReader)
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Log(w.Code)
		t.Error("Server should have returned an 409 status code!")
	}
}

func TestRegisterRouteInvalidFields(t *testing.T) {
	util.LoadEnvVars("../../.env")
	database, err := server.NewDatabase()
	defer database.Close()
	util.DeleteDatabaseEnvVars()

	if err != nil {
		t.Error(err)
	}

	// empty username
	r := server.NewRouter(database)
	w := httptest.NewRecorder()
	rawJSON := []byte(`{"username":"","password":"superyacht_420","email":"jeff.bezos@amazon.com"}`)
	bodyReader := bytes.NewReader(rawJSON)
	req, _ := http.NewRequest("POST", "/register", bodyReader)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("Server should have returned a 400 status code!")
	}

	// empty password
	w = httptest.NewRecorder()
	rawJSON = []byte(`{"username":"jeffbezos","password":"","email":"jeff.bezos@amazon.com"}`)
	bodyReader = bytes.NewReader(rawJSON)
	req, _ = http.NewRequest("POST", "/register", bodyReader)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("Server should have returned a 400 status code!")
	}

	// empty email
	w = httptest.NewRecorder()
	rawJSON = []byte(`{"username":"jeffbezos","password":"","email":""}`)
	bodyReader = bytes.NewReader(rawJSON)
	req, _ = http.NewRequest("POST", "/register", bodyReader)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("Server should have returned a 400 status code!")
	}

	// empty email
	w = httptest.NewRecorder()
	rawJSON = []byte(`{"username":"jeffbezos","password":"","email":""}`)
	bodyReader = bytes.NewReader(rawJSON)
	req, _ = http.NewRequest("POST", "/register", bodyReader)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("Server should have returned a 400 status code!")
	}
}
