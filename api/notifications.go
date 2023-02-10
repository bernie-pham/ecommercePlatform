package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bernie-pham/ecommercePlatform/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *Server) ListNofications(ctx *gin.Context) {
	session_payload, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		err := errors.New("invalid access token")
		log.Error().
			Err(err).
			Msg("failed to get userID from access payload")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	auth_payload := session_payload.(*token.Payload)
	notifications, err := server.store.ListNotifications(ctx, int64(auth_payload.UserID))
	if err != nil && err != sql.ErrNoRows {
		log.Error().
			Err(err).
			Msg("failed to list notifications")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, notifications)
}
