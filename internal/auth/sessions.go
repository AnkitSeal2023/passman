package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"passman/internal/db"
)

func isValidSessionToken(queries *db.Queries, session_token string, username string) error {
	user, err := queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		log.Printf("sessions.go:16 : %v", err)
		return err
	}
	if user.SessionToken.Valid && user.SessionToken.String == session_token {
		return nil
	}

	return errors.New("Invalid Session ID")
}

func GenerateSessionId() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
