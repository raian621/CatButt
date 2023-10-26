package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"catbook.com/auth"
	"catbook.com/db"
	"catbook.com/server"
	"catbook.com/util"
)

func TestSignIn(t *testing.T) {
	util.LoadEnvVars("../../.env")
	database, err := server.NewDatabase()
	util.DeleteDatabaseEnvVars()
	defer database.Close()

	if err != nil {
		t.Error(err)
	}

	userRegInfo := auth.UserRegistrationInfo{
		Username: "jeffbezos",
		Password: "megayacht_420",
		Email:    "jeff.bezos@amazon.com",
	}

	db.CreateUser(&userRegInfo, database)
	user, err := db.GetUserByUsername(userRegInfo.Username, database)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err := db.DeleteUser(user.UserId, database); err != nil {
			t.Error(err)
			return
		}
	}()

	rawJSON, err := json.Marshal(auth.UserCredentials{
		Username: userRegInfo.Username,
		Password: userRegInfo.Password,
	})

	if err != nil {
		t.Error(err)
		return
	}

	r := server.NewRouter(database)
	bodyReader := bytes.NewReader(rawJSON)
	req, _ := http.NewRequest("POST", "/signin", bodyReader)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d\n", http.StatusOK, w.Code)
	}
	cookies := w.Result().Cookies()
	sessionIdSet := false
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "sessionid" {
			sessionIdSet = true
			sessionCookie = cookie
			t.Log(sessionCookie)
			break
		}
	}

	if !sessionIdSet {
		t.Error("sessionid cookie should have been set")
		return
	}

	bodyReader = bytes.NewReader(rawJSON)
	req, _ = http.NewRequest("POST", "/signin", bodyReader)
	req.AddCookie(sessionCookie)
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected %d status code, got %d\n", http.StatusAccepted, w.Code)
		return
	}
}
