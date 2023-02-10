package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoTokenMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symKey string) (TokenMaker, error) {
	if len(symKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoTokenMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symKey),
	}
	return maker, nil
}

func (maker *PasetoTokenMaker) CreateToken(duration time.Duration, email string, access_level int, token_type string, userID int) (string, *Payload, error) {
	payload, err := NewPayload(email, duration, access_level, userID)
	if err != nil {
		return "", nil, err
	}
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, token_type)
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

func (maker *PasetoTokenMaker) VerifyToken(token string) (*Payload, string, error) {
	payload := &Payload{}
	var footer string
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, &footer)
	if err != nil {
		return nil, "", err
	}
	err = payload.IsValidThrough()
	if err != nil {
		return nil, "", err
	}
	return payload, footer, nil
}
