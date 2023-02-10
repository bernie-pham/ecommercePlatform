package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *Server) ListProducts(ctx *gin.Context) {
	products, err := server.store.ListAllProducts(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to list all products")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, products)
}
