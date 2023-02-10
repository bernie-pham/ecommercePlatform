package api

import (
	"database/sql"
	"net/http"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateMerchantRequest struct {
	UserID      int    `json:"user_id" binding:"required"`
	CountryCode int    `json:"country_code" binding:"required"`
	Description string `json:"desc" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

func (server *Server) CreateMerchant(ctx *gin.Context) {
	var req CreateMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateMerchantParams{
		UserID:       int64(req.UserID),
		CountryCode:  int32(req.CountryCode),
		Description:  req.Description,
		MerchantName: req.Name,
	}

	merchant, err := server.store.CreateMerchant(ctx, arg)
	if err != nil {
		log.Error().
			Err(err).Msg("failed to create merchant")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, merchant)
}

type UpdateMerchantRequest struct {
	ID          int    `json:"merchant_id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	CountryCode int    `json:"country_code"`
}

func (server *Server) UpdateMerchant(ctx *gin.Context) {
	var req UpdateMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	name, isName := getOptionalString(req.Name)
	desc, isDesc := getOptionalString(req.Description)
	countryCode, isCountryCode := getOptionalInt(req.CountryCode)
	arg := db.UpdateMerchantParams{
		ID: int64(req.ID),
		MerchantName: sql.NullString{
			String: name,
			Valid:  isName,
		},
		CountryCode: sql.NullInt32{
			Int32: int32(countryCode),
			Valid: isCountryCode,
		},
		Description: sql.NullString{
			String: desc,
			Valid:  isDesc,
		},
	}
	merchant, err := server.store.UpdateMerchant(ctx, arg)
	if err != nil {
		log.Error().
			Err(err).Msg("failed to update merchant")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, merchant)
}
