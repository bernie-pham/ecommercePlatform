package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/token"
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

type getMerchantOrderRequest struct {
	MerchantOrderID int64 `uri:"id"`
}

type merchantOrderDetails struct {
	OrderTotalPrice float32                                 `json:"order_total_price"`
	Status          string                                  `json:"status"`
	Created_At      time.Time                               `json:"created_at"`
	Updated_At      time.Time                               `json:"updated_at"`
	OrderItems      []db.ListOrderItemsByMerchantOrderIDRow `json:"order_items"`
}

func (server *Server) GetMerchantOrderDetails(ctx *gin.Context) {
	session_payload, ok := ctx.Get(authorizationHeaderKey)
	if !ok {
		err := errors.New("invalid access token")
		log.Error().
			Err(err).
			Msg("failed to get auth payload from session")
		return
	}
	authPayload := session_payload.(*token.Payload)
	var req getMerchantOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Error().
			Err(err).
			Msg("Bad Request: failed to get arg from request")
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}
	arg := db.GetMerchantOrderParams{
		ID:         req.MerchantOrderID,
		MerchantID: int64(authPayload.UserID),
	}
	merchantOrder, err := server.store.GetMerchantOrder(ctx, arg)
	if err != nil && err != sql.ErrNoRows {
		log.Error().
			Err(err).
			Msg("failed to get merchant order")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	orderItems, err := server.store.
		ListOrderItemsByMerchantOrderID(ctx, merchantOrder.ID)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to get merchant order item")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	rsp := merchantOrderDetails{
		OrderTotalPrice: merchantOrder.TotalPrice,
		Status:          string(merchantOrder.OrderStatus),
		Created_At:      merchantOrder.CreatedAt,
		Updated_At:      merchantOrder.UpdatedAt,
		OrderItems:      orderItems,
	}
	ctx.JSON(http.StatusOK, rsp)
}
