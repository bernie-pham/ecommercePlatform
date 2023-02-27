package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bernie-pham/ecommercePlatform/pb"
	"github.com/bernie-pham/ecommercePlatform/session"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RefreshToken(
	ctx context.Context, req *pb.RefreshTokenReq,
) (*pb.RefreshTokenResponse, error) {
	token := req.GetRefreshToken()

	payload, err := server.sessionRepo.GetRefreshPayload(ctx, token)
	if err != nil {
		if err == redis.Nil {
			log.Error().
				Err(err).
				Str("token", token).
				Msg("invalid token")
			return nil, err
		}
		log.Error().
			Err(err).
			Str("token", token).
			Msg("failed to get refresh paylaod from redis with key")
		return nil, err
	}
	var sessionPayload session.SessionPayload
	err = json.Unmarshal([]byte(payload), &sessionPayload)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to unmarshal payload")
		return nil, err
	}
	prevTokenID := strings.Split(token, ":")[1]
	err = server.sessionRepo.DeleteRefreshToken(ctx, sessionPayload.UserID, prevTokenID)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to delete all related token")
		return nil, err
	}
	session_payload, access_token_id, err := session.NewPayload(
		sessionPayload.UserID,
		sessionPayload.AccessLevel,
		sessionPayload.Email,
		sessionPayload.UserAgent,
	)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to create access token")
		return nil, err
	}
	refresh_token_id, err := session.NewRefreshTokenID()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to refresh token id")
		return nil, err
	}
	err = server.sessionRepo.SetAccessToken(
		ctx,
		sessionPayload.UserID,
		access_token_id,
		server.config.AccessTimeout,
		session_payload,
	)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to set access token to redis")
		return nil, err
	}
	err = server.sessionRepo.SetRefreshToken(
		ctx,
		sessionPayload.UserID,
		refresh_token_id,
		server.config.RefreshTimeout,
		session_payload,
	)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to refresh token id")
		return nil, err
	}
	rsp := &pb.RefreshTokenResponse{
		AccessToken:         fmt.Sprintf("%d:%s", sessionPayload.UserID, access_token_id),
		AccessTokenTimeout:  timestamppb.New(time.Now().Add(server.config.AccessTimeout)),
		NewRefreshToken:     fmt.Sprintf("%d:%s", sessionPayload.UserID, refresh_token_id),
		RefreshTokenTimeout: timestamppb.New(time.Now().Add(server.config.RefreshTimeout)),
	}
	return rsp, nil
}
