package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/mileusna/useragent"
)

type Session struct {
	SessionId string
	UserId    string
	Expires   time.Time
	OS        string
	Browser   string
	IP        string
}

func NewSession(userId, userAgent, clientIP string, database *sql.DB) (*Session, error) {
	ua := useragent.Parse(userAgent)
	sessionId, err := GenerateSessionId()
	if err != nil {
		return nil, err
	}

	session := Session{
		SessionId: sessionId,
		UserId:    userId,
		Expires:   time.Now().Add(time.Hour * 24 * 7), // sessions expire after 1 week
		OS:        ua.OS,
		Browser:   ua.Name,
		IP:        clientIP,
	}

	return &session, nil
}

func CreateSession(session *Session, database *sql.DB) error {
	_, err := database.Exec(`
		INSERT INTO sessions (id, user_id, expires, os, browser, ip_address)
			VALUES ($1, $2, $3, $4, $5, $6);
		`,
		session.SessionId,
		session.UserId,
		session.Expires,
		session.OS,
		session.Browser,
		session.IP,
	)

	return err
}

func DestroySession(sessionId string, database *sql.DB) error {
	_, err := database.Exec(`DELETE FROM sessions WHERE id=$1`, sessionId)
	return err
}

func GenerateSessionId() (string, error) {
	sessionIdBytes := make([]byte, 24)
	binary.BigEndian.PutUint32(sessionIdBytes[:4], uint32(time.Now().Unix()))
	if _, err := rand.Read(sessionIdBytes[4:]); err != nil {
		return "", err
	}

	sessionId := base64.RawURLEncoding.EncodeToString(sessionIdBytes)

	return sessionId, nil
}

func ValidSession(sessionId string, database *sql.DB) (bool, error) {
	row := database.QueryRow(`SELECT expires FROM sessions WHERE id=$1;`, sessionId)
	var expires time.Time
	if err := row.Scan(&expires); err != nil {
		return false, err
	}

	if expires.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}

func RefreshSession(sessionId string, database *sql.DB) error {
	_, err := database.Exec("UPDATE sessions SET expires=$1 WHERE id=$2",
		time.Now().Add(time.Hour*24*7), sessionId)
	return err
}
