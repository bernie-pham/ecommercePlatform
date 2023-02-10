package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreatePSizeRequest struct {
	SizeValue string `json:"size_value"`
}

func (server *Server) CreatePSize(ctx *gin.Context) {
	var req CreatePSizeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	PSize, err := server.store.CreatePSize(ctx, req.SizeValue)
	if err != nil {
		log.Error().Err(err).Msg("failed to create product size")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, PSize)
}

func (server *Server) ListPSizes(ctx *gin.Context) {
	PSizes, err := server.store.ListPSizes(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list product sizes")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, PSizes)
}
