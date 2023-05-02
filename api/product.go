package api

import (
	"net/http"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ListProductParamsRequest struct {
	PageID   int32 `form:"page_id"`
	PageSize int32 `form:"page_size"`
}

func (server *Server) ListProducts(ctx *gin.Context) {
	var req ListProductParamsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var limit, offset int32
	limit = req.PageSize
	offset = req.PageSize * req.PageID

	arg := db.ListAllProductsParams{
		Limit:  limit,
		Offset: offset,
	}

	products, err := server.store.ListAllProducts(ctx, arg)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to list all products")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, products)
}
