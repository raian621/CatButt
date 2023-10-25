package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/mileusna/useragent"
)

type Session struct {
	SessionId string
	Username  string
	Expires   time.Time
	OS        string
	Browser   string
	IP        string
}

func CreateSession(username, userAgent, clientIP string, database *sql.DB) *Session {
	ua := useragent.Parse(userAgent)
	session := Session{
		SessionId: uuid.New().String(),
		Username:  username,
		Expires:   time.Now().Add(time.Hour * 24 * 7), // sessions expire after 1 week
		OS:        ua.OS,
		Browser:   ua.Name,
		IP:        clientIP,
	}

	return &session
}

func DestroySession(sessionId string) {

}
