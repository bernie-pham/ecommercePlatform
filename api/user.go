package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/bernie-pham/ecommercePlatform/async"
	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/token"
	ultils "github.com/bernie-pham/ecommercePlatform/ultil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type CreateUserRequest struct {
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}
type LoginUserRequest struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type userResponse struct {
	ID                int       `json:"user_id"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type LoginUserResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpireAt  time.Time `json:"access_token_expire_at"`
	RefreshToken         string    `json:"refresh_token"`
	RefreshTokenExpireAt time.Time `json:"refresh_token_expire_at"`
	userResponse
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                int(user.ID),
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordUpdatedAt,
		Phone:             user.Phone,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedpassowrd, err := ultils.HashPassword(req.Password)
	if err != nil {
		log.Err(err).Msg("failed to hash password")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user_arg := db.CreateUserParams{
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
		HashedPassword: hashedpassowrd,
	}

	user, err := server.store.CreateUser(ctx, user_arg)
	if err != nil {
		log.Err(err).Msg("failed to create user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// password verification
	err = ultils.VerifyHashedPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrWrongPassword))
		return
	}
	// generate session token
	access_token, access_payload, err := server.tokenMaker.CreateToken(
		server.config.AccessTimeout,
		user.Email,
		int(user.AccessLevel),
		token.AccessTokenType,
		int(user.ID),
	)

	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to generate access token: %s", user.Email)
	}

	refresh_token, refresh_payload, err := server.tokenMaker.CreateToken(
		server.config.RefreshTimeout,
		user.Email,
		int(user.AccessLevel),
		token.RefreshTokenType,
		int(user.ID))

	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to generate refresh token: %s", user.Email)
	}

	userResponse := newUserResponse(user)
	loginUserResponse := LoginUserResponse{
		AccessToken:          access_token,
		AccessTokenExpireAt:  access_payload.Expire_at,
		RefreshToken:         refresh_token,
		RefreshTokenExpireAt: refresh_payload.Expire_at,
		userResponse:         userResponse,
	}
	ctx.JSON(http.StatusOK, loginUserResponse)
}

func (server *Server) ForgotPassword(ctx *gin.Context) {
	// take a email string out of the request url
	email := ctx.DefaultQuery("email", "none")
	if email == "none" {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}
	// ascertain where given user's email is existed or not
	user, err := server.store.GetUser(ctx, email)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to get user with email: %s", email)
	}
	// Create verification row
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to generate uuid for user: %s", email)
	}

	arg := db.CreateVerificationParams{
		ID:        id,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(server.config.VerifyTimeout),
	}

	verification, err := server.store.CreateVerification(ctx, arg)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to create verification row for user: %s", email)
	}

	// Generate a reset password link
	resetPasswordLink := fmt.Sprintf(
		"http://%s/reset_password?reset_token=%s",
		server.config.HTTPServerAddr,
		verification.ID,
	)
	// Send an Email which contain reset password link
	emailPayload := &async.EmailDeliveryPayload{
		Email:        user.Email,
		Url:          resetPasswordLink,
		Subject:      "Password Recovery",
		Msg:          "This email is generated automatically by ecommerce system",
		EmailTemplte: server.config.MailTempltePath,
	}
	server.taskDistributor.DistributeTaskSendMail(ctx, emailPayload)
	// TODO create send email func
	log.Info().
		Msg("Reset Password link generated")

	ctx.JSON(http.StatusOK, nil)
}
func (server *Server) ResetPassword(ctx *gin.Context) {
	// Get Reset ID
	reset_id := ctx.DefaultQuery("reset_token", "none")
	if reset_id == "none" || len(reset_id) == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}
	// Get New Password
	new_password := ctx.DefaultQuery("new_password", "none")
	if new_password == "none" || len(new_password) == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}

	// Ascertain Reset ID
	verification, err := server.store.GetVerification(ctx, uuid.MustParse(reset_id))
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to get verification")
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}
	// Reset Password TX
	hashed_password, err := ultils.HashPassword(new_password)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to hash a password")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = server.store.ResetPasswdTX(ctx, db.ResetPasswdTXParams{
		Email:             verification.Email,
		VerificationID:    verification.ID,
		NewHashedPassword: hashed_password,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to reset password")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

type UpdateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Fullname string `json:"full_name"`
	Phone    string `json:"phone"`
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fullname, isFullname := getOptionalString(req.Fullname)
	phone, isPhone := getOptionalString(req.Phone)
	arg := db.UpdateUserParams{
		FullName: sql.NullString{
			String: fullname,
			Valid:  isFullname,
		},
		Phone: sql.NullString{
			String: phone,
			Valid:  isPhone,
		},
		Email: req.Email,
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrNotFound))
			return
		}
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}
