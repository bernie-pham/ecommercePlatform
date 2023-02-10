package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreatePCriteriaRequest struct {
	Criteria string `json:"criteria"`
}

func (server *Server) CreatePCriteria(ctx *gin.Context) {
	var req CreatePCriteriaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	PCriteria, err := server.store.CreatePCriteria(ctx, req.Criteria)
	if err != nil {
		log.Error().Err(err).Msg("failed to create product colour")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, PCriteria)
}

func (server *Server) ListPCriterias(ctx *gin.Context) {
	PCriterias, err := server.store.ListPCriterias(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list product criterias")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, PCriterias)
}
