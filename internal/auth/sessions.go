package auth

import (
	"errors"
)

func isValidSessionId(session_id string) error {
	if len((session_id)) > 0 {
		return nil
	}
	return errors.New("Invalid Session ID")
}
