package interceptor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bernie-pham/ecommercePlatform/session"
	"github.com/bernie-pham/ecommercePlatform/token"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (interceptor *Interceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	accessLevel, ok := interceptor.accessLevels[method]
	if !ok {
		// non-authenticated required method, ex: login
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	value := md["authorization"]
	if len(value) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header is not found")
	}
	values := strings.Fields(value[0])

	log.Info().Str("accessLevel ", fmt.Sprint(accessLevel[0])).Msgf("authorization: %v", values[0])
	auth_type := values[0]
	if strings.ToLower(auth_type) != "bearer" {
		return nil, status.Error(codes.Unauthenticated, "authoritzation type is not valid")
	}

	auth_token := values[1]
	auth_payload, token_type, err := interceptor.tokenMaker.VerifyToken(auth_token)
	if err != nil {
		if err == token.ErrExpiredToken {
			return nil, status.Error(codes.Unauthenticated, "authorization token has been expired")
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if token_type != token.AccessTokenType {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not valid")
	}

	access_level := auth_payload.Access_level
	notfound := true
	for _, level := range accessLevel {
		if level == access_level {
			notfound = false
			break
		}
	}

	if notfound {
		return nil, status.Error(codes.PermissionDenied, "unauthorized method")
	}

	newCtx := context.WithValue(ctx, token.AuthPayloadKey, auth_payload)
	return newCtx, nil
}

func (interceptor *Interceptor) authorizeWithRedis(ctx context.Context, method string) (context.Context, error) {
	accessLevel, ok := interceptor.accessLevels[method]
	if !ok {
		// non-authenticated required method, ex: login
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	sessions := md.Get("access_token_id")
	if len(sessions) == 0 {
		return nil, status.Error(codes.Unauthenticated, "session id is not found in metadata")
	}
	access_token_id := sessions[0]

	session_object, err := interceptor.sessionRepo.GetAccessPayload(ctx, access_token_id)
	if err != nil {
		if err == redis.Nil {
			return nil, status.Error(codes.Unauthenticated, "expired session")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user session: %v", err)
	}
	var session_payload session.SessionPayload
	err = json.Unmarshal([]byte(session_object), &session_payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshall user session: %v", err)
	}
	// fmt.Printf("session_object %+v", session_payload)
	notfound := true
	for _, level := range accessLevel {
		if level == int(session_payload.AccessLevel) {
			notfound = false
			break
		}
	}
	if notfound {
		return nil, status.Error(codes.PermissionDenied, "unauthorized method")
	}

	newCtx := context.WithValue(ctx, token.AuthPayloadKey, session_payload)
	return newCtx, nil
}
