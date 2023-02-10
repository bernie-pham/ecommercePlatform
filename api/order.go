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
	// user select item in cart and then process to order.
	// body, err := ioutil.ReadAll(ctx.Request.Body)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
	// 	return
	// }
	var req CreateOrderRequest
	// err = json.Unmarshal(body, &req)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
	// 	return
	// }

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
	var tasks []async.NotificationPayload
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
		merchant_id, err := server.store.GetMerchantByCartID(ctx, int64(id))
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to get merchant id by cart id")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
			return
		}
		task := async.NotificationPayload{
			RecipientID: int(merchant_id.Int32),
			Title:       "Ecommerce System: You got an order",
			Msg: fmt.Sprintf(
				"You got an order with Product Entry - %v, Quantity: %v",
				item.ProductEntryID, item.Quantity),
		}
		tasks = append(tasks, task)
		cartItems = append(cartItems, item)
	}
	arg := db.CreateOrderTXParams{
		UserID: auth_payload.UserID,
		Items:  cartItems,
		DealID: req.DealID,
	}
	err := server.store.CreateOrderTX(ctx, arg)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to create order for user: %s", auth_payload.UserID)
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	// after Create Order Transaction success, notify merchant for prepare stock for delivery
	// execute task
	for _, task := range tasks {
		server.taskDistributor.DistributeTaskNotification(ctx, &task)
	}

}
