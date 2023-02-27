package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisSessionRepository struct {
	Redis *redis.Client
}

func NewRedisSessionRepo(redisClient *redis.Client) SessionRepository {
	return &RedisSessionRepository{
		Redis: redisClient,
	}
}

func (r *RedisSessionRepository) SetRefreshToken(
	ctx context.Context,
	userID int64,
	tokenID string,
	expiresIn time.Duration,
	payload []byte,
) error {
	key := fmt.Sprintf("%v:%s", userID, tokenID)
	// at the momment, we dont set any value for refresh token
	// But will do later
	if err := r.Redis.SetEx(ctx, key, payload, expiresIn).Err(); err != nil {
		log.Error().
			Err(err).
			Msgf("could not set refresh token to redis for userid/tokenid: %s:%s", userID, tokenID)
		return err
	}
	return nil
}
func (r *RedisSessionRepository) SetAccessToken(
	ctx context.Context,
	userID int64,
	tokenID string,
	expiresIn time.Duration,
	payload []byte,
) error {
	key := fmt.Sprintf("%v:%s", userID, tokenID)
	// at the momment, we dont set any value for refresh token
	// But will do later
	if err := r.Redis.SetEx(ctx, key, payload, expiresIn).Err(); err != nil {
		log.Error().
			Err(err).
			Msgf("could not set access token to redis for userid/tokenid: %s:%s", userID, tokenID)
		return err
	}
	return nil
}

// DeleteRefreshToken used to delete old user's refresh token
func (r *RedisSessionRepository) DeleteRefreshToken(
	ctx context.Context,
	userID int64,
	prevTokenID string,
) error {
	key := fmt.Sprintf("%d:%s", userID, prevTokenID)
	result := r.Redis.Del(ctx, key)
	if result.Err() != nil {
		log.Error().
			Err(result.Err()).
			Msgf("failed to remove old refresh token with key: %s", key)
		return result.Err()
	}
	if result.Val() < 1 {
		log.Error().
			Err(result.Err()).
			Msgf("invalid refresh token: %s", key)
		return errors.New("invalid refresh token")
	}
	return nil
}

// DeleteUserRefreshTokens looking for all refresh/access token stored in redis, begining with user_id.
// this function is used in case of logout, any security violation related
func (r *RedisSessionRepository) DeleteUserRelatedTokens(
	ctx context.Context,
	userID int64,
) error {
	pattern := fmt.Sprintf("%v:*", userID)

	// in case there is less than 50 access/refresh token has been created.
	// if more than, we should check the returned cursor and keep scanning for the next 50 counts and so on until the end.
	iter := r.Redis.Scan(ctx, 0, pattern, 50).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			log.Error().
				Err(err).
				Msgf("failed to remote key in redis with key: %v", iter.Val())
			failCount++
		}
	}

	if err := iter.Err(); err != nil {
		log.Error().
			Err(err).
			Msgf("failed to remote key in redis with key: %v", iter.Val())
	}
	// TODO: implement redis distribute/handler for this task in case any corrupted occurs
	if failCount > 0 {
		return fmt.Errorf("failed to remove all key for user_id: %v", userID)
	}
	return nil
}

func (r *RedisSessionRepository) GetAccessPayload(ctx context.Context, key string) (string, error) {
	result := r.Redis.Get(ctx, key)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Result()
}
func (r *RedisSessionRepository) GetRefreshPayload(ctx context.Context, key string) (string, error) {
	result := r.Redis.Get(ctx, key)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Result()
}
