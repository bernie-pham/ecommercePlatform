package gapi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bernie-pham/ecommercePlatform/pb"
	"github.com/bernie-pham/ecommercePlatform/session"
	ultils "github.com/bernie-pham/ecommercePlatform/ultil"
	val "github.com/bernie-pham/ecommercePlatform/validator"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
)

func (server *Server) LoginMerchant(
	ctx context.Context,
	req *pb.LoginMerchantRequest,
) (*pb.LoginMerchantResponse, error) {
	violations := validateLoginRequest(req)
	if violations != nil {
		return nil, InvalidArgumentError(violations)
	}
	email := req.GetEmail()
	password := req.GetPassword()

	user, err := server.store.GetUser(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, MerchantEmailNotFound("email")
		}
		return nil, err
	}
	err = ultils.VerifyHashedPassword(user.HashedPassword, password)
	if err != nil {
		return nil, fmt.Errorf("wrong password, try again")
	}

	// Login successful
	// Manage session with DB
	// accessToken, payload, err := server.tokenMaker.CreateToken(
	// 	server.config.AccessTimeout,
	// 	email,
	// 	int(user.AccessLevel),
	// 	token.AccessTokenType,
	// 	int(user.ID),
	// )
	// if err != nil {
	// 	return nil, errors.New("failed to generate accessToken")
	// }
	// refreshToken, refresh_payload, err := server.tokenMaker.CreateToken(
	// 	server.config.RefreshTimeout,
	// 	email,
	// 	int(user.AccessLevel),
	// 	token.RefreshTokenType,
	// 	int(user.ID))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate refreshToken")
	// }

	// rsp := &pb.LoginMerchantResponse{
	// 	AccessToken:         accessToken,
	// 	AccessTokenTimeout:  timestamppb.New(payload.Expire_at),
	// 	RefreshToken:        refreshToken,
	// 	RefreshTokenTimeout: timestamppb.New(refresh_payload.Expire_at),
	// }
	// var userAgent, ClientIP string
	var userAgent string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// fmt.Printf("%+v%+v", md.Get("user-agent"), ok)
		userAgent = md.Get("user-agent")[0]
	}
	// generate access token
	access_payload, access_token_id, err := session.NewPayload(
		user.ID,
		int8(user.AccessLevel),
		user.Email,
		userAgent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create access payload: %v", err)
	}

	refresh_token_id, err := session.NewRefreshTokenID()
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token id: %v", err)
	}

	// before set new access/refresh token, check existing token from previous login session and delete all of them
	err = server.sessionRepo.DeleteUserRelatedTokens(ctx, user.ID)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to remove all previous tokens from userID: %d", user.ID)
		return nil, fmt.Errorf("failed to remove old tokens: %v", err)
	}

	err = server.sessionRepo.SetAccessToken(ctx, user.ID, access_token_id, server.config.AccessTimeout, access_payload)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to set access token to redis for userID: %d", user.ID)
		return nil, fmt.Errorf("failed to set access token: %v", err)
	}

	err = server.sessionRepo.SetRefreshToken(ctx, user.ID, refresh_token_id, server.config.RefreshTimeout, access_payload)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to set refresh token to redis for userID: %d", user.ID)
		return nil, fmt.Errorf("failed to set refresh token: %v", err)
	}

	rsp := &pb.LoginMerchantResponse{
		AccessToken:  fmt.Sprintf("%v:%s", user.ID, access_token_id),
		RefreshToken: fmt.Sprintf("%v:%s", user.ID, refresh_token_id),
	}

	return rsp, nil
}

func validateLoginRequest(req *pb.LoginMerchantRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, FieldViolation("email", err))
	}
	return
}
