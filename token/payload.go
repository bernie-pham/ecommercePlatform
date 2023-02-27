package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TypeKey string

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token invalid")
	AuthPayloadKey  = TypeKey("auth_payload")
)

type Payload struct {
	ID             uuid.UUID
	Email          string
	UserID         int
	CurrentOrderID int
	Access_level   int
	Issue_at       time.Time
	Expire_at      time.Time
}

func NewPayload(email string, duration time.Duration, access_level int, userID int) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	issue_at := time.Now()
	payload := &Payload{
		ID:           id,
		Email:        email,
		UserID:       userID,
		Access_level: access_level,
		Issue_at:     issue_at,
		Expire_at:    issue_at.Add(duration),
	}
	return payload, nil
}

func (payload *Payload) IsValidThrough() error {
	if time.Now().After(payload.Expire_at) {
		return ErrExpiredToken
	}
	return nil
}
