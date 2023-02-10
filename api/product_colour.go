package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreatePColourRequest struct {
	ColourName string `json:"colour_name"`
}

func (server *Server) CreatePColour(ctx *gin.Context) {
	var req CreatePColourRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	pColour, err := server.store.CreatePColour(ctx, req.ColourName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create product colour")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, pColour)
}

func (server *Server) ListPColours(ctx *gin.Context) {
	pColours, err := server.store.ListPColours(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list product colours")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, pColours)
}
