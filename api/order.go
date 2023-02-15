package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/bernie-pham/ecommercePlatform/async"
	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateOrderRequest struct {
	Items  []int `json:"items"  binding:"required"`
	DealID int   `json:"deal_id"`
}

// AddOrderItem will take list of item.
func (server *Server) CreateOrder(ctx *gin.Context) {
	session_value, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		err := errors.New("invalid access token")
		log.Error().
			Err(err).
			Msg("failed to get userID from access payload")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	auth_payload := session_value.(*token.Payload)

	var req CreateOrderRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Msg("failed to get item id list from request")
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}
	// checking these item in list are avaialble or not?
	// checking price, sale, quantity
	var cartItems []db.CartItem
	for _, id := range req.Items {
		item, err := server.store.GetCartItemByID(ctx, int64(id))
		if err != nil {
			if err == sql.ErrNoRows {
				log.Error().
					Err(err).
					Msg("cart item does not exist")
				ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
				return
			}
			log.Error().
				Err(err).
				Msg("failed to get cart item")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
			return
		}
		cartItems = append(cartItems, item)
	}
	arg := db.CreateOrderTXParams{
		UserID: auth_payload.UserID,
		Items:  cartItems,
		DealID: req.DealID,
		AfterCreateFunc: func(merchantID_OrderID_map map[int32]int64) error {
			for merchantID, merchant_order_id := range merchantID_OrderID_map {
				arg := async.NotificationPayload{
					RecipientID: int(merchantID),
					Title:       "Ecommerce System: You got an order",
					Msg: fmt.Sprintf(
						"You got an order with ID: %v", merchant_order_id),
				}
				err := server.taskDistributor.DistributeTaskNotification(ctx, &arg)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	err := server.store.CreateOrd	erTX_v2(ctx, arg)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to create order for user: %s", auth_payload.UserID)
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
}
