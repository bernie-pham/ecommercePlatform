package token

import (
	"time"
)

var (
	RefreshTokenType = "refresh_token"
	AccessTokenType  = "access_token"
)

type TokenMaker interface {
	CreateToken(duration time.Duration, email string, access_level int, token_type string, userID int) (string, *Payload, error)
	VerifyToken(token string) (*Payload, string, error)
}
