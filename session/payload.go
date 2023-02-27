package session

import (
	"encoding/json"

	"github.com/google/uuid"
)

type SessionPayload struct {
	UserID      int64
	AccessLevel int8
	Email       string
	UserAgent   string
}

func NewPayload(userID int64, accessLevel int8, email, userAgent string) ([]byte, string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, "", err
	}
	payload := &SessionPayload{
		UserID:      userID,
		AccessLevel: accessLevel,
		Email:       email,
		UserAgent:   userAgent,
	}
	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}
	return marshalledPayload, id.String(), nil
}

func NewRefreshTokenID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
