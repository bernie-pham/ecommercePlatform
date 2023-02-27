package session

import (
	"context"
	"time"
)

type SessionRepository interface {
	SetRefreshToken(ctx context.Context, userID int64, tokenID string, expiresIn time.Duration, payload []byte) error
	SetAccessToken(ctx context.Context, userID int64, tokenID string, expiresIn time.Duration, payload []byte) error
	DeleteRefreshToken(ctx context.Context, userID int64, prevTokenID string) error
	// DeleteUserRefreshTokens looking for all refresh/access token stored in redis, begining with user_id.
	// this function is used in case of logout, any security violation related
	DeleteUserRelatedTokens(ctx context.Context, userID int64) error
	GetAccessPayload(ctx context.Context, key string) (string, error)
	GetRefreshPayload(ctx context.Context, key string) (string, error)
}
