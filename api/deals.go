package api

import (
	"net/http"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type getOrdDealRequest struct {
	Code string `form:"code" binding:"required"`
}

type getOrderDealRsp struct {
	Description   string  `json:"description"`
	Code          string  `json:"code"`
	Discount_rate float32 `json:"dis_rate"`
	Limit         float64 `json:"limit"`
}

// getDeal URI GET /deals?code=[paramter]
func (server *Server) getOrdDeal(ctx *gin.Context) {
	var req getOrdDealRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetDealByCodeParams{
		Code: req.Code,
		Type: "ord_dis",
	}

	deal, err := server.store.GetDealByCode(ctx, arg)

	if err != nil {
		log.Error().
			Err(err).
			Str("deal's code", arg.Code).
			Msg("failed to get Deal info")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}

	rsp := &getOrderDealRsp{
		Description:   deal.Name,
		Code:          deal.Code,
		Discount_rate: deal.DiscountRate,
		Limit:         deal.DealLimit.Float64,
	}
	ctx.JSON(http.StatusOK, rsp)
}
