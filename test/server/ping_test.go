package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"catbook.com/server"
)

func TestPingRoute(t *testing.T) {
	r := server.NewRouter(nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error("Incorrect status code!")
	}

	var message map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &message); err != nil {
		t.Error(err)
	}

	msg, ok := message["message"].(string)
	if !ok || msg != "pong" {
		t.Error("Incorrect response recieved")
	}
}
